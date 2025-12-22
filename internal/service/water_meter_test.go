package service

import (
	"amr-data-bridge/config"
	"amr-data-bridge/internal/db"
	"amr-data-bridge/internal/dto"
	"amr-data-bridge/internal/repository"
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// mockWaterMeterRepository is a mock implementation of the WaterMeterRepository interface.
type mockWaterMeterRepository struct {
	GetWaterMetersFunc               func(ctx context.Context, arg db.GetWaterMetersParams) ([]db.GetWaterMetersRow, error)
	GetWaterMeterBySerialFunc        func(ctx context.Context, serialNumber string) (db.WaterMeter, error)
	UpdateWaterMeterActiveStatusFunc func(ctx context.Context, arg db.UpdateWaterMeterActiveStatusParams) error
}

func (m *mockWaterMeterRepository) GetWaterMeters(ctx context.Context, arg db.GetWaterMetersParams) ([]db.GetWaterMetersRow, error) {
	if m.GetWaterMetersFunc != nil {
		return m.GetWaterMetersFunc(ctx, arg)
	}
	return nil, errors.New("GetWaterMetersFunc not implemented")
}

func (m *mockWaterMeterRepository) GetWaterMeterBySerial(ctx context.Context, serialNumber string) (db.WaterMeter, error) {
	if m.GetWaterMeterBySerialFunc != nil {
		return m.GetWaterMeterBySerialFunc(ctx, serialNumber)
	}
	return db.WaterMeter{}, errors.New("GetWaterMeterBySerialFunc not implemented")
}

func (m *mockWaterMeterRepository) UpdateWaterMeterActiveStatus(ctx context.Context, arg db.UpdateWaterMeterActiveStatusParams) error {
	if m.UpdateWaterMeterActiveStatusFunc != nil {
		return m.UpdateWaterMeterActiveStatusFunc(ctx, arg)
	}
	return errors.New("UpdateWaterMeterActiveStatusFunc not implemented")
}

func boolp(b bool) *bool {
	return &b
}

func TestWaterMeterService_GetWaterMeters(t *testing.T) {
	// Mock preferences, as they are required by the mapper
	prefs := &config.Preferences{
		Export: config.ExportPreferences{
			WaterMeterFields: []string{
				"DevEUI",
				"SerialNumber",
				"IsActive",
				"LastSeen",
			},
		},
	}

	mockTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		req           dto.WaterMetersRequest
		setupMock     func(*mockWaterMeterRepository)
		want          []dto.WaterMeterResponse
		wantErr       bool
		expectedError string
	}{
		{
			name: "Success - fetch active meters",
			req: dto.WaterMetersRequest{
				Limit:  10,
				Active: boolp(true),
			},
			setupMock: func(m *mockWaterMeterRepository) {
				m.GetWaterMetersFunc = func(ctx context.Context, arg db.GetWaterMetersParams) ([]db.GetWaterMetersRow, error) {
					if arg.Limit != 10 || !arg.Active.Valid || !arg.Active.Bool {
						t.Fatalf("mock received unexpected params: %+v", arg)
					}
					return []db.GetWaterMetersRow{
						{DevEUI: "ABC", SerialNumber: "SN123", IsActive: true, LastSeen: pgtype.Timestamp{Time: mockTime, Valid: true}},
					}, nil
				}
			},
			want: []dto.WaterMeterResponse{
				{DevEUI: "ABC", SerialNumber: "SN123", IsActive: true, LastSeen: mockTime.Format(time.RFC3339)},
			},
			wantErr: false,
		},
		{
			name: "Repository returns an error",
			req:  dto.WaterMetersRequest{},
			setupMock: func(m *mockWaterMeterRepository) {
				m.GetWaterMetersFunc = func(ctx context.Context, arg db.GetWaterMetersParams) ([]db.GetWaterMetersRow, error) {
					return nil, errors.New("database error")
				}
			},
			want:          nil,
			wantErr:       true,
			expectedError: "database error",
		},
		{
			name: "Success - no meters found",
			req:  dto.WaterMetersRequest{},
			setupMock: func(m *mockWaterMeterRepository) {
				m.GetWaterMetersFunc = func(ctx context.Context, arg db.GetWaterMetersParams) ([]db.GetWaterMetersRow, error) {
					return []db.GetWaterMetersRow{}, nil
				}
			},
			want:    []dto.WaterMeterResponse{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockWaterMeterRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			// We need a non-nil repository.WaterMeterRepository to satisfy the constructor
			var repo repository.WaterMeterRepository = mockRepo
			s := NewWaterMeterService(repo, prefs)

			got, err := s.GetWaterMeters(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetWaterMeters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.expectedError {
				t.Errorf("GetWaterMeters() error = %v, expectedError %v", err, tt.expectedError)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWaterMeters() got = %v, want %v", got, tt.want)
			}
		})
	}
}
