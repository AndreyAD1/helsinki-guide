package services

import (
	"context"
	"strings"

	i "github.com/AndreyAD1/helsinki-guide/internal"
	"github.com/AndreyAD1/helsinki-guide/internal/infrastructure/repositories"
	s "github.com/AndreyAD1/helsinki-guide/internal/infrastructure/specifications"
)

type BuildingService struct {
	buildingCollection repositories.BuildingRepository
	actorCollection repositories.ActorRepository
}

func NewBuildingService(
	buildingCollection repositories.BuildingRepository,
	actorCollection repositories.ActorRepository,	
) BuildingService {
	return BuildingService{buildingCollection, actorCollection}
}

type BuildingPreview struct {
	Address string
	Name    string
}

type BuildingDTO struct {
	NameFi         *string `valueLanguage:"fi" nameFi:"Nimi" nameEn:"Name" nameRu:"Имя"`
	NameEn         *string `valueLanguage:"en" nameFi:"Nimi" nameEn:"Name" nameRu:"Имя"`
	NameRu         *string `valueLanguage:"ru" nameFi:"Nimi" nameEn:"Name" nameRu:"Имя"`
	Address        string  `valueLanguage:"all" nameFi:"Katuosoite" nameEn:"Address" nameRu:"Адрес"`
	CompletionYear *int    `valueLanguage:"all" nameFi:"Käyttöönottovuosi" nameEn:"Completion_year" nameRu:"Год_постройки"`
	Authors        *[]string `valueLanguage:"all" nameFi:"Suunnittelijat" nameEn:"Authors" nameRu:"Авторы"`
	HistoryFi      *string `valueLanguage:"fi" nameFi:"Rakennushistoria" nameEn:"Building_history" nameRu:"История_здания"`
	HistoryEn      *string `valueLanguage:"en" nameFi:"Rakennushistoria" nameEn:"Building_history" nameRu:"История_здания"`
	HistoryRu      *string `valueLanguage:"ru" nameFi:"Rakennushistoria" nameEn:"Building_history" nameRu:"История_здания"`
}

func NewBuildingDTO(b i.Building, authors []i.Actor, address string) BuildingDTO {
	var authorNames []string
	for _, author := range authors {
		authorNames = append(authorNames, author.Name)
	}
	var authorPtr *[]string
	if authorNames != nil {
		authorPtr = &authorNames
	}
	return BuildingDTO{
		b.NameFi,
		b.NameEn,
		b.NameRu,
		address,
		b.CompletionYear,
		authorPtr,
		b.HistoryFi,
		b.HistoryEn,
		b.HistoryRu,
	}
}

func (bs BuildingService) GetBuildingPreviews(
	ctx context.Context,
	addressPrefix string,
	limit,
	offset int,
) ([]BuildingPreview, error) {
	addressPrefix = strings.TrimLeft(addressPrefix, " ")
	spec := s.NewBuildingSpecificationByAlikeAddress(addressPrefix, limit, offset)
	buildings, err := bs.buildingCollection.Query(ctx, spec)
	if err != nil {
		return nil, err
	}

	previews := make([]BuildingPreview, len(buildings))
	for i, building := range buildings {
		name := ""
		if building.NameFi != nil {
			name = *building.NameFi
		}
		previews[i] = BuildingPreview{building.Address.StreetAddress, name}
	}
	return previews, nil
}

func (bs BuildingService) GetBuildingsByAddress(
	ctx context.Context,
	address string,
) ([]BuildingDTO, error) {
	address = strings.TrimSpace(address)
	spec := s.NewBuildingSpecificationByAddress(address)
	buildings, err := bs.buildingCollection.Query(ctx, spec)
	if err != nil {
		return nil, err
	}
	buildingsDto := make([]BuildingDTO, len(buildings))
	for i, building := range buildings {
		spec := s.NewActorSpecificationByBuilding(building.ID)
		authors, err := bs.actorCollection.Query(ctx, spec)
		if err != nil {
			return nil, err
		}
		buildingsDto[i] = NewBuildingDTO(building, authors, address)
	}
	return buildingsDto, nil
}
