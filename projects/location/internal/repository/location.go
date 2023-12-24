package repository

import (
	"context"
	"database/sql"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/domain"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/service/adapters"
	global "go.opentelemetry.io/otel"

	"github.com/jmoiron/sqlx"
)

var _ adapters.LocationRepository = &locationRepository{}

const (
	getDriversByRadius = `
		SELECT id, lat, lng
		FROM drivers_locations
		WHERE earth_box(ll_to_earth($1, $2), $3) @> ll_to_earth(lat, lng);`
	//getDriversByRadius = `SELECT id, ST_X(location) as lat, ST_Y(location) as lng FROM drivers_locations WHERE ST_DWithin(location, ST_MakePoint($1, $2)::geography, $3);`
	updateDrivers = `
		INSERT INTO drivers_locations (id, lat, lng)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET lat = EXCLUDED.lat, lng = EXCLUDED.lng`
	//updateDrivers      = `INSERT INTO drivers_locations (id, location) VALUES ($1, ST_MakePoint($2, $3)::geography) ON CONFLICT (id) DO UPDATE SET location = EXCLUDED.location;`
)

type locationRepository struct {
	db *sqlx.DB
}

func NewLocalRepository(repos *sqlx.DB) adapters.LocationRepository {
	return &locationRepository{
		db: repos,
	}
}

// GetDrivers получает всех водителей в заданном радиусе от точки
func (r *locationRepository) GetDrivers(ctx context.Context, lat float32, lng float32, radius float32) ([]domain.Driver, error) {
	tr := global.Tracer(domain.TracerName)
	_, span := tr.Start(ctx, "location.repository: GetDrivers")
	defer span.End()

	// ATTENTION: Считаем, что radius в метрах
	rows, err := r.db.Query(getDriversByRadius, float64(lat), float64(lng), float64(radius))
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var drivers []domain.Driver
	for rows.Next() {
		var driver domain.Driver
		err = rows.Scan(&driver.Id, &driver.Lat, &driver.Lng)
		if err != nil {
			return nil, err
		}
		drivers = append(drivers, driver)
	}

	return drivers, nil
}

// UpdateDriverLocation Обновляет значение позиции у водителей. При отсутствии добавляет запись.
func (r *locationRepository) UpdateDriverLocation(ctx context.Context, driverId openapitypes.UUID, lat float32, lng float32) error {
	tr := global.Tracer(domain.TracerName)
	_, span := tr.Start(ctx, "location.repository: UpdateDriverLocation")
	defer span.End()

	_, err := r.db.Exec(updateDrivers, driverId, lat, lng)
	return err
}
