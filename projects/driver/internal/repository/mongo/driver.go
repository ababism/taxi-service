package mongo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/juju/zaputil/zapctx"
	"github.com/pkg/errors"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	global "go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// GetTripsByStatus returns all Trips with given status
func (r *DriveRepository) GetTripsByStatus(ctx context.Context, status domain.TripStatus) ([]domain.Trip, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.repository.mongo: GetTripsByStatus")
	defer span.End()

	filter := bson.M{"status": status}

	cursor, err := r.driverCollection.Find(newCtx, filter)
	if err != nil {
		// TODO: errors.Wrap добавить в остальные функции?
		errExplanation := "Error finding trips by status"
		logger.Error(errExplanation, zap.Error(err))
		return nil, errors.Wrap(err, errExplanation)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, newCtx)

	var trips []domain.Trip
	for cursor.Next(newCtx) {
		var trip domain.Trip
		err = cursor.Decode(&trip)
		if err != nil {
			errExplanation := fmt.Sprintf("Error decoding trip by status %s", status)
			logger.Error(errExplanation, zap.Error(err))
			return nil, errors.Wrap(err, errExplanation)
		}
		trips = append(trips, trip)
	}

	if err = cursor.Err(); err != nil {
		errExplanation := "Cursor error"
		logger.Error(errExplanation, zap.Error(err))
		return nil, err
	}

	return trips, nil
}

// GetTripByID returns Trip from driver collection by trip_id
func (r *DriveRepository) GetTripByID(ctx context.Context, tripId uuid.UUID) (*domain.Trip, error) {
	log := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.repository.mongo: GetTripByID")
	defer span.End()

	var trip domain.Trip
	filter := bson.M{"trip_id": tripId}
	err := r.driverCollection.FindOne(newCtx, filter).Decode(&trip)
	if err != nil {
		log.Error("Error fetching trip from MongoDB", zap.Error(err))
		return nil, err
	}

	return &trip, nil
}

// UpdateTrip updates all fields of the given trip in the database
func (r *DriveRepository) UpdateTrip(ctx context.Context, tripId uuid.UUID, updatedTrip domain.Trip) error {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.repository.mongo: UpdateTrip")
	defer span.End()

	filter := bson.M{"trip_id": tripId}

	updateDoc, err := bson.Marshal(updatedTrip)
	if err != nil {
		logger.Error("Failed to marshal updated trip", zap.Error(err))
		return err
	}

	var updateMap map[string]interface{}
	err = bson.Unmarshal(updateDoc, &updateMap)
	if err != nil {
		logger.Error("Failed to unmarshal updated trip", zap.Error(err))
		return err
	}

	for key, value := range updateMap {
		if value == nil {
			delete(updateMap, key)
		}
	}

	update := bson.M{"$set": updateMap}

	_, err = r.driverCollection.UpdateOne(newCtx, filter, update)
	if err != nil {
		logger.Error("Failed to update trip in MongoDB", zap.Error(err))
		return err
	}

	return nil
}

// ChangeTripStatus changes Trip (with tripId from parameter) status in db to status parameter
func (r *DriveRepository) ChangeTripStatus(ctx context.Context, tripId uuid.UUID, status domain.TripStatus) error {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.repository.mongo: ChangeTripStatus")
	defer span.End()

	filter := bson.M{"trip_id": tripId}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := r.driverCollection.UpdateOne(newCtx, filter, update)
	if err != nil {
		logger.Error("Failed to update trip status in MongoDB", zap.Error(err))
		return err
	}

	return nil
}

// InsertTrip inserts Trip
func (r *DriveRepository) InsertTrip(ctx context.Context, trip domain.Trip) error {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.repository.mongo: InsertTrip")
	defer span.End()

	mongoTrip := models.ToMongoTripModel(trip)
	_, err := r.driverCollection.InsertOne(newCtx, mongoTrip)
	if err != nil {
		logger.Error("Failed to create trip in MongoDB", zap.Error(err))
		return err
	}

	return nil
}
