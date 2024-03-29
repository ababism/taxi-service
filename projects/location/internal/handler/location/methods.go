package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/domain"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/handler/models"
	"go.uber.org/zap"
	"net/http"
)

func abortWithErrorResponse(ctx *gin.Context, logger *zap.Logger, statusCode int, message string) {
	logger.Error(fmt.Sprintf("%s: %d %s", ctx.Request.URL, statusCode, message))
	ctx.AbortWithStatusJSON(statusCode, models.Error{Message: message})
}

func abortWithBadResponse(ctx *gin.Context, logger *zap.Logger, statusCode int, message string) {
	logger.Debug(fmt.Sprintf("%s: %d %s", ctx.Request.URL, statusCode, message))
	ctx.AbortWithStatusJSON(statusCode, models.Error{Message: message})
}

func answerWithGoodResponse(ctx *gin.Context, logger *zap.Logger, statusCode int) {
	logger.Debug(fmt.Sprintf("%s: %d %v", ctx.Request.URL, statusCode))
	ctx.AbortWithStatusJSON(statusCode, http.NoBody)
}

func CallAbortByErrorCode(ctx *gin.Context, logger *zap.Logger, statusCode int, resp error) {
	errorType := statusCode / 100 % 10

	switch errorType {
	case 1, 2, 3:
		answerWithGoodResponse(ctx, logger, statusCode)
	case 4:
		abortWithBadResponse(ctx, logger, statusCode, resp.Error())
	case 5:
		abortWithErrorResponse(ctx, logger, statusCode, resp.Error())
	default:
		abortWithErrorResponse(ctx, logger, statusCode, resp.Error())
	}
}

func MapErrorToCode(err error) int {
	switch {
	case errors.Is(err, domain.ErrInternal):
		return http.StatusInternalServerError
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrTokenInvalid):
		return http.StatusUnauthorized
	case errors.Is(err, domain.ErrAccessDenied):
		return http.StatusUnauthorized
	case errors.Is(err, domain.ErrAlreadyExists):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrIncorrectBody):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrBadLatitude):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrBadLongitude):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrBadUUID):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
