package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/configuration"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/handlers"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/logger"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/metrics"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/middlewares"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	bot                 handlers.InternalBot
	handlers            handlers.HandlerContainer
	shutdownFuncs       []func()
	tgUpdateTimeout     int
	updateReadersNumber int
	httpServer          *http.Server
	metrics             *metrics.Metrics
}

func NewServer(ctx context.Context, config configuration.StartupConfig) (*Server, error) {
	bot, err := tgbotapi.NewBotAPI(config.BotAPIToken)
	if err != nil {
		return nil, fmt.Errorf("can not connect to the Telegram API: %w", err)
	}
	bot.Debug = false

	dbpool, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to create a connection pool: DB URL '%s': %w",
			config.DatabaseURL,
			err,
		)
	}
	if err := dbpool.Ping(ctx); err != nil {
		logMsg := fmt.Sprintf(
			"unable to connect to the DB '%v'",
			config.DatabaseURL,
		)
		slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		return nil, fmt.Errorf("%v: %w", logMsg, err)
	}
	buildingRepo := repositories.NewBuildingRepo(dbpool)
	actorRepo := repositories.NewActorRepo(dbpool)
	userRepo := repositories.NewUserRepo(dbpool)
	buildingService := services.NewBuildingService(buildingRepo, actorRepo)
	userService := services.NewUserService(userRepo)

	registry := prom.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	registeredMetrics := metrics.NewMetrics(registry)
	prometheusHandler := promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{ErrorLog: log.Default()},
	)
	authMetricsHandler := middlewares.GetBasicAuthHandler(
		prometheusHandler,
		config.MetricsUser,
		config.MetricsPassword,
	)

	srvMux := http.NewServeMux()
	srvMux.Handle("/metrics", authMetricsHandler)
	httpServer := http.Server{
		Addr:    ":" + strconv.Itoa(config.MetricsPort),
		Handler: srvMux,
	}

	botWithMetrics := handlers.NewBotWithMetrics(bot, registeredMetrics)

	handlerContainer := handlers.NewCommandContainer(
		botWithMetrics,
		buildingService,
		userService,
		registeredMetrics,
	)
	server := Server{
		botWithMetrics,
		handlerContainer,
		[]func(){dbpool.Close},
		config.TGUpdateTimeout,
		config.UpdateReadersNumber,
		&httpServer,
		registeredMetrics,
	}
	return &server, nil
}

func (s *Server) Shutdown(timeout time.Duration) {
	// set the timeout to prevent a system hang
	timeoutFunc := time.AfterFunc(timeout, func() {
		logMsg := fmt.Sprintf(
			"timeout %v has been elapsed, force exit",
			timeout.Seconds(),
		)
		slog.Error(logMsg)
		os.Exit(0)
	})
	defer timeoutFunc.Stop()
	for _, f := range s.shutdownFuncs {
		f()
	}
}

