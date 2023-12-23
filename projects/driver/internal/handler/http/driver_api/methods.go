package driver_api

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/handler/http/models"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
)

// TODO maybe c.AbortWithStatusJSON ?
func AbortWithBadResponse(c *gin.Context, logger *zap.Logger, statusCode int, message string) {
	logger.Debug(fmt.Sprintf("%s: %d %s", c.Request.URL, statusCode, message))
	c.AbortWithStatusJSON(statusCode, models.Error{Message: message})
}

func AbortWithErrorResponse(c *gin.Context, logger *zap.Logger, statusCode int, message string) {
	logger.Error(fmt.Sprintf("%s: %d %s", c.Request.URL, statusCode, message))
	c.AbortWithStatusJSON(statusCode, models.Error{Message: message})
}

//func BindRequestBody(c *gin.Context, logger *zap.Logger, obj any) bool {
//	if err := c.BindJSON(obj); err != nil {
//		AbortWithBadResponse(c, logger, http.StatusBadRequest, domain.ErrIncorrectBody.Error())
//		return false
//	}
//	return true
//}

//func ErrorResponse(c *gin.Context, statusCode int, err error) {
//	logrus.Infof("%s: %d %s", c.Request.URL, statusCode, err.Error())
//	c.AbortWithStatusJSON(statusCode, models.Error{Message: err.Error()})
//}

func MapErrorToCode(err error) int {
	switch {
	case errors.Is(err, domain.ErrInternal):
		return http.StatusInternalServerError
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrIncorrectBody):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrTokenInvalid):
		return http.StatusUnauthorized
	case errors.Is(err, domain.ErrAccessDenied):
		return http.StatusUnauthorized
	case errors.Is(err, domain.ErrAlreadyExists):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrBadUUID):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// ParseUUIDFromParam makes Error response if it couldn't parse token, returns true if everything is ok
func ParseUUIDFromParam(c *gin.Context, l *zap.Logger, key string) (uuid.UUID, bool) {
	id := c.Param(key)
	itemUUID, err := uuid.Parse(id)
	if err != nil {
		AbortWithBadResponse(c, l, MapErrorToCode(domain.ErrBadUUID), domain.ErrBadUUID.Error())
		return uuid.Nil, false
	}
	return itemUUID, true
}

//func parseUUID(c *gin.Context, id string) (uuid.UUID, bool) {
//	itemUUID, err := uuid.Parse(id)
//	if err != nil {
//		AbortWithBadResponse(c, MapErrorToCode(c, adapters.ErrBadUUID), adapters.ErrBadUUID.Error())
//		return uuid.Nil, false
//	}
//	return itemUUID, true
//}
