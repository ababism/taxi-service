package http

import (
	"fmt"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/config"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/handler/generated"
	locationAPI "gitlab/ArtemFed/mts-final-taxi/projects/location/internal/handler/location_api"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/service/adapters"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	httpPrefix = "api"
	version    = "1"
)

type Handler struct {
	logger              *zap.Logger
	cfg                 *config.Config
	driverHandler       *locationAPI.LocationHandler
	userServiceProvider adapters.LocationService
}

// HandleError is a sample error handler function
func HandleError(c *gin.Context, err error, statusCode int) {
	c.JSON(statusCode, gin.H{"error": err.Error()})
}

func InitHandler(
	router gin.IRouter,
	logger *zap.Logger,
	middlewares []generated.MiddlewareFunc,
	locationService adapters.LocationService,
) {
	locationHandler := locationAPI.NewLocationHandler(logger, locationService)

	ginOpts := generated.GinServerOptions{
		BaseURL:      fmt.Sprintf("%s/%s", httpPrefix, getVersion()),
		Middlewares:  middlewares,
		ErrorHandler: HandleError,
	}

	generated.RegisterHandlersWithOptions(router, locationHandler, ginOpts)
}

func getVersion() string {
	return fmt.Sprintf("v%s", strings.Split(version, ".")[0])
}
