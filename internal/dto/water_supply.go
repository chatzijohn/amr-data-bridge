package dto

// Request DTO — what comes from HTTP/API
type WaterSupplyRequest struct {
	SupplyNumber string  `json:"supplyNumber" binding:"required"`
	Latitude     float64 `json:"latitude" binding:"required"`
	Longitude    float64 `json:"longitude" binding:"required"`
	SerialNumber string  `json:"serialNumber,omitempty"`
	Active       bool    `json:"Active,omitempty"`
}
