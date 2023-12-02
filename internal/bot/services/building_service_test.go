package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	spec "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/specifications"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/types"
	"github.com/AndreyAD1/helsinki-guide/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBuildingService_GetBuildingPreviews(t *testing.T) {
	type fields struct {
		buildingCollection *repositories.BuildingRepository_mock
		actorCollection    *repositories.ActorRepository_mock
	}
	type args struct {
		ctx           context.Context
		addressPrefix string
		limit         int
		offset        int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		foundBuildings []types.Building
		repositoryError error
		want    []BuildingPreview
	}{
		{
			"no previews - no arguments",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{},
			[]types.Building{},
			nil,
			[]BuildingPreview{},
		},
		{
			"no previews",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{context.Background(), "test", 5, 10},
			[]types.Building{},
			nil,
			[]BuildingPreview{},
		},
		{
			"repository error",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{},
			[]types.Building{},
			errors.New("test error"),
			[]BuildingPreview{},
		},
		{
			"one preview",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{context.Background(), "test", 5, 10},
			[]types.Building{
				{
					NameFi: utils.GetPointer("test name"), 
					Address: types.Address{StreetAddress: "test address"},
				},
			},
			nil,
			[]BuildingPreview{{"test address", "test name"}},
		},
		{
			"two previews",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{context.Background(), "test", 5, 10},
			[]types.Building{
				{
					NameFi: utils.GetPointer("test name 1"), 
					Address: types.Address{StreetAddress: "test address 1"},
				},
				{
					NameFi: utils.GetPointer("test name 2"), 
					Address: types.Address{StreetAddress: "test address 2"},
				},
			},
			nil,
			[]BuildingPreview{
				{"test address 1", "test name 1"}, 
				{"test address 2", "test name 2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchSpec := func(s *spec.BuildingSpecificationByAlikeAddress) bool {
				addressMatch := s.AddressPrefix == tt.args.addressPrefix
				limitMatch := s.Limit == tt.args.limit
				offsetMatch := s.Offset == tt.args.offset
				return addressMatch && limitMatch && offsetMatch
			}
			tt.fields.buildingCollection.EXPECT().Query(
				tt.args.ctx, 
				mock.MatchedBy(matchSpec),
			).Return(tt.foundBuildings, tt.repositoryError)
			bs := BuildingService{
				buildingCollection: tt.fields.buildingCollection,
				actorCollection:    tt.fields.actorCollection,
			}
			got, err := bs.GetBuildingPreviews(
				tt.args.ctx, 
				tt.args.addressPrefix, 
				tt.args.limit, 
				tt.args.offset,
			)
			if tt.repositoryError == nil {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			} else {
				require.ErrorIs(t, err, tt.repositoryError)
				require.Nil(t, got)
			}
			
		})
	}
}

func TestBuildingService_GetBuildingsByAddress(t *testing.T) {
	type fields struct {
		buildingCollection *repositories.BuildingRepository_mock
		actorCollection    *repositories.ActorRepository_mock
	}
	type args struct {
		ctx           context.Context
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		foundBuildings []types.Building
		foundAuthors []types.Actor
		repositoryBuildingError error
		repositoryActorError error
		want    []BuildingDTO
	}{
		{
			"no address",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{},
			[]types.Building{},
			[]types.Actor{},
			nil,
			nil,
			[]BuildingDTO{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildingSpec := func(s *spec.BuildingSpecificationByAddress) bool {
				return s.Address == tt.args.address
			}

			tt.fields.buildingCollection.EXPECT().Query(
				tt.args.ctx, 
				mock.MatchedBy(buildingSpec),
			).Return(tt.foundBuildings, tt.repositoryBuildingError)
			for _, building := range tt.foundBuildings {
				authorSpec := func(s spec.ActorSpecificationByBuilding) bool {
					return s.BuildingID == building.ID
				}
				tt.fields.actorCollection.EXPECT().Query(
					tt.args.ctx,
					authorSpec,
				).Return(tt.foundAuthors, tt.repositoryActorError)
			}

			bs := BuildingService{
				buildingCollection: tt.fields.buildingCollection,
				actorCollection:    tt.fields.actorCollection,
			}
			got, err := bs.GetBuildingsByAddress(tt.args.ctx, tt.args.address)
			if tt.repositoryBuildingError != nil {
				require.Error(t, err)
				return
			}
			if tt.repositoryActorError != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			fmt.Println(tt.want, got)
			require.Equal(t, tt.want, got)
		})
	}
}
