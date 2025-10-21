package mapper

import (
	"amr-data-bridge/internal"
	"amr-data-bridge/internal/db"
	"amr-data-bridge/internal/dto"
	"time"
)

// WaterMeterToDTO converts a DB row to a DTO,
// including only the fields specified in preferences.yaml.
func WaterMeterToDTO(m db.GetWaterMetersRow, prefs *internal.Preferences) dto.WaterMeterResponse {
	var (
		lastSeen     string
		supplyNumber string
	)

	if m.LastSeen.Valid {
		lastSeen = m.LastSeen.Time.UTC().Format(time.RFC3339)
	}

	if m.SupplyNumber.Valid {
		supplyNumber = m.SupplyNumber.String
	}

	// Build a quick lookup map of allowed fields.
	allowed := make(map[string]struct{}, len(prefs.Export.WaterMeterFields))
	for _, f := range prefs.Export.WaterMeterFields {
		allowed[f] = struct{}{}
	}

	// Conditionally populate DTO fields based on preferences.
	var res dto.WaterMeterResponse

	if _, ok := allowed["DevEUI"]; ok {
		res.DevEUI = m.DevEUI
	}
	if _, ok := allowed["SupplyNumber"]; ok {
		res.SupplyNumber = supplyNumber
	}
	if _, ok := allowed["SerialNumber"]; ok {
		res.SerialNumber = m.SerialNumber
	}
	if _, ok := allowed["BrandName"]; ok {
		res.BrandName = m.BrandName
	}
	if _, ok := allowed["LtPerPulse"]; ok {
		res.LtPerPulse = m.LtPerPulse
	}
	if _, ok := allowed["IsActive"]; ok {
		res.IsActive = m.IsActive
	}
	if _, ok := allowed["AlarmStatus"]; ok {
		res.AlarmStatus = m.AlarmStatus
	}
	if _, ok := allowed["NoFlow"]; ok {
		res.NoFlow = m.NoFlow
	}
	if _, ok := allowed["CurrentReading"]; ok {
		res.CurrentReading = m.CurrentReading.Int32
	}
	if _, ok := allowed["LastSeen"]; ok {
		res.LastSeen = lastSeen
	}

	return res
}

// WaterMetersToDTO maps a slice of DB results to DTOs,
// applying the same preferences-based filtering.
func WaterMetersToDTO(models []db.GetWaterMetersRow, prefs *internal.Preferences) []dto.WaterMeterResponse {
	out := make([]dto.WaterMeterResponse, 0, len(models))
	for _, m := range models {
		out = append(out, WaterMeterToDTO(m, prefs))
	}
	return out
}
