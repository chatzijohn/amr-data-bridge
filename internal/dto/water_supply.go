package dto

// Request DTO — what comes from HTTP/API
type WaterSupplyRequest struct {
	SupplyNumber string  `json:"supplyNumber" binding:"required"`
	Latitude     float64 `json:"latitude" binding:"required"`
	Longitude    float64 `json:"longitude" binding:"required"`
	SerialNumber string  `json:"serialNumber,omitempty"`
	Active       bool    `json:"active,omitempty"`
}

// Response DTO — what we send back to the client
type WaterSupplyResponse struct {
	ID           int64   `json:"id"`
	SupplyNumber string  `json:"supplyNumber"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	SerialNumber string  `json:"serialNumber,omitempty"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}
