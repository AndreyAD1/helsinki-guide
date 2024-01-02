package services

import (
	"context"
	"errors"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_SetLanguage(t *testing.T) {
	type fields struct {
		userCollection *repositories.UserRepository_mock
	}
	type args struct {
		ctx      context.Context
		userID   int64
		language Language
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		repositoryError error
	}{
		{
			"success",
			fields{
				repositories.NewUserRepository_mock(t),
			},
			args{
				context.Background(),
				123,
				Finnish,
			},
			nil,
		},
		{
			"error",
			fields{
				repositories.NewUserRepository_mock(t),
			},
			args{
				context.Background(),
				123,
				Finnish,
			},
			errors.New("some DB error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedUser := repositories.User{
				TelegramID:        tt.args.userID,
				PreferredLanguage: string(tt.args.language),
			}
			tt.fields.userCollection.EXPECT().
				AddOrUpdate(tt.args.ctx, expectedUser).
				Return(nil, tt.repositoryError)
			s := UserService{userCollection: tt.fields.userCollection}
			err := s.SetLanguage(tt.args.ctx, tt.args.userID, tt.args.language)
			require.ErrorIs(t, err, tt.repositoryError)
		})
	}
}

func TestUserService_GetPreferredLanguage(t *testing.T) {
	type fields struct {
		userCollection *repositories.UserRepository_mock
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		foundUsers      []repositories.User
		repositoryError error
		want            *Language
		expectedError   error
	}{
		{
			"DB error",
			fields{repositories.NewUserRepository_mock(t)},
			args{context.Background(), 123},
			nil,
			repositories.ErrNotImplemented,
			nil,
			repositories.ErrNotImplemented,
		},
		{
			"no users",
			fields{repositories.NewUserRepository_mock(t)},
			args{context.Background(), 123},
			nil,
			nil,
			nil,
			nil,
		},
		{
			"no preferred language",
			fields{repositories.NewUserRepository_mock(t)},
			args{context.Background(), 123},
			[]repositories.User{{}},
			nil,
			nil,
			nil,
		},
		{
			"unexpected preferred language",
			fields{repositories.NewUserRepository_mock(t)},
			args{context.Background(), 123},
			[]repositories.User{{PreferredLanguage: "unexpected"}},
			nil,
			nil,
			nil,
		},
		{
			"a valid preferred language",
			fields{repositories.NewUserRepository_mock(t)},
			args{context.Background(), 123},
			[]repositories.User{{PreferredLanguage: "fi"}},
			nil,
			&Finnish,
			nil,
		},
		{
			"severa users",
			fields{repositories.NewUserRepository_mock(t)},
			args{context.Background(), 123},
			[]repositories.User{{PreferredLanguage: "en"}, {PreferredLanguage: "fi"}},
			nil,
			&English,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchSpec := repositories.UserByIDIsEqual(tt.args.userID)
			tt.fields.userCollection.EXPECT().Query(
				tt.args.ctx,
				mock.MatchedBy(matchSpec),
			).Return(tt.foundUsers, tt.repositoryError)
			us := UserService{userCollection: tt.fields.userCollection}
			got, err := us.GetPreferredLanguage(tt.args.ctx, tt.args.userID)
			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, tt.want, got)
		})
	}
}
