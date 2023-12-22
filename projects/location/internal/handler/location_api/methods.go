package location_api

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/domain"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/handler/models"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO maybe c.AbortWithStatusJSON ?
func NewErrorResponse(c *gin.Context, l *zap.Logger, statusCode int, message string) {
	l.Info(fmt.Sprintf("%s: %d %s", c.Request.URL, statusCode, message))
	c.AbortWithStatusJSON(statusCode, models.Error{Message: message})
}

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
		NewErrorResponse(c, l, MapErrorToCode(domain.ErrBadUUID), domain.ErrBadUUID.Error())
		return uuid.Nil, false
	}
	return itemUUID, true
}

//func parseUUID(c *gin.Context, id string) (uuid.UUID, bool) {
//	itemUUID, err := uuid.Parse(id)
//	if err != nil {
//		NewErrorResponse(c, MapErrorToCode(c, adapters.ErrBadUUID), adapters.ErrBadUUID.Error())
//		return uuid.Nil, false
//	}
//	return itemUUID, true
//}
