package handler

import (
	"amr-data-bridge/internal/dto"
	"amr-data-bridge/internal/service"
	"amr-data-bridge/internal/transport/http/middleware"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type WaterMeterHandler struct {
	svc *service.WaterMeterService
}

func NewWaterMeterHandler(svc *service.WaterMeterService) *WaterMeterHandler {
	return &WaterMeterHandler{svc: svc}
}

// GetWaterMeters godoc
// @Summary      Get active water meters
// @Description  Returns a list of all active water meters from the database.
// @Tags         water-meters
// @Produce      json
// @Success      200  {array}  dto.WaterMeterResponse
// @Failure      400  {object}  middleware.HttpError
// @Failure      500  {object}  middleware.HttpError
// @Router       /watermeters [get]
func (h *WaterMeterHandler) GetWaterMeters(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	q := r.URL.Query()

	req := dto.GetWaterMetersRequest{}

	// Parse 'limit' query param (optional)
	if limitStr := q.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return middleware.NewHttpError(http.StatusBadRequest, "invalid limit parameter")
		}
		req.Limit = int32(limit)
	}

	// Parse 'active' query param (optional)
	if activeStr := q.Get("active"); activeStr != "" {
		active := strings.ToLower(activeStr) == "true"
		req.Active = &active
	}

	// Parse 'type' query param (optional)
	req.Type = q.Get("type")

	// Validate request DTO
	if err := validate.Struct(req); err != nil {
		return middleware.NewHttpError(http.StatusBadRequest, "validation error: "+err.Error())
	}

	// Call service with DTO
	meters, err := h.svc.GetWaterMeters(ctx, req)
	if err != nil {
		return middleware.NewHttpError(http.StatusInternalServerError, "failed to fetch water meters")
	}

	// Encode response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(meters); err != nil {
		return middleware.NewHttpError(http.StatusInternalServerError, "failed to encode response")
	}

	return nil
}
