package location_api

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/zaputil/zapctx"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/domain"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/handler/generated"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/handler/models"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/service/adapters"
	global "go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"net/http"
)

type LocationHandler struct {
	logger          *zap.Logger
	locationService adapters.LocationService
}

func NewLocationHandler(logger *zap.Logger, locationService adapters.LocationService) *LocationHandler {
	return &LocationHandler{logger: logger, locationService: locationService}
}

func (h *LocationHandler) GetDrivers(ctx *gin.Context, params generated.GetDriversParams) {
	tr := global.Tracer("gitlab/ArtemFed/mts-final-taxi")
	ctxTrace, span := tr.Start(ctx, "http: GetDrivers")
	defer span.End()

	ctxTraceLog := zapctx.WithLogger(ctxTrace, h.logger)

	drivers, err := h.locationService.GetDrivers(ctxTraceLog, params.Lat, params.Lng, params.Radius)
	if err != nil {
		NewErrorResponse(ctx, h.logger, MapErrorToCode(err), err.Error())
		return
	}

	resp := make([]generated.Driver, 0)
	for _, driver := range drivers {
		resp = append(resp, models.ToDriverResponse(driver))
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *LocationHandler) UpdateDriverLocation(ctx *gin.Context, driverId openapitypes.UUID) {
	tr := global.Tracer("gitlab/ArtemFed/mts-final-taxi")
	ctxTrace, span := tr.Start(ctx, "http: UpdateDriverLocation")
	defer span.End()
	ctxTraceLog := zapctx.WithLogger(ctxTrace, h.logger)

	var body generated.UpdateDriverLocationJSONRequestBody
	if err := ctx.BindJSON(body); err != nil {
		NewErrorResponse(ctx, h.logger, http.StatusBadRequest, domain.ErrIncorrectBody.Error())
		return
	}

	err := h.locationService.UpdateDriverLocation(ctxTraceLog, driverId, body.Lng, body.Lat)
	if err != nil {
		NewErrorResponse(ctx, h.logger, MapErrorToCode(err), err.Error())
		return
	}

	ctx.JSON(http.StatusOK, http.NoBody)
}

//func (h *DriverAPI) CreateUser(c *gin.Context) {
//	req := models.UserCreateRequest{}
//	if !restapi.BindRequestBody(c, &req) {
//		return
//	}
//
//	res, err := h.s.Create(c, domain.User{
//		Name: req.Name,
//	})
//
//	if err != nil {
//		restapi.NewErrorResponse(c, restapi.MapErrorToCode(c, err), err.Error())
//		return
//	}
//	c.JSON(http.StatusOK, models.UserResponse{Id: res.Id, Name: res.Name, Balance: res.Balance})
//}
//
//func (h *DriverAPI) GetUser(c *gin.Context) {
//	// parse from "/get-user/:id"
//	h.mylogger.Info(c.Param("id"))
//	id, ok := restapi.ParseUUIDFromParam(c)
//	if !ok {
//		return
//	}
//	user, err := h.s.Get(c, id)
//	if err != nil {
//		restapi.NewErrorResponse(c, restapi.MapErrorToCode(c, err), err.Error())
//		return
//	}
//	c.JSON(http.StatusOK, models.MakeUserResponse(user))
//}
