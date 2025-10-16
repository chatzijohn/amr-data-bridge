package mapper

import (
	"amr-data-bridge/internal/db"
	"amr-data-bridge/internal/dto"
)

func WaterMeterToDTO(m db.WaterMeter) dto.WaterMeterResponse {

	return dto.WaterMeterResponse{
		DevEUI:         m.DevEUI,
		SerialNumber:   m.SerialNumber,
		BrandName:      m.SerialNumber,
		LtPerPulse:     m.LtPerPulse,
		IsActive:       m.IsActive,
		AlarmStatus:    m.AlarmStatus,
		NoFlow:         m.NoFlow,
		CurrentReading: m.CurrentReading.Int32,
	}
}

func WaterMetersToDTO(models []db.WaterMeter) []dto.WaterMeterResponse {
	out := make([]dto.WaterMeterResponse, 0, len(models))
	for _, m := range models {
		out = append(out, WaterMeterToDTO(m))
	}
	return out
}
