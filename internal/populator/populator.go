package populator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/configuration"
	"github.com/AndreyAD1/helsinki-guide/internal/infrastructure/repositories"
	"github.com/AndreyAD1/helsinki-guide/internal/infrastructure/repositories/types"
)

type Populator struct {
	buildingRepo      repositories.BuildingRepository
	actorRepo         repositories.ActorRepository
	neighbourhoodRepo repositories.NeighbourhoodRepository
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
		repositories.NewNeighbourhoodRepo(dbpool),
	}
	return &populator, nil
}

func (p *Populator) Run(
	ctx context.Context,
	sheetName,
	fiFilename,
	enFilename,
	ruFilename string,
) error {
	var rowSet []*excelize.Rows
	for _, filename := range []string{fiFilename, enFilename, ruFilename} {
		source, err := excelize.OpenFile(filename)
		if err != nil {
			log.Printf("can not open a file %v: %v", filename, err)
			return err
		}
		defer func() {
			if err := source.Close(); err != nil {
				log.Printf("can not close the file %s: %v", filename, err)
			}
		}()
		rows, err := source.Rows(sheetName)
		if err != nil {
			return fmt.Errorf(
				"can not get rows of a sheet '%v': %w", sheetName, err)
		}
		defer func() {
			if err := rows.Close(); err != nil {
				log.Printf("can not close a sheet '%v'", sheetName)
			}
		}()
		rowSet = append(rowSet, rows)
	}
	fiRows, enRows, ruRows := rowSet[0], rowSet[1], rowSet[2]
	// skip the first line
	fiRows.Next()
	enRows.Next()
	ruRows.Next()

	for i := 2; fiRows.Next(); i++ {
		enRows.Next()
		ruRows.Next()

		fiRow, err := fiRows.Columns()
		enRow, err := enRows.Columns()
		ruRow, err := ruRows.Columns()

		if len(fiRow) < longitudeIdx+1 {
			log.Printf("a final or unexpected row %v: %v", i, fiRow)
			break
		}

		initialUses, err := getUses(
			fiRow[initialUseIdx],
			enRow[initialUseIdx],
			ruRow[initialUseIdx],
		)
		if err != nil {
			log.Printf("initial uses error at the line %v: %v", i, err)
			return fmt.Errorf("initial uses error at the line %v: %w", i, err)
		}
		currentuses, err := getUses(
			fiRow[currentUseIdx],
			enRow[currentUseIdx],
			ruRow[currentUseIdx],
		)
		if err != nil {
			log.Printf("current uses error at the line %v: %v", i, err)
			return fmt.Errorf("current uses error at the line %v: %w", i, err)
		}

		address, err := p.getAddress(fiRow)
		if err != nil {
			return err
		}
		constructionYear, err := getYear(fiRow[constructionYearIdx])
		if err != nil {
			return err
		}
		completionYear, err := getYear(fiRow[completionYearIdx])
		if err != nil {
			return err
		}
		latitude, err := getPointerFloat32(fiRow[latitudeIdx])
		if err != nil {
			return nil
		}
		longitude, err := getPointerFloat32(fiRow[longitudeIdx])
		if err != nil {
			return nil
		}

		authorIDs := []int64{}
		for _, author := range getAuthors(fiRow, enRow, ruRow) {
			savedAuthor, err := p.actorRepo.Add(ctx, author)
			if err != nil && !errors.Is(err, repositories.ErrDuplicate) {
				return err
			}
			authorIDs = append(authorIDs, savedAuthor.ID)
		}

		building := types.Building{
			Code:                  getPointerStr(strings.Trim(fiRow[codeIdx], "'")),
			Address:               address,
			NameFi:                getPointerStr(fiRow[nameIdx]),
			NameEn:                getPointerStr(enRow[nameIdx]),
			NameRu:                getPointerStr(ruRow[nameIdx]),
			ConstructionStartYear: constructionYear,
			CompletionYear:        completionYear,
			ComplexFi:             getPointerStr(fiRow[complexIdx]),
			ComplexEn:             getPointerStr(enRow[complexIdx]),
			ComplexRu:             getPointerStr(ruRow[complexIdx]),
			HistoryFi:             getPointerStr(fiRow[historyIdx]),
			HistoryEn:             getPointerStr(enRow[historyIdx]),
			HistoryRu:             getPointerStr(ruRow[historyIdx]),
			ReasoningFi:           getPointerStr(fiRow[reasoningIdx]),
			ReasoningEn:           getPointerStr(enRow[reasoningIdx]),
			ReasoningRu:           getPointerStr(ruRow[reasoningIdx]),
			ProtectionStatusFi:    getPointerStr(fiRow[protectionStatusIdx]),
			ProtectionStatusEn:    getPointerStr(enRow[protectionStatusIdx]),
			ProtectionStatusRu:    getPointerStr(ruRow[protectionStatusIdx]),
			InfoSourceFi:          getPointerStr(fiRow[infoSourceIdx]),
			InfoSourceEn:          getPointerStr(enRow[infoSourceIdx]),
			InfoSourceRu:          getPointerStr(ruRow[infoSourceIdx]),
			SurroundingsFi:        getPointerStr(fiRow[surroundingsIdx]),
			SurroundingsEn:        getPointerStr(enRow[surroundingsIdx]),
			SurroundingsRu:        getPointerStr(ruRow[surroundingsIdx]),
			FoundationFi:          getPointerStr(fiRow[foundationIdx]),
			FoundationEn:          getPointerStr(enRow[foundationIdx]),
			FoundationRu:          getPointerStr(ruRow[foundationIdx]),
			FrameFi:               getPointerStr(fiRow[frameIdx]),
			FrameEn:               getPointerStr(enRow[frameIdx]),
			FrameRu:               getPointerStr(ruRow[frameIdx]),
			FloorDescriptionFi:    getPointerStr(fiRow[floorDescriptionIdx]),
			FloorDescriptionEn:    getPointerStr(enRow[floorDescriptionIdx]),
			FloorDescriptionRu:    getPointerStr(ruRow[floorDescriptionIdx]),
			FacadesFi:             getPointerStr(fiRow[facadeIdx]),
			FacadesEn:             getPointerStr(enRow[facadeIdx]),
			FacadesRu:             getPointerStr(ruRow[facadeIdx]),
			SpecialFeaturesFi:     getPointerStr(fiRow[specialFeaturesIdx]),
			SpecialFeaturesEn:     getPointerStr(enRow[specialFeaturesIdx]),
			SpecialFeaturesRu:     getPointerStr(ruRow[specialFeaturesIdx]),
			Latitude_ETRSGK25:     latitude,
			Longitude_ERRSGK25:    longitude,
			AuthorIDs:             authorIDs,
			InitialUses:           initialUses,
			CurrentUses:           currentuses,
		}
		_, err = p.buildingRepo.Add(ctx, building)
		if err != nil {
			return err
		}
	}
	return nil
}

func getPointerStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func getPointerFloat32(s string) (*float32, error) {
	if s == "" {
		return nil, nil
	}
	f64, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return nil, err
	}
	f32 := float32(f64)
	return &f32, nil
}

func (p *Populator) getAddress(row []string) (types.Address, error) {
	municipality := &row[municipalityIdx]
	if *municipality == "" {
		municipality = nil
	}
	neighbourbourhood := types.Neighbourhood{
		Name:         row[neighbourhoodIdx],
		Municipality: municipality,
	}
	saved, err := p.neighbourhoodRepo.Add(context.Background(), neighbourbourhood)
	if err != nil && !errors.Is(err, repositories.ErrDuplicate) {
		municipalStr := ""
		if neighbourbourhood.Municipality != nil {
			municipalStr = *neighbourbourhood.Municipality
		}
		log.Printf(
			"can not save a neighbourhood: %v-%v",
			neighbourbourhood.Name,
			municipalStr,
		)
		return types.Address{}, err
	}
	address := types.Address{
		StreetAddress:   row[streetIdx],
		NeighbourhoodID: &saved.ID,
	}
	return address, nil
}

func getAuthors(fiRow, enRow, ruRow []string) []types.Actor {
	authorNames := strings.Split(fiRow[authorIdx], " ja ")

	authors := []types.Actor{}
	for _, useFi := range authorNames {
		authorName := strings.TrimSpace(useFi)
		if authorName == "" {
			continue
		}
		titleFi := getPointerStr(fiRow[authorTitleIdx])
		titleEn := getPointerStr(enRow[authorTitleIdx])
		titleRu := getPointerStr(ruRow[authorTitleIdx])
		author := types.Actor{
			Name:    authorName,
			TitleFi: titleFi,
			TitleEn: titleEn,
			TitleRu: titleRu,
		}
		authors = append(authors, author)
	}
	return authors
}

func getYear(year string) (*int, error) {
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		log.Printf("can not get a year %v", year)
		return nil, err
	}
	if yearInt == 9999 {
		return nil, nil
	}
	return &yearInt, nil
}

func getUses(usesFi, usesEn, usesRu string) ([]types.UseType, error) {
	useFiList := strings.Split(usesFi, ",")
	useEnList := strings.Split(usesEn, ",")
	useRuList := strings.Split(usesRu, ",")
	if len(useFiList) != len(useEnList) || len(useFiList) != len(useRuList) {
		return nil, fmt.Errorf(
			"unexpected en or ru uses: %v: %v: expect: %v",
			useEnList,
			useRuList,
			useFiList,
		)
	}

	uses := []types.UseType{}
	for i, useFi := range useFiList {
		useFi := strings.ToLower(strings.TrimSpace(useFi))
		if useFi == "" {
			continue
		}
		useEn := strings.ToLower(strings.TrimSpace(useEnList[i]))
		useRu := strings.ToLower(strings.TrimSpace(useRuList[i]))
		useType := types.UseType{NameFi: useFi, NameEn: useEn, NameRu: useRu}
		uses = append(uses, useType)
	}
	return uses, nil
}
