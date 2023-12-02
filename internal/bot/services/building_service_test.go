package services

import (
	"context"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	spec "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/specifications"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/types"
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
		want    []BuildingPreview
		wantErr bool
	}{
		{
			"no previews",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{},
			[]BuildingPreview{},
			false,
		},
		{
			"no previews",
			fields{
				repositories.NewBuildingRepository_mock(t),
				repositories.NewActorRepository_mock(t),
			},
			args{context.Background(), "test", 5, 10},
			[]BuildingPreview{},
			false,
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
			).Return([]types.Building{}, nil)
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
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}
