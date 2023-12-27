package api

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/zaputil/zapctx"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"gitlab.com/ArtemFed/mts-final-taxi/pkg/app"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/domain"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/handler/generated"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/handler/models"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/service/adapters"
	global "go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"net/http"
)

type LocationHandler struct {
	logger          *zap.Logger
	locationService adapters.LocationService
	appCfg          *app.Config
}

func NewLocationHandler(logger *zap.Logger, locationService adapters.LocationService, appCfg *app.Config) *LocationHandler {
	return &LocationHandler{logger: logger, locationService: locationService, appCfg: appCfg}
}

func (h *LocationHandler) GetDrivers(ctx *gin.Context, params generated.GetDriversParams) {
	tr := global.Tracer(domain.TracerName)
	ctxTrace, span := tr.Start(ctx, "location.http: GetDrivers")
	defer span.End()

	ctxTraceLog := zapctx.WithLogger(ctxTrace, h.logger)
	drivers, err := h.locationService.GetDrivers(ctxTraceLog, params.Lat, params.Lng, params.Radius)
	if err != nil {
		CallAbortByErrorCode(ctx, h.logger, MapErrorToCode(err), err)
		return
	}

	resp := make([]generated.Driver, len(drivers))
	for i, driver := range drivers {
		resp[i] = models.ToDriverResponse(driver)
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *LocationHandler) UpdateDriverLocation(ctx *gin.Context, driverId openapitypes.UUID) {
	tr := global.Tracer(domain.TracerName)
	ctxTrace, span := tr.Start(ctx, "location.http: UpdateDriverLocation")
	defer span.End()

	ctxTraceLog := zapctx.WithLogger(ctxTrace, h.logger)

	var body generated.UpdateDriverLocationJSONRequestBody
	if err := ctx.BindJSON(&body); err != nil {
		CallAbortByErrorCode(ctx, h.logger, http.StatusBadRequest, domain.ErrIncorrectBody)
		return
	}

	err := h.locationService.UpdateDriverLocation(ctxTraceLog, driverId, body.Lat, body.Lng)
	if err != nil {
		CallAbortByErrorCode(ctx, h.logger, MapErrorToCode(err), err)
		return
	}

	ctx.JSON(http.StatusOK, http.NoBody)
}
