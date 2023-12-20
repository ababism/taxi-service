package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
)

//func BindRequestBody(c *gin.Context, obj any) bool {
//	if err := c.BindJSON(obj); err != nil {
//		NewErrorResponse(c, http.StatusBadRequest, domain.ErrIncorrectBody.Error())
//		return false
//	}
//	return true
//}

//func NewErrorResponse(c *gin.Context, statusCode int, message string) {
//	logrus.Infof("%s: %d %s", c.Request.URL, statusCode, message)
//	c.AbortWithStatusJSON(statusCode, models.Error{Message: message})
//}
//func ErrorResponse(c *gin.Context, statusCode int, err error) {
//	logrus.Infof("%s: %d %s", c.Request.URL, statusCode, err.Error())
//	c.AbortWithStatusJSON(statusCode, models.Error{Message: err.Error()})
//}

func MapErrorToCode(c *gin.Context, err error) int {
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
//func ParseUUIDFromParam(c *gin.Context) (uuid.UUID, bool) {
//	id := c.Param("id")
//	itemUUID, err := uuid.Parse(id)
//	if err != nil {
//		NewErrorResponse(c, MapErrorToCode(c, domain.ErrBadUUID), domain.ErrBadUUID.Error())
//		return uuid.Nil, false
//	}
//	return itemUUID, true
//}

//func parseUUID(c *gin.Context, id string) (uuid.UUID, bool) {
//	itemUUID, err := uuid.Parse(id)
//	if err != nil {
//		NewErrorResponse(c, MapErrorToCode(c, adapters.ErrBadUUID), adapters.ErrBadUUID.Error())
//		return uuid.Nil, false
//	}
//	return itemUUID, true
//}
