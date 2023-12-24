package driver_api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/juju/zaputil/zapctx"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/handler/generated"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/handler/http/models"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
	global "go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const idParam = "user_id"

var _ generated.ServerInterface = &DriverHandler{}

type DriverHandler struct {
	logger        *zap.Logger
	driverService adapters.DriverService
	WaitTimeout   time.Duration
}

func NewDriverHandler(logger *zap.Logger, driverService adapters.DriverService, socketTimeout time.Duration) *DriverHandler {
	return &DriverHandler{logger: logger, driverService: driverService, WaitTimeout: socketTimeout}
}

// TODO websocket
// GetTrips long pull получает в ответ список доступных (DRIVE_SEARCH) поездок за время поллинга
func (h *DriverHandler) GetTrips(ginCtx *gin.Context, params generated.GetTripsParams) {
	tr := global.Tracer(domain.ServiceName)
	ctxTrace, span := tr.Start(ginCtx, "http: GetTripByID")
	defer span.End()

	ctx := zapctx.WithLogger(ctxTrace, h.logger)

	upg := websocket.Upgrader{}
	conn, err := upg.Upgrade(ginCtx.Writer, ginCtx.Request, nil)
	defer conn.Close()

	if err != nil {
		http.Error(ginCtx.Writer, "Could not upgrade to WebSocket", http.StatusBadRequest)
		return
	}

	driverId := params.UserId.String()
	incomingTrip, ok := domain.IncomingTrips.GetTripChannel(&driverId)
	if ok {
		h.logger.Error("driver should not be in event map")
	} else {
		domain.IncomingTrips.AddTrip(&driverId, make(chan *uuid.UUID))
	}
	// Не забудем в конце очистить информацию о водителе в ожидании
	defer domain.IncomingTrips.DeleteTripChannel(&driverId)

	timer := time.NewTimer(h.WaitTimeout).C
	tripsResp := make([]domain.Trip, 0)
	go func() {
		tripId := <-incomingTrip
		trip, err := h.driverService.GetTripByID(ctx, params.UserId, *tripId)
		if err == nil {
			tripsResp = append(tripsResp, *trip)
			resp := models.ToTripResponse(*trip)
			conn.WriteJSON(resp)
		}
	}()
	<-timer

	resp := models.ToTripsResponse(tripsResp)

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
