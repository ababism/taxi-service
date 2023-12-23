package driver_api

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/zaputil/zapctx"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/handler/generated"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/handler/http/models"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
	global "go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"net/http"
)

const idParam = "user_id"

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
func (h *DriverHandler) GetTrips(ginCtx *gin.Context, params generated.GetTripsParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(ginCtx, "http: GetTripByID")
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	trips, err := h.driverService.GetTrips(ctx, params.UserId)
	if err != nil {
		AbortWithBadResponse(ginCtx, h.logger, MapErrorToCode(err), err.Error())
		return
	}
	resp := models.ToTripsResponse(trips)

	ginCtx.JSON(http.StatusOK, resp)
}

func (h *DriverHandler) GetTripByID(ginCtx *gin.Context, tripId openapi_types.UUID, params generated.GetTripByIDParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(ginCtx, "http: GetTripByID")
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	trip, err := h.driverService.GetTripByID(ctx, params.UserId, tripId)
	if err != nil {
		AbortWithBadResponse(ginCtx, h.logger, MapErrorToCode(err), err.Error())
		return
	}
	resp := models.ToTripResponse(*trip)

	ginCtx.JSON(http.StatusOK, resp)
}

func (h *DriverHandler) AcceptTrip(ginCtx *gin.Context, tripId openapi_types.UUID, params generated.AcceptTripParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(ginCtx, "http: AcceptTrip")
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	err := h.driverService.AcceptTrip(ctx, params.UserId, tripId)
	if err != nil {
		AbortWithBadResponse(ginCtx, h.logger, MapErrorToCode(err), err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, http.NoBody)
}

func (h *DriverHandler) CancelTrip(ginCtx *gin.Context, tripId openapi_types.UUID, params generated.CancelTripParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(ginCtx, "http: CancelTrip")
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	err := h.driverService.CancelTrip(ctx, params.UserId, tripId, params.Reason)
	if err != nil {
		AbortWithBadResponse(ginCtx, h.logger, MapErrorToCode(err), err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, http.NoBody)
}

func (h *DriverHandler) EndTrip(ginCtx *gin.Context, tripId openapi_types.UUID, params generated.EndTripParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(ginCtx, "http: EndTrip")
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	err := h.driverService.EndTrip(ctx, params.UserId, tripId)
	if err != nil {
		AbortWithBadResponse(ginCtx, h.logger, MapErrorToCode(err), err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, http.NoBody)
}

func (h *DriverHandler) StartTrip(ginCtx *gin.Context, tripId openapi_types.UUID, params generated.StartTripParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(ginCtx, "http: StartTrip")
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	err := h.driverService.StartTrip(ctx, params.UserId, tripId)
	if err != nil {
		AbortWithBadResponse(ginCtx, h.logger, MapErrorToCode(err), err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, http.NoBody)
}
