package services

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	s "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/specifications"
	i "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/types"
)

type BuildingService struct {
	buildingCollection repositories.BuildingRepository
	actorCollection    repositories.ActorRepository
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
	NameFi            *string   `valueLanguage:"fi" nameFi:"Nimi" nameEn:"Name" nameRu:"Имя"`
	NameEn            *string   `valueLanguage:"en" nameFi:"Nimi" nameEn:"Name" nameRu:"Имя"`
	NameRu            *string   `valueLanguage:"ru" nameFi:"Nimi" nameEn:"Name" nameRu:"Имя"`
	Address           string    `valueLanguage:"all" nameFi:"Katuosoite" nameEn:"Address" nameRu:"Адрес"`
	DescriptionFi     *string   `valueLanguage:"fi" nameFi:"Kerrosluku" nameEn:"Description" nameRu:"Описание"`
	DescriptionEn     *string   `valueLanguage:"en" nameFi:"Kerrosluku" nameEn:"Description" nameRu:"Описание"`
	DescriptionRu     *string   `valueLanguage:"ru" nameFi:"Kerrosluku" nameEn:"Description" nameRu:"Описание"`
	CompletionYear    *int      `valueLanguage:"all" nameFi:"Käyttöönottovuosi" nameEn:"Completion_year" nameRu:"Год_постройки"`
	Authors           *[]string `valueLanguage:"all" nameFi:"Suunnittelijat" nameEn:"Authors" nameRu:"Авторы"`
	HistoryFi         *string   `valueLanguage:"fi" nameFi:"Rakennushistoria" nameEn:"Building_history" nameRu:"История_здания"`
	HistoryEn         *string   `valueLanguage:"en" nameFi:"Rakennushistoria" nameEn:"Building_history" nameRu:"История_здания"`
	HistoryRu         *string   `valueLanguage:"ru" nameFi:"Rakennushistoria" nameEn:"Building_history" nameRu:"История_здания"`
	NotableFeaturesFi *string   `valueLanguage:"fi" nameFi:"Huomattavia_ominaisuuksia" nameEn:"Notable_features" nameRu:"Примечательные_особенности"`
	NotableFeaturesEn *string   `valueLanguage:"en" nameFi:"Huomattavia_ominaisuuksia" nameEn:"Notable_features" nameRu:"Примечательные_особенности"`
	NotableFeaturesRu *string   `valueLanguage:"ru" nameFi:"Huomattavia_ominaisuuksia" nameEn:"Notable_features" nameRu:"Примечательные_особенности"`
	FacadesFi         *string   `valueLanguage:"fi" nameFi:"Julkisivut" nameEn:"Facades" nameRu:"Фасады"`
	FacadesEn         *string   `valueLanguage:"en" nameFi:"Julkisivut" nameEn:"Facades" nameRu:"Фасады"`
	FacadesRu         *string   `valueLanguage:"ru" nameFi:"Julkisivut" nameEn:"Facades" nameRu:"Фасады"`
	DetailsFi         *string   `valueLanguage:"fi" nameFi:"Erityispiirteet" nameEn:"Interesting_details" nameRu:"Интересные_детали"`
	DetailsEn         *string   `valueLanguage:"en" nameFi:"Erityispiirteet" nameEn:"Interesting_details" nameRu:"Интересные_детали"`
	DetailsRu         *string   `valueLanguage:"ru" nameFi:"Erityispiirteet" nameEn:"Interesting_details" nameRu:"Интересные_детали"`
	SurroundingsFi    *string   `valueLanguage:"fi" nameFi:"Ymparistonkuvaus" nameEn:"Surroundings" nameRu:"Окрестности"`
	SurroundingsEn    *string   `valueLanguage:"en" nameFi:"Ymparistonkuvaus" nameEn:"Surroundings" nameRu:"Окрестности"`
	SurroundingsRu    *string   `valueLanguage:"ru" nameFi:"Ymparistonkuvaus" nameEn:"Surroundings" nameRu:"Окрестности"`
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
		NameFi:            b.NameFi,
		NameEn:            b.NameEn,
		NameRu:            b.NameRu,
		Address:           address,
		DescriptionFi:     b.FloorDescriptionFi,
		DescriptionEn:     b.FloorDescriptionEn,
		DescriptionRu:     b.FloorDescriptionRu,
		CompletionYear:    b.CompletionYear,
		Authors:           authorPtr,
		HistoryFi:         b.HistoryFi,
		HistoryEn:         b.HistoryEn,
		HistoryRu:         b.HistoryRu,
		NotableFeaturesFi: b.ReasoningFi,
		NotableFeaturesEn: b.ReasoningEn,
		NotableFeaturesRu: b.ReasoningRu,
		FacadesFi:         b.FacadesFi,
		FacadesEn:         b.FacadesEn,
		FacadesRu:         b.FacadesRu,
		DetailsFi:         b.SpecialFeaturesFi,
		DetailsEn:         b.SpecialFeaturesEn,
		DetailsRu:         b.SpecialFeaturesRu,
		SurroundingsFi:    b.SurroundingsFi,
		SurroundingsEn:    b.SurroundingsEn,
		SurroundingsRu:    b.SpecialFeaturesRu,
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
	slog.DebugContext(ctx, fmt.Sprintf("buildings: %v", buildings))
	buildingsDto := make([]BuildingDTO, len(buildings))
	for i, building := range buildings {
		spec := s.NewActorSpecificationByBuilding(building.ID)
		authors, err := bs.actorCollection.Query(ctx, spec)
		if err != nil {
			return nil, err
		}
		buildingsDto[i] = NewBuildingDTO(building, authors, address)
		slog.DebugContext(ctx, fmt.Sprintf("authors: %v", authors))
	}

	return buildingsDto, nil
}
