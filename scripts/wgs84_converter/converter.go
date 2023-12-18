package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/specifications"
	"github.com/jackc/pgx/v5/pgxpool"
)

func run() error {
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf(
			"unable to create a connection pool: DB URL '%s': %w",
			databaseURL,
			err,
		)
	}
	if err := dbpool.Ping(ctx); err != nil {
		logMsg := fmt.Sprintf(
			"unable to connect to the DB '%v'",
			databaseURL,
		)
		log.Println(logMsg)
		return fmt.Errorf("%v: %w", logMsg, err)
	}
	buildingRepo := repositories.NewBuildingRepo(dbpool)
	converterClient := NewCoordinateConverterClient(converterURL)
	spec := specifications.NewBuildingSpecificationByAlikeAddress("", 500, 0)
	buildings, err := buildingRepo.Query(ctx, spec)
	if err != nil {
		return err
	}
	for _, building := range buildings {
		latitude := building.Latitude_ETRSGK25
		longitude := building.Longitude_ETRSGK25
		latitudeWGS84, longitudeWGS84, err := converterClient.Convert(
			ctx,
			latitude,
			longitude,
		)
		building.Latitude_WGS84 = latitudeWGS84
		building.Longitude_ETRSGK25 = longitudeWGS84
		if _, err = buildingRepo.Update(ctx, building); err != nil {
			log.Println("can not update a building '%v'", building.ID)
		}
	}

	return nil
}
