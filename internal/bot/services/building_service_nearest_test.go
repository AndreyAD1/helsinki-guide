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

func TestBuildingService_GetNearestPreviews(t *testing.T) {
	type fields struct {
		buildingCollection *repositories.BuildingRepository_mock
		actorCollection    *repositories.ActorRepository_mock
	}
	type args struct {
		ctx           context.Context
		distance int
		latitude float64
		longitude float64
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
			args{context.Background(), 100, 0.0, 0.0, 5, 10},
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
			args{context.Background(), 100, 0.0, 0.0, 5, 10},
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
			args{context.Background(), 100, 0.0, 0.0, 5, 10},
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
			matchSpecFunc := spec.NearestSpecIsEqual(
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
			got, err := bs.GetNearestBuildingPreviews(
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