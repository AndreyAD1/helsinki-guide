package services

import (
	"context"
	"errors"
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
		name            string
		fields          fields
		args            args
		foundBuildings  []types.Building
		repositoryError error
		want            []BuildingPreview
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
					NameFi:  utils.GetPointer("test name"),
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
					NameFi:  utils.GetPointer("test name 1"),
					Address: types.Address{StreetAddress: "test address 1"},
				},
				{
					NameFi:  utils.GetPointer("test name 2"),
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
			matchSpec := spec.AlikeAddressSpecIsEqual(
				tt.args.addressPrefix,
				tt.args.limit,
				tt.args.offset,
			)
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
		ctx     context.Context
		address string
	}
	tests := []struct {
		name                    string
		fields                  fields
		args                    args
		foundBuildings          []types.Building
		foundAuthors            []types.Actor
		repositoryBuildingError error
		repositoryActorError    error
		want                    []BuildingDTO
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
		{
			"building error",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{},
			[]types.Building{},
			[]types.Actor{},
			errors.New("building error"),
			errors.New("actor error"),
			[]BuildingDTO{},
		},
		{
			"actor error",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{},
			[]types.Building{{ID: 1}},
			[]types.Actor{},
			nil,
			errors.New("actor error"),
			[]BuildingDTO{},
		},
		{
			"one building - no authors",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{context.Background(), "test address"},
			[]types.Building{{ID: 1}},
			[]types.Actor{},
			nil,
			nil,
			[]BuildingDTO{{Address: "test address"}},
		},
		{
			"one building - one author",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{context.Background(), "test address"},
			[]types.Building{{ID: 1}},
			[]types.Actor{{Name: "author 1"}},
			nil,
			nil,
			[]BuildingDTO{{Address: "test address", Authors: &[]string{"author 1"}}},
		},
		{
			"one building - two authors",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{context.Background(), "test address"},
			[]types.Building{{ID: 1}},
			[]types.Actor{{Name: "author 1"}, {Name: "author 2"}},
			nil,
			nil,
			[]BuildingDTO{{Address: "test address", Authors: &[]string{"author 1", "author 2"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildingSpecFunc := spec.BuildingByAddressIsEqual(tt.args.address)
			tt.fields.buildingCollection.EXPECT().Query(
				tt.args.ctx,
				mock.MatchedBy(buildingSpecFunc),
			).Return(tt.foundBuildings, tt.repositoryBuildingError)

			if len(tt.foundBuildings) > 0 {
				authorSpec := spec.ActorByBuildingIsEqual(tt.foundBuildings[0].ID)
				tt.fields.actorCollection.On(
					"Query",
					tt.args.ctx,
					mock.MatchedBy(authorSpec),
				).Return(tt.foundAuthors, tt.repositoryActorError)
			}
			bs := BuildingService{
				buildingCollection: tt.fields.buildingCollection,
				actorCollection:    tt.fields.actorCollection,
			}
			got, err := bs.GetBuildingsByAddress(tt.args.ctx, tt.args.address)
			if tt.repositoryBuildingError != nil {
				require.ErrorIs(t, err, tt.repositoryBuildingError)
				return
			}
			if tt.repositoryActorError != nil {
				require.ErrorIs(t, err, tt.repositoryActorError)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestBuildingService_GetBuildingsByAddress_ManyBuildings(t *testing.T) {
	type fields struct {
		buildingCollection *repositories.BuildingRepository_mock
		actorCollection    *repositories.ActorRepository_mock
	}
	type args struct {
		ctx     context.Context
		address string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		foundBuildings []types.Building
		foundAuthors   []types.Actor
		want           []BuildingDTO
	}{
		{
			"two buildings - same authors",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{context.Background(), "test address"},
			[]types.Building{
				{ID: 1, NameFi: utils.GetPointer("test 1")},
				{ID: 2, NameFi: utils.GetPointer("test 2")},
			},
			[]types.Actor{{Name: "author 1"}, {Name: "author 2"}},
			[]BuildingDTO{
				{
					NameFi:  utils.GetPointer("test 1"),
					Address: "test address",
					Authors: &[]string{"author 1", "author 2"},
				},
				{
					NameFi:  utils.GetPointer("test 2"),
					Address: "test address",
					Authors: &[]string{"author 1", "author 2"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.buildingCollection.EXPECT().Query(
				tt.args.ctx,
				mock.MatchedBy(spec.BuildingByAddressIsEqual(tt.args.address)),
			).Return(tt.foundBuildings, nil)

			authorSpec0 := spec.ActorByBuildingIsEqual(tt.foundBuildings[0].ID) 
			authorSpec1 := spec.ActorByBuildingIsEqual(tt.foundBuildings[1].ID) 
			tt.fields.actorCollection.EXPECT().Query(
				tt.args.ctx,
				mock.MatchedBy(authorSpec0),
			).Return(tt.foundAuthors, nil).
				On("Query", tt.args.ctx, mock.MatchedBy(authorSpec1)).
				Return(tt.foundAuthors, nil)

			bs := BuildingService{
				buildingCollection: tt.fields.buildingCollection,
				actorCollection:    tt.fields.actorCollection,
			}
			got, err := bs.GetBuildingsByAddress(tt.args.ctx, tt.args.address)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