func (s *Server) RunBot(ctx context.Context) error {
	ctx, cancelCtx := context.WithCancel(ctx)
	idleConnectionsClosed := make(chan struct{})

	go func() {
		signalCh := make(chan os.Signal, 4)
		signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		select {
		case sig := <-signalCh:
			slog.InfoContext(ctx, fmt.Sprintf("receive an OS signal '%v'", sig))
		case <-ctx.Done():
			slog.InfoContext(ctx, fmt.Sprintf("start shutdown because of context"))
		}

		shutdownCtx, shutdownCtxCancel := context.WithTimeout(ctx, 5*time.Second)
		defer shutdownCtxCancel()
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			slog.ErrorContext(ctx, "a metrics shutdown error", slog.Any(logger.ErrorKey, err))
		}
		close(idleConnectionsClosed)
		cancelCtx()
	}()

	go func() {
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			slog.ErrorContext(ctx, "can not start a server", slog.Any(logger.ErrorKey, err))
		}
		cancelCtx()
	}()

	if err := s.setBotCommands(ctx); err != nil {
		return fmt.Errorf("can not set bot commands: %w", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = s.tgUpdateTimeout
	updates := s.bot.GetUpdatesChan(u)
	var wg sync.WaitGroup
	for i := 0; i < s.updateReadersNumber; i++ {
		wg.Add(1)
		go s.receiveUpdates(ctx, updates, &wg)
	}
	slog.InfoContext(
		ctx,
		fmt.Sprintf("start to listen for updates in %v goroutines", s.updateReadersNumber),
	)

	<-idleConnectionsClosed
	wg.Wait()
	slog.InfoContext(ctx, "stopped listening for updates")
	return nil
}

func (s *Server) setBotCommands(ctx context.Context) error {
	commands := []tgbotapi.BotCommand{}
	for commandName, handler := range s.handlers.HandlersPerCommand {
		command := tgbotapi.BotCommand{
			Command:     commandName,
			Description: handler.Description,
		}
		commands = append(commands, command)
	}
	setCommandsConfig := tgbotapi.NewSetMyCommands(commands...)
	result, err := s.bot.Request(setCommandsConfig)
	if err != nil {
		slog.ErrorContext(
			ctx,
			fmt.Sprintf(
				"can not set commands: command config: %v",
				setCommandsConfig,
			),
			slog.Any(logger.ErrorKey, err),
		)
		return fmt.Errorf("can not make a request to set commands: %w", err)
	}
	if !result.Ok {
		err = fmt.Errorf(result.Description)
		slog.ErrorContext(
			ctx,
			fmt.Sprintf(
				"can not set commands: an error code '%v': a response body '%s'",
				result.ErrorCode,
				result.Result,
			),
			slog.Any(logger.ErrorKey, err),
		)
		return err
	}

	return nil
}

func (s *Server) receiveUpdates(
	ctx context.Context,
	updates tgbotapi.UpdatesChannel,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			s.metrics.ChatUpdates.Inc()
			s.handleUpdate(ctx, update)
		}
	}
}

func (s *Server) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		s.handleMessage(ctx, update.Message)
	case update.CallbackQuery != nil:
		s.handleButton(ctx, update.CallbackQuery)
	default:
		slog.DebugContext(ctx, fmt.Sprintf("an unexpected update %v", update))
	}
}

func (s *Server) handleMessage(ctx context.Context, message *tgbotapi.Message) {
	user := message.From
	if user == nil {
		s.metrics.UnexpectedUpdates.With(prom.Labels{"error": "no user"}).Inc()
		return
	}
	handlerName := message.Command()
	if message.Location != nil {
		handlerName = "nearestAddresses"
	}
	handler, ok := s.handlers.GetCommandHandler(handlerName)
	if ok {
		slog.DebugContext(ctx, "handle a command", slog.String("command", message.Command()))
		handler(ctx, message)
		return
	}
	s.handlers.ProcessCommonMessage(ctx, message)
}

func (s *Server) handleButton(ctx context.Context, query *tgbotapi.CallbackQuery) {
	var queryData handlers.NextButton
	if err := json.Unmarshal([]byte(query.Data), &queryData); err != nil {
		slog.WarnContext(
			ctx,
			fmt.Sprintf("unexpected callback data %v", query),
			slog.Any(logger.ErrorKey, err),
		)
		s.metrics.UnexpectedUpdates.With(
			prom.Labels{"error": "unexpected callback data"},
		).Inc()
		return
	}
	handler, ok := s.handlers.GetButtonHandler(queryData.Name)
	if !ok {
		logMsg := fmt.Sprintf(
			"the unexpected button name %v from the chat %v: initial message %v",
			queryData.Name,
			query.Message.Chat.ID,
			query.Message.MessageID,
		)
		slog.WarnContext(ctx, logMsg)
		s.metrics.UnexpectedUpdates.With(
			prom.Labels{"error": "unexpected button name"},
		).Inc()
		return
	}
	handler(ctx, query)
}
