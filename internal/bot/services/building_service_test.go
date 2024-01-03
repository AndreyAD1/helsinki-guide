package services

import (
	"context"
	"errors"
	"testing"

	r "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	"github.com/AndreyAD1/helsinki-guide/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBuildingService_GetBuildingPreviews(t *testing.T) {
	type fields struct {
		buildingCollection *r.BuildingRepository_mock
		actorCollection    *r.ActorRepository_mock
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
		foundBuildings  []r.Building
		repositoryError error
		want            []BuildingPreview
	}{
		{
			"no previews - no arguments",
			fields{
				r.NewBuildingRepository_mock(t),
				r.NewActorRepository_mock(t),
			},
			args{},
			[]r.Building{},
			nil,
			[]BuildingPreview{},
		},
		{
			"no previews",
			fields{
				r.NewBuildingRepository_mock(t),
				r.NewActorRepository_mock(t),
			},
			args{context.Background(), "test", 5, 10},
			[]r.Building{},
			nil,
			[]BuildingPreview{},
		},
		{
			"repository error",
			fields{
				r.NewBuildingRepository_mock(t),
				r.NewActorRepository_mock(t),
			},
			args{},
			[]r.Building{},
			errors.New("test error"),
			[]BuildingPreview{},
		},
		{
			"one preview",
			fields{
				r.NewBuildingRepository_mock(t),
				r.NewActorRepository_mock(t),
			},
			args{context.Background(), "test", 5, 10},
			[]r.Building{
				{
					NameFi:  utils.GetPointer("test name"),
					Address: r.Address{StreetAddress: "test address"},
				},
			},
			nil,
			[]BuildingPreview{{"test address", "test name"}},
		},
		{
			"two previews",
			fields{
				r.NewBuildingRepository_mock(t),
				r.NewActorRepository_mock(t),
			},
			args{context.Background(), "test", 5, 10},
			[]r.Building{
				{
					NameFi:  utils.GetPointer("test name 1"),
					Address: r.Address{StreetAddress: "test address 1"},
				},
				{
					NameFi:  utils.GetPointer("test name 2"),
					Address: r.Address{StreetAddress: "test address 2"},
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
			matchSpec := r.AlikeAddressSpecIsEqual(
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
		buildingCollection *r.BuildingRepository_mock
		actorCollection    *r.ActorRepository_mock
	}
	type args struct {
		ctx     context.Context
		address string
	}
	tests := []struct {
		name                    string
		fields                  fields
		args                    args
		foundBuildings          []r.Building
		foundAuthors            []r.Actor
		repositoryBuildingError error
		repositoryActorError    error
		want                    []BuildingDTO
	}{
		{
			"no address",
			fields{
				r.NewBuildingRepository_mock(t),
				r.NewActorRepository_mock(t),
			},
			args{},
			[]r.Building{},
			[]r.Actor{},
			nil,
			nil,
			[]BuildingDTO{},
		},
		{
			"building error",
			fields{
				r.NewBuildingRepository_mock(t),
				r.NewActorRepository_mock(t),
			},
			args{},
			[]r.Building{},
			[]r.Actor{},
			errors.New("building error"),
			errors.New("actor error"),
			[]BuildingDTO{},
		},
		{
			"actor error",
			fields{
				r.NewBuildingRepository_mock(t),
				r.NewActorRepository_mock(t),
			},
			args{},
			[]r.Building{{ID: 1}},
			[]r.Actor{},
			nil,
			errors.New("actor error"),
			[]BuildingDTO{},
		},
		{
			"one building - no authors",
			fields{
				r.NewBuildingRepository_mock(t),
				r.NewActorRepository_mock(t),
			},
			args{context.Background(), "test address"},
			[]r.Building{{ID: 1, Address: r.Address{StreetAddress: "test address"}}},
			[]r.Actor{},
			nil,
			nil,
			[]BuildingDTO{{Address: "test address"}},
		},
		{
			"one building - one author",
			fields{
				r.NewBuildingRepository_mock(t),
				r.NewActorRepository_mock(t),
			},
			args{context.Background(), "test address"},
			[]r.Building{{ID: 1, Address: r.Address{StreetAddress: "test address"}}},
			[]r.Actor{{Name: "author 1"}},
			nil,
			nil,
			[]BuildingDTO{{Address: "test address", Authors: &[]string{"author 1"}}},
		},
		{
			"one building - two authors",
			fields{
				r.NewBuildingRepository_mock(t),
				r.NewActorRepository_mock(t),
			},
			args{context.Background(), "test address"},
			[]r.Building{{ID: 1, Address: r.Address{StreetAddress: "test address"}}},
			[]r.Actor{{Name: "author 1"}, {Name: "author 2"}},
			nil,
			nil,
			[]BuildingDTO{{Address: "test address", Authors: &[]string{"author 1", "author 2"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildingSpecFunc := r.BuildingByAddressIsEqual(tt.args.address)
			tt.fields.buildingCollection.EXPECT().Query(
				tt.args.ctx,
				mock.MatchedBy(buildingSpecFunc),
			).Return(tt.foundBuildings, tt.repositoryBuildingError)

			if len(tt.foundBuildings) > 0 {
				authorSpec := r.ActorByBuildingIsEqual(tt.foundBuildings[0].ID)
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
		buildingCollection *r.BuildingRepository_mock
		actorCollection    *r.ActorRepository_mock
	}
	type args struct {
		ctx     context.Context
		address string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		foundBuildings []r.Building
		foundAuthors   []r.Actor
		want           []BuildingDTO
	}{
		{
			"two buildings - same authors",
			fields{
				r.NewBuildingRepository_mock(t),
				r.NewActorRepository_mock(t),
			},
			args{context.Background(), "test address"},
			[]r.Building{
				{ID: 1, NameFi: utils.GetPointer("test 1")},
				{ID: 2, NameFi: utils.GetPointer("test 2")},
			},
			[]r.Actor{{Name: "author 1"}, {Name: "author 2"}},
			[]BuildingDTO{
				{
					NameFi:  utils.GetPointer("test 1"),
					Address: "",
					Authors: &[]string{"author 1", "author 2"},
				},
				{
					NameFi:  utils.GetPointer("test 2"),
					Address: "",
					Authors: &[]string{"author 1", "author 2"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.buildingCollection.EXPECT().Query(
				tt.args.ctx,
				mock.MatchedBy(r.BuildingByAddressIsEqual(tt.args.address)),
			).Return(tt.foundBuildings, nil)

			authorSpec0 := r.ActorByBuildingIsEqual(tt.foundBuildings[0].ID)
			authorSpec1 := r.ActorByBuildingIsEqual(tt.foundBuildings[1].ID)
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
