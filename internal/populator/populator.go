package populator

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"

	"github.com/AndreyAD1/helsinki-guide/internal"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/configuration"
	"github.com/AndreyAD1/helsinki-guide/internal/infrastructure/repositories"
)

var (
	codeIdx = 1
	nameIdx = 2
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

func (p *Populator) Run(ctx context.Context, sheetName, fiFilename, enFilename, ruFilename string) error {
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

	for fiRows.Next() {
		enRows.Next()
		ruRows.Next()

		fiRow, err := fiRows.Columns()
		enRow, err := enRows.Columns()
		ruRow, err := ruRows.Columns()

		address, err := getAddress(fiRow)
		if err != nil {
			return err
		}

		authorIDs := []int64{}
		for _, author := range getAuthors(fiRow, enRow, ruRow) {
			savedAuthor, err := p.actorRepo.Add(ctx, author)
			if err != nil && err != repositories.ErrDuplicate {
				return err
			}
			authorIDs = append(authorIDs, savedAuthor.ID)
		}

		building := internal.Building{
			Code: &fiRow[codeIdx],
			Address: address,
			NameFi: &fiRow[nameIdx],
			NameEn: &enRow[nameIdx],
			NameRu: &ruRow[nameIdx],
			AuthorIDs: authorIDs,
		}
		_, err = p.buildingRepo.Add(ctx, building)
		if err != nil {
			return err
		}
	}
	return nil
}

func getAddress(row []string) (internal.Address, error) {
	
	return internal.Address{}, nil
}

func getAuthors(fiRow, enRow, ruRow []string) []internal.Actor {
	return []internal.Actor{}
}