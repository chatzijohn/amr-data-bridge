package dto

// Request DTO — what comes from HTTP/API
type WaterMetersRequest struct {
	Limit  int32  `json:"limit" validate:"gte=0,lte=1000"`
	Active *bool  `json:"active" validate:"omitempty"`
	Type   string `json:"type" validate:"omitempty,oneof=json csv"` // Allowed types
}

// Response DTO — what you return to API clients
type WaterMeterResponse struct {
	DevEUI         string `json:"devEUI"`
	SerialNumber   string `json:"serialNumber"`
	BrandName      string `json:"brandName"`
	LtPerPulse     int32  `json:"ltPerPulse"`
	IsActive       bool   `json:"isActive"`
	AlarmStatus    bool   `json:"alarmStatus"`
	NoFlow         bool   `json:"noFlow"`
	CurrentReading int32  `json:"currentReading"`
}
