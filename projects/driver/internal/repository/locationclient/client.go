package locationclient

import (
	"context"
	"encoding/json"
	"github.com/juju/zaputil/zapctx"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/repository/locationclient/generated"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
	global "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

var _ adapters.LocationClient = Client{}

type Client struct {
	httpDoer *generated.ClientWithResponses
}

func NewClient(client *generated.ClientWithResponses) *Client {
	return &Client{httpDoer: client}
}

const KeyRequestID = "RequestID"

func GetRequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(KeyRequestID).(string)
	return requestID, ok
}

func (c Client) GetDrivers(ctx context.Context, driverLocation domain.LatLngLiteral, radius float32) ([]domain.DriverLocation, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.repository.locationclient: GetDrivers")

	defer span.End()
	logger.Debug("Start GetDrivers:")

	requestID, _ := GetRequestIDFromContext(newCtx)
	var (
		resp *generated.GetDriversResponse
		err  error
	)
	if false {
		span.AddEvent("passed requestId for GetDrivers handler from location service",
			trace.WithAttributes(attribute.String(KeyRequestID, requestID)))

		reqEditor := func(newCtx context.Context, req *http.Request) error {
			req.Header.Set(KeyRequestID, requestID)
			return nil
		}
		resp, err = c.httpDoer.GetDriversWithResponse(newCtx, &generated.GetDriversParams{
			Lat:    driverLocation.Lat,
			Lng:    driverLocation.Lng,
			Radius: radius,
		}, reqEditor)
	} else {
		logger.Error("can't find RequestID in ctx")
		resp, err = c.httpDoer.GetDriversWithResponse(newCtx, &generated.GetDriversParams{
			Lat:    driverLocation.Lat,
			Lng:    driverLocation.Lng,
			Radius: radius,
		})
	}
	if err != nil {
		logger.Error("error while getting drivers from location service:", zap.Error(err))
		return nil, err
	}
	logger.Debug("Checking response dl:")

	var locationErrorMessage Error
	if resp.HTTPResponse.StatusCode != http.StatusOK {
		if err = json.Unmarshal(resp.Body, &locationErrorMessage); err != nil {
			logger.Error("error while decoding location error message JSON:", zap.Error(err))
			return nil, err
		}
		logger.Error("can't get drivers from location service ended:", zap.Int("status", resp.HTTPResponse.StatusCode), zap.Error(locationErrorMessage))
		return nil, locationErrorMessage
	}
	logger.Debug("Unmarshal dl:")

	//var driverLocations GetDriversResponse
	//var driverLocations []generated.Driver
	type GetDriversResponse struct {
		Drivers []generated.Driver `json:"drivers"`
	}
	var response GetDriversResponse

	err = json.Unmarshal(resp.Body, &response)
	if err != nil {
		logger.Error("error while decoding driver location JSON:", zap.Error(err))
		return nil, err
	}

	logger.Debug("DRIVEr LOCATIONS:", zap.Any("dLoc", response.Drivers))

	res, err := ToDriverLocationsDomain(response.Drivers)
	if err != nil {
		logger.Error("error while converting driver locations to domain:", zap.Error(err))
		return nil, err
	}
	return res, nil
}
