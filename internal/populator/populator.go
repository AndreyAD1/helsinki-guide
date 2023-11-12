package populator

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AndreyAD1/helsinki-guide/internal"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/configuration"
	"github.com/AndreyAD1/helsinki-guide/internal/infrastructure/repositories"
)

type Populator struct {
	buildingRepo repositories.BuildingRepository
	actorRepo    repositories.ActorRepository
}

func NewPopulator(ctx context.Context, config configuration.PopulatorConfig) (*Populator, error) {
	dbpool, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		log.Printf(
			"unable to create a connection pool: DB URL '%s': %v",
			config.DatabaseURL,
			err,
		)
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
		log.Println(logMsg)
		return nil, fmt.Errorf("%v: %w", logMsg, err)
	}
	populator := Populator{
		repositories.NewBuildingRepo(dbpool),
		repositories.NewActorRepo(dbpool),
	}
	return &populator, nil
}

func (p *Populator) Run(ctx context.Context) error {
	// authorIDs := []int64{}
	// for _, author := range getAuthors(authorCell) {
	// 	savedAuthor, err := p.actorRepo.Add(ctx, author)
	// 	if err != nil && err != repositories.ErrDuplicate {
	// 		return err
	// 	}
	// }

	var building internal.Building
	_, err := p.buildingRepo.Add(ctx, building)
	if err != nil {
		return err
	}
	return nil
}
