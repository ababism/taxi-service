package mongo

import (
	"context"
	"github.com/google/uuid"
	"github.com/juju/zaputil/zapctx"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	global "go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// GetTripByID returns Trip from driver collection by trip_id
func (r *DriveRepository) GetTripByID(ctx context.Context, tripId uuid.UUID) (domain.Trip, error) {
	log := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.repository.mongo: GetTripByID")
	defer span.End()

	//collection := r.client.Database("your_database_name").Collection(r.driverCollection.Name())

	var trip domain.Trip
	filter := bson.M{"trip_id": tripId}
	err := r.driverCollection.FindOne(newCtx, filter).Decode(&trip)
	if err != nil {
		log.Error("Error fetching trip from MongoDB", zap.Error(err))
		return domain.Trip{}, err
	}

	return trip, nil
}

// ChangeTripStatus changes Trip (with tripId from parameter) status in db to status parameter
func (r *DriveRepository) ChangeTripStatus(ctx context.Context, tripId uuid.UUID, status domain.TripStatus) error {
	log := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.repository.mongo: ChangeTripStatus")
	defer span.End()

	filter := bson.M{"trip_id": tripId}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := r.driverCollection.UpdateOne(newCtx, filter, update)
	if err != nil {
		log.Error("Failed to update trip status in MongoDB", zap.Error(err))
		return err
	}

	return nil
}

// CreateTrip creates Trip
func (r *DriveRepository) CreateTrip(ctx context.Context, trip domain.Trip) error {
	log := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.repository.mongo: CreateTrip")
	defer span.End()

	mongoTrip := models.ToMongoTripModel(trip)
	// Insert the trip into MongoDB
	_, err := r.driverCollection.InsertOne(newCtx, mongoTrip)
	if err != nil {
		log.Error("Failed to create trip in MongoDB", zap.Error(err))
		return err
	}

	return nil
}
