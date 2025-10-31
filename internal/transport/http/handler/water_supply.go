package handler

import (
	"amr-data-bridge/internal/dto"
	"amr-data-bridge/internal/service"
	"amr-data-bridge/internal/transport/http/middleware"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// WaterSupplyHandler handles HTTP routes for water supplies.
type WaterSupplyHandler struct {
	svc *service.WaterSupplyService
	// prefs    *config.Preferences
	validate *validator.Validate
}

// NewWaterSupplyHandler creates a new WaterSupplyHandler.
func NewWaterSupplyHandler(
	svc *service.WaterSupplyService,
	// prefs *config.Preferences,
	v *validator.Validate,
) *WaterSupplyHandler {
	return &WaterSupplyHandler{
		svc: svc,
		// prefs:    prefs,
		validate: v,
	}
}

// ImportWaterSupplies godoc
// @Summary      Import water supplies
// @Description  Imports or updates water supplies from JSON input. Each record contains supply number, water meter serial number, and geometry (latitude, longitude).
// @Tags         water-supplies
// @Accept       json
// @Produce      json
// @Param        body  body  []dto.WaterSupplyRequest  true  "List of water supplies"
// @Success      200  {array}  dto.WaterSupplyResponse
// @Failure      400  {object}  middleware.HttpError
// @Failure      500  {object}  middleware.HttpError
// @Router       /water-supplies/import [post]
func (h *WaterSupplyHandler) ImportWaterSupplies(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var payload []dto.WaterSupplyRequest

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return middleware.NewHttpError(http.StatusBadRequest, "invalid JSON: "+err.Error())
	}

	for i, item := range payload {
		if err := h.validate.Struct(item); err != nil {
			return middleware.NewHttpError(http.StatusBadRequest,
				fmt.Sprintf("validation error in item %d: %s", i, err.Error()))
		}
	}

	// 💡 Let the service handle transaction management internally
	result, err := h.svc.ImportWaterSupplies(ctx, payload)
	if err != nil {
		return middleware.NewHttpError(http.StatusInternalServerError, err.Error())
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		return middleware.NewHttpError(http.StatusInternalServerError, "failed to encode JSON response")
	}
	return nil
}
