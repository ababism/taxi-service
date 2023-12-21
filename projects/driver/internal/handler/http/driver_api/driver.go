package driver_api

import (
	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/handler/generated"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/handler/http/models"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
	global "go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"net/http"
)

const id_param = "user_id"

var _ generated.ServerInterface = &DriverHandler{}

type DriverHandler struct {
	logger        *zap.Logger
	driverService adapters.DriverService
}

func NewDriverHandler(logger *zap.Logger, driverService adapters.DriverService) *DriverHandler {
	return &DriverHandler{logger: logger, driverService: driverService}
}

// TODO websocket
// GetTrips long pull получает в ответ список доступных (DRIVE_SEARCH) поездок
func (h *DriverHandler) GetTrips(c *gin.Context, params generated.GetTripsParams) {
	tr := global.Tracer("gitlab/ArtemFed/mts-final-taxi")
	newCtx, span := tr.Start(c, "http: GetTripByID")
	defer span.End()

	trips, err := h.driverService.GetTrips(newCtx, params.UserId)
	if err != nil {
		NewErrorResponse(c, h.logger, MapErrorToCode(c, err), err.Error())
		return
	}
	resp := models.ToTripsResponse(trips)

	c.JSON(http.StatusOK, resp)
}

func (h *DriverHandler) GetTripByID(c *gin.Context, tripId openapi_types.UUID, params generated.GetTripByIDParams) {
	tr := global.Tracer("gitlab/ArtemFed/mts-final-taxi")
	newCtx, span := tr.Start(c, "http: GetTripByID")
	defer span.End()
	//c = zapctx.WithLogger(newCtx, h.logger)

	trip, err := h.driverService.GetTripByID(newCtx, params.UserId, tripId)
	if err != nil {
		NewErrorResponse(c, h.logger, MapErrorToCode(c, err), err.Error())
		return
	}
	resp := models.ToTripResponse(trip)

	c.JSON(http.StatusOK, resp)
}

func (h *DriverHandler) AcceptTrip(c *gin.Context, tripId openapi_types.UUID, params generated.AcceptTripParams) {
	tr := global.Tracer("gitlab/ArtemFed/mts-final-taxi")
	newCtx, span := tr.Start(c, "http: AcceptTrip")
	defer span.End()

	_, err := h.driverService.AcceptTrip(newCtx, params.UserId, tripId)
	if err != nil {
		NewErrorResponse(c, h.logger, MapErrorToCode(c, err), err.Error())
		return
	}

	c.JSON(http.StatusOK, http.NoBody)
}

func (h *DriverHandler) CancelTrip(c *gin.Context, tripId openapi_types.UUID, params generated.CancelTripParams) {
	tr := global.Tracer("gitlab/ArtemFed/mts-final-taxi")
	newCtx, span := tr.Start(c, "http: CancelTrip")
	defer span.End()

	_, err := h.driverService.CancelTrip(newCtx, params.UserId, tripId, params.Reason)
	if err != nil {
		NewErrorResponse(c, h.logger, MapErrorToCode(c, err), err.Error())
		return
	}

	c.JSON(http.StatusOK, http.NoBody)
}

func (h *DriverHandler) EndTrip(c *gin.Context, tripId openapi_types.UUID, params generated.EndTripParams) {
	tr := global.Tracer("gitlab/ArtemFed/mts-final-taxi")
	newCtx, span := tr.Start(c, "http: EndTrip")
	defer span.End()

	_, err := h.driverService.EndTrip(newCtx, params.UserId, tripId)
	if err != nil {
		NewErrorResponse(c, h.logger, MapErrorToCode(c, err), err.Error())
		return
	}

	c.JSON(http.StatusOK, http.NoBody)
}

func (h *DriverHandler) StartTrip(c *gin.Context, tripId openapi_types.UUID, params generated.StartTripParams) {
	tr := global.Tracer("gitlab/ArtemFed/mts-final-taxi")
	newCtx, span := tr.Start(c, "http: StartTrip")
	defer span.End()

	_, err := h.driverService.StartTrip(newCtx, params.UserId, tripId)
	if err != nil {
		NewErrorResponse(c, h.logger, MapErrorToCode(c, err), err.Error())
		return
	}

	c.JSON(http.StatusOK, http.NoBody)
}
