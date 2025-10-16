package handler

import (
	"amr-data-bridge/internal/dto"
	"amr-data-bridge/internal/export"
	"amr-data-bridge/internal/service"
	"amr-data-bridge/internal/transport/http/middleware"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	var req dto.WaterMetersRequest

	// Parse limit
	if limitStr := q.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return middleware.NewHttpError(http.StatusBadRequest, "invalid 'limit' parameter")
		}
		req.Limit = int32(limit)
	}

	// Parse active
	if activeStr := q.Get("active"); activeStr != "" {
		switch strings.ToLower(activeStr) {
		case "true":
			val := true
			req.Active = &val
		case "false":
			val := false
			req.Active = &val
		default:
			return middleware.NewHttpError(http.StatusBadRequest, "'active' must be true, false, or omitted")
		}
	}

	// Parse and normalize type
	req.Type = strings.ToLower(strings.TrimSpace(q.Get("type")))
	if req.Type == "" {
		req.Type = "json" // default format
	}

	// Validate request DTO (enforces oneof=json csv)
	if err := validate.Struct(req); err != nil {
		return middleware.NewHttpError(http.StatusBadRequest, "validation error: "+err.Error())
	}

	// Call service
	meters, err := h.svc.GetWaterMeters(ctx, req)
	if err != nil {
		return middleware.NewHttpError(http.StatusInternalServerError, "failed to fetch water meters")
	}

	// Encode response based on validated type
	switch req.Type {
	case "csv":
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		filename := fmt.Sprintf("water_meters_%s.csv", time.Now().Format("20060102_150405"))
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		if err := export.ToCSV(w, meters); err != nil {
			return middleware.NewHttpError(http.StatusInternalServerError, "failed to export CSV")
		}
		return nil

	case "json":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(meters); err != nil {
			return middleware.NewHttpError(http.StatusInternalServerError, "failed to encode JSON response")
		}
		return nil

	default:
		return middleware.NewHttpError(http.StatusBadRequest, "unsupported export type")
	}
}
