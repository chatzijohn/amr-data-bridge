package service

import (
	"amr-data-bridge/internal/dto"
	"amr-data-bridge/internal/repository"
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

// mockWaterSupplyRepository is a mock implementation of the WaterSupplyRepository interface.
type mockWaterSupplyRepository struct {
	ImportWaterSuppliesFunc func(ctx context.Context, req []dto.WaterSupplyRequest) ([]dto.WaterSupplyResponse, error)
}

func (m *mockWaterSupplyRepository) ImportWaterSupplies(ctx context.Context, req []dto.WaterSupplyRequest) ([]dto.WaterSupplyResponse, error) {
	if m.ImportWaterSuppliesFunc != nil {
		return m.ImportWaterSuppliesFunc(ctx, req)
	}
	return nil, errors.New("ImportWaterSuppliesFunc not implemented")
}

func TestWaterSupplyService_ImportWaterSupplies(t *testing.T) {
	mockTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC).Format(time.RFC3339)

	tests := []struct {
		name          string
		req           []dto.WaterSupplyRequest
		setupMock     func(*mockWaterSupplyRepository)
		want          []dto.WaterSupplyResponse
		wantErr       bool
		expectedError string
	}{
		{
			name: "Success - import two supplies",
			req: []dto.WaterSupplyRequest{
				{SupplyNumber: "S1", Latitude: 1.0, Longitude: 1.0},
				{SupplyNumber: "S2", Latitude: 2.0, Longitude: 2.0},
			},
			setupMock: func(m *mockWaterSupplyRepository) {
				m.ImportWaterSuppliesFunc = func(ctx context.Context, req []dto.WaterSupplyRequest) ([]dto.WaterSupplyResponse, error) {
					if len(req) != 2 {
						t.Fatalf("mock expected 2 requests, got %d", len(req))
					}
					return []dto.WaterSupplyResponse{
						{ID: 1, SupplyNumber: "S1", CreatedAt: mockTime, UpdatedAt: mockTime},
						{ID: 2, SupplyNumber: "S2", CreatedAt: mockTime, UpdatedAt: mockTime},
					}, nil
				}
			},
			want: []dto.WaterSupplyResponse{
				{ID: 1, SupplyNumber: "S1", CreatedAt: mockTime, UpdatedAt: mockTime},
				{ID: 2, SupplyNumber: "S2", CreatedAt: mockTime, UpdatedAt: mockTime},
			},
			wantErr: false,
		},
		{
			name: "Repository returns an error",
			req: []dto.WaterSupplyRequest{
				{SupplyNumber: "S1", Latitude: 1.0, Longitude: 1.0},
			},
			setupMock: func(m *mockWaterSupplyRepository) {
				m.ImportWaterSuppliesFunc = func(ctx context.Context, req []dto.WaterSupplyRequest) ([]dto.WaterSupplyResponse, error) {
					return nil, errors.New("repository transaction failed")
				}
			},
			want:          nil,
			wantErr:       true,
			expectedError: "repository transaction failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockWaterSupplyRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			var repo repository.WaterSupplyRepository = mockRepo
			s := NewWaterSupplyService(repo)

			got, err := s.ImportWaterSupplies(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("ImportWaterSupplies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.expectedError {
				t.Errorf("ImportWaterSupplies() error = %v, expectedError %v", err, tt.expectedError)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImportWaterSupplies() got = %v, want %v", got, tt.want)
			}
		})
	}
}
