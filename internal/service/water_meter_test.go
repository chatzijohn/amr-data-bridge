package service_test

import (
	"amr-data-bridge/internal/db"
	"amr-data-bridge/internal/service"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type mockStore struct{}

func mockWaterMeter(id int32, devEUI string) db.WaterMeter {
	now := pgtype.Timestamp{Time: time.Now(), Valid: true}
	return db.WaterMeter{
		ID:                id,
		DevEUI:            devEUI,
		SerialNumber:      "SN12345",
		BrandName:         "AquaSmart",
		LtPerPulse:        10,
		CurrentReading:    pgtype.Int4{Int32: 500, Valid: true},
		IsActive:          true,
		AlarmStatus:       false,
		NoFlow:            false,
		DeviceHandshake:   pgtype.Int4{Int32: 1, Valid: true},
		DeviceLogging:     pgtype.Int4{Int32: 1, Valid: true},
		ServerHandshake:   pgtype.Int4{Int32: 1, Valid: true},
		ServerLogging:     pgtype.Int4{Int32: 1, Valid: true},
		InactivityTimeout: 60,
		LastSeen:          now,
		Rssi:              pgtype.Int4{Int32: -70, Valid: true},
		Snr:               pgtype.Float8{Float64: 8.5, Valid: true},
		SpreadingFactor:   pgtype.Int4{Int32: 7, Valid: true},
		GatewayId:         pgtype.Text{String: "gateway-1", Valid: true},
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func (m *mockStore) GetWaterMeters(ctx context.Context, _ db.GetWaterMetersParams) ([]db.WaterMeter, error) {
	return []db.WaterMeter{mockWaterMeter(1, "ABC123"), mockWaterMeter(2, "A3CD2E")}, nil
}

func TestGetWaterMeters(t *testing.T) {
	params := db.GetWaterMetersParams{
		Limit:  1,
		Active: false, // test default/no filter
	}

	svc := service.NewWaterMeterService(&mockStore{})
	meters, err := svc.GetWaterMeters(context.Background(), params)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if params.Limit == 1 && len(meters) > 1 {
		t.Errorf("expected 1 meter, got %d", len(meters))
	}

	if len(meters) != 2 {
		t.Errorf("expected 2 meters, got %d", len(meters))
	}

	if meters[0].DevEUI != "ABC123" {
		t.Errorf("expected first meter name to be 'Meter A', got %s", meters[0].DevEUI)
	}

	data, _ := json.MarshalIndent(meters, "", "  ")
	t.Logf("Result:\n%s", data)
}
