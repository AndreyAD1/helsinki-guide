package services

import (
	"context"
	"errors"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	"github.com/AndreyAD1/helsinki-guide/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBuildingService_GetNearestPreviews(t *testing.T) {
	type fields struct {
		buildingCollection *repositories.BuildingRepository_mock
		actorCollection    *repositories.ActorRepository_mock
	}
	type args struct {
		ctx       context.Context
		distance  int
		latitude  float64
		longitude float64
		limit     int
		offset    int
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		foundBuildings  []repositories.Building
		repositoryError error
		want            []BuildingDTO
	}{
		{
			"no previews - no arguments",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{},
			[]repositories.Building{},
			nil,
			[]BuildingDTO{},
		},
		{
			"no previews",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{context.Background(), 100, 0.0, 0.0, 5, 10},
			[]repositories.Building{},
			nil,
			[]BuildingDTO{},
		},
		{
			"repository error",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{},
			[]repositories.Building{},
			errors.New("test error"),
			[]BuildingDTO{},
		},
		{
			"one preview",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{context.Background(), 100, 0.0, 0.0, 5, 10},
			[]repositories.Building{
				{
					NameFi:  utils.GetPointer("test name"),
					Address: repositories.Address{StreetAddress: "test address"},
				},
			},
			nil,
			[]BuildingDTO{{Address: "test address", NameFi: utils.GetPointer("test name")}},
		},
		{
			"two previews",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{context.Background(), 100, 0.0, 0.0, 5, 10},
			[]repositories.Building{
				{
					NameFi:  utils.GetPointer("test name 1"),
					Address: repositories.Address{StreetAddress: "test address 1"},
				},
				{
					NameFi:  utils.GetPointer("test name 2"),
					Address: repositories.Address{StreetAddress: "test address 2"},
				},
			},
			nil,
			[]BuildingDTO{
				{Address: "test address 1", NameFi: utils.GetPointer("test name 1")},
				{Address: "test address 2", NameFi: utils.GetPointer("test name 2")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchSpecFunc := repositories.NearestSpecIsEqual(
				tt.args.distance,
				tt.args.latitude,
				tt.args.longitude,
				tt.args.limit,
				tt.args.offset,
			)
			tt.fields.buildingCollection.EXPECT().Query(
				tt.args.ctx,
				mock.MatchedBy(matchSpecFunc),
			).Return(tt.foundBuildings, tt.repositoryError)
			bs := BuildingService{
				buildingCollection: tt.fields.buildingCollection,
				actorCollection:    tt.fields.actorCollection,
			}
			got, err := bs.GetNearestBuildings(
				tt.args.ctx,
				tt.args.distance,
				tt.args.latitude,
				tt.args.longitude,
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
