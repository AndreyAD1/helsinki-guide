package services

import (
	"context"
	"reflect"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
)

func TestBuildingService_GetBuildingPreviews(t *testing.T) {
	type fields struct {
		buildingCollection repositories.BuildingRepository
		actorCollection    repositories.ActorRepository
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := BuildingService{
				buildingCollection: tt.fields.buildingCollection,
				actorCollection:    tt.fields.actorCollection,
			}
			got, err := bs.GetBuildingPreviews(tt.args.ctx, tt.args.addressPrefix, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildingService.GetBuildingPreviews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildingService.GetBuildingPreviews() = %v, want %v", got, tt.want)
			}
		})
	}
}
