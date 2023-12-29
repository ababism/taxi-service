package driver_api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juju/zaputil/zapctx"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/handler/generated"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/handler/http/models"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
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

	//upg := websocket.Upgrader{}
	//conn, err := upg.Upgrade(ginCtx.Writer, ginCtx.Request, nil)
	//defer conn.Close()
	//if err != nil {
	//	http.Error(ginCtx.Writer, "Could not upgrade to WebSocket", http.StatusBadRequest)
	//	return
	//}

	driverId := params.UserId.String()
	trEventCh, ok := domain.AvailableTripEvents.GetTripChannel(driverId)
	// Не забудем в конце очистить информацию о водителе в ожидании
	defer domain.AvailableTripEvents.DeleteTripChannel(driverId)
	if ok {
		close(trEventCh)
		h.logger.Fatal("[driver.handler: GetTrips] new driver should not be in event map")
		return
	} else {
		trEventCh = make(chan uuid.UUID)
		domain.AvailableTripEvents.AddTripChannel(driverId, trEventCh)
	}

	//timer := time.NewTimer(h.WaitTimeout).C
	timer := time.After(h.WaitTimeout)
	tripsResp := make([]domain.Trip, 0)

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			select {
			case tripId, ok := <-trEventCh:
				if !ok {
					return
				}
				h.logger.Debug("[driver.api] GetTrips: got event for driver", zap.String("trip_id", tripId.String()))
				trip, err := h.driverService.GetTripByID(ctx, params.UserId, tripId)
				if err == nil {
					tripsResp = append(tripsResp, *trip)
					//resp := models.ToTripResponse(*trip)
					//conn.WriteJSON(resp)
				}
				continue
			}
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
