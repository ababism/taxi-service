package http

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/config"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/handler/generate"
	driverAPI "gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/handler/http/driver_api"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
)

const (
	httpPrefix = "api"
	version    = "1"
)

type Handler struct {
	logger              *zap.Logger
	cfg                 *config.Config
	driverHandler       *driverAPI.DriverHandler
	userServiceProvider adapters.DriverService
}

// HandleError is a sample error handler function
func HandleError(c *gin.Context, err error, statusCode int) {
	c.JSON(statusCode, gin.H{"error": err.Error()})
}

func InitHandler(
	router gin.IRouter,
	logger *zap.Logger,
	middlewares []generate.MiddlewareFunc,
	driverService *adapters.DriverService,
) {
	driverHandler := driverAPI.NewDriverHandler(logger, driverService)

	ginOpts := generate.GinServerOptions{
		BaseURL:      fmt.Sprintf("%s/%s", httpPrefix, getVersion()),
		Middlewares:  middlewares,
		ErrorHandler: HandleError,
	}

	generate.RegisterHandlersWithOptions(router, driverHandler, ginOpts)
}

func getVersion() string {
	return fmt.Sprintf("v%s", strings.Split(version, ".")[0])
}
