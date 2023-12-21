package driver_api

import (
	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/handler/generated"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/handler/http/models"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
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

// TODO Add otel and logger
// GetTrips long pull получает в ответ список доступных (DRIVE_SEARCH) поездок
func (h *DriverHandler) GetTrips(c *gin.Context, params generated.GetTripsParams) {
	//TODO implement me
	panic("implement me")
}

func (h *DriverHandler) GetTripByID(c *gin.Context, tripId openapi_types.UUID, params generated.GetTripByIDParams) {
	//tools.ParseUUIDFromParam(c, h.logger, id_param)
	trip, err := h.driverService.Get(c, params.UserId, tripId)

	if err != nil {
		NewErrorResponse(c, h.logger, MapErrorToCode(c, err), err.Error())
		return
	}
	resp := models.ToTripResponse(trip)

	c.JSON(http.StatusOK, resp)
}

func (h *DriverHandler) AcceptTrip(c *gin.Context, tripId openapi_types.UUID, params generated.AcceptTripParams) {
	//TODO implement me
	panic("implement me")
}

func (h *DriverHandler) CancelTrip(c *gin.Context, tripId openapi_types.UUID, params generated.CancelTripParams) {
	//TODO implement me
	panic("implement me")
}

func (h *DriverHandler) EndTrip(c *gin.Context, tripId openapi_types.UUID, params generated.EndTripParams) {
	//TODO implement me
	panic("implement me")
}

func (h *DriverHandler) StartTrip(c *gin.Context, tripId openapi_types.UUID, params generated.StartTripParams) {
	//TODO implement me
	panic("implement me")
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
