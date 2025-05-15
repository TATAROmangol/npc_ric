package delete

import (
	"context"
	"errors"
	"forms/internal/service/delete/mocks"
	"forms/pkg/logger"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestDeleteService_DeleteInstitution(t *testing.T) {
	repo := mocks.NewMockDeleteRepo(gomock.NewController(t))

	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	type MockBehavior func(institutionId int)

	tests := []struct {
		name    string
		mockBehavior    MockBehavior
		institutionId int
		wantErr bool
	}{
		{
			name:    "valid institutionId",
			mockBehavior: func(institutionId int) {
				repo.EXPECT().DeleteInstitution(gomock.Any(), institutionId).Return(nil)
			},
			institutionId: 1,
			wantErr: false,
		},
		{
			name:    "invalid institutionId",
			mockBehavior: func(institutionId int) {
				repo.EXPECT().DeleteInstitution(gomock.Any(), institutionId).Return(errors.New("invalid institutionId"))
			},
			institutionId: 0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.institutionId)
			ds := &DeleteService{
				DeleteRepo: repo,
			}
			if err := ds.DeleteInstitution(ctx, tt.institutionId); (err != nil) != tt.wantErr {
				t.Errorf("DeleteService.DeleteInstitution() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteService_DeleteMentor(t *testing.T) {
	repo := mocks.NewMockDeleteRepo(gomock.NewController(t))

	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	type MockBehavior func(mentorId int)

	tests := []struct {
		name    string
		mockBehavior    MockBehavior
		mentorId int
		wantErr bool
	}{
		{
			name:    "valid mentorId",
			mockBehavior: func(mentorId int) {
				repo.EXPECT().DeleteMentor(gomock.Any(), mentorId).Return(nil)
			},
			mentorId: 1,
			wantErr: false,
		},
		{
			name:    "invalid mentorId",
			mockBehavior: func(mentorId int) {
				repo.EXPECT().DeleteMentor(gomock.Any(), mentorId).Return(errors.New("invalid mentorId"))
			},
			mentorId: 0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.mentorId)
			ds := &DeleteService{
				DeleteRepo: repo,
			}
			if err := ds.DeleteMentor(ctx, tt.mentorId); (err != nil) != tt.wantErr {
				t.Errorf("DeleteService.DeleteMentor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
