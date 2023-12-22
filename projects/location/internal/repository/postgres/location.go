package postgres

import (
	"context"
	"database/sql"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/domain"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/service/adapters"
	global "go.opentelemetry.io/otel"

	"github.com/jmoiron/sqlx"
)

var _ adapters.LocationRepository = &localRepository{}

const (
	getDriversByRadius = `SELECT id, lat, lng FROM drivers_locations WHERE earth_box(ll_to_earth($1, $2), $3) @> ll_to_earth(lat, lng)", lat, lng, radius`
	//getDriversByRadius = `SELECT id, ST_X(location) as lat, ST_Y(location) as lng FROM drivers_locations WHERE ST_DWithin(location, ST_MakePoint($1, $2)::geography, $3);`
	updateDrivers = `INSERT INTO drivers_locations (id, lat, lng) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET lat = EXCLUDED.lat, lng = EXCLUDED.lng`
	//updateDrivers      = `INSERT INTO drivers_locations (id, location) VALUES ($1, ST_MakePoint($2, $3)::geography) ON CONFLICT (id) DO UPDATE SET location = EXCLUDED.location;`
)

type localRepository struct {
	db *sqlx.DB
}

// GetDrivers получает всех водителей в заданном радиусе от точки
func (r *localRepository) GetDrivers(ctx context.Context, lat float32, lng float32, radius float32) ([]domain.Driver, error) {
	tr := global.Tracer("gitlab/ArtemFed/mts-final-taxi")
	_, span := tr.Start(ctx, "location-service: GetDrivers Repository")
	defer span.End()

	// ATTENTION: Считаем, что radius в метрах
	rows, err := r.db.Query(getDriversByRadius, lng, lat, radius)
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

// UpdateDriverLocation Обновляет значение позиции у водителей. При отсутствии создаёт
func (r *localRepository) UpdateDriverLocation(ctx context.Context, driverId openapitypes.UUID, lat float32, lng float32) error {
	tr := global.Tracer("gitlab/ArtemFed/mts-final-taxi")
	_, span := tr.Start(ctx, "location-service: UpdateDriverLocation Repository")
	defer span.End()

	_, err := r.db.Exec(updateDrivers, driverId, lng, lat)
	return err
}

func NewLocalRepository(repos *sqlx.DB) adapters.LocationRepository {
	return &localRepository{
		db: repos,
	}
}

//func (r localRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
//	q := `
//	INSERT INTO users (username, balance)
//	VALUES ($1, $2)
//	RETURNING *
//	`
//	logrus.Trace(formatQuery(q))
//	var resUser models.User
//	err := r.db.Get(&resUser, q, user.Name, 0)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return domain.User{}, domain.ErrInternal
//		}
//		return domain.User{}, err
//	}
//	return resUser.ToDomain(), nil
//}
//
//func (r localRepository) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
//	q := `
//	SELECT * FROM users WHERE id = $1
//	`
//	logrus.Trace(formatQuery(q))
//	var user models.User
//	err := r.db.Get(&user, q, id)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return domain.User{}, domain.ErrNotFound
//		}
//		return domain.User{}, err
//	}
//	return user.ToDomain(), nil
//}
//
//func (r localRepository) GetBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
//	q := `
//	SELECT * FROM users WHERE id = $1
//	`
//	logrus.Trace(formatQuery(q))
//	var user models.User
//	err := r.db.Get(&user, q, userID)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return 0, domain.ErrNotFound
//		}
//		return 0, err
//	}
//	return user.Balance, nil
//}
//
//// AddToBalance uses Serializable isolation level
//func (r localRepository) AddToBalance(ctx context.Context, userID uuid.UUID, amount float64) (float64, error) {
//	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
//	if err != nil {
//		return 0, domain.ErrInternal
//	}
//
//	// Defer a rollback in case anything fails. (from go.dev)
//	// defer выполняется перед паникой
//	defer func(tx *sql.Tx) {
//		_ = tx.Rollback()
//	}(tx)
//
//	var user models.User
//	q := `
//	UPDATE users
//	SET balance = balance + $1
//	WHERE id = $2
//	RETURNING id, username, balance
//	`
//	logrus.Trace(formatQuery(q))
//	row := tx.QueryRow(q, amount, userID)
//	err = row.Scan(&user.Id, &user.Username, &user.Balance)
//	if err != nil {
//		_ = tx.Rollback()
//		logrus.Error(err)
//		if err.Error() == pqBalanceNegativeError {
//			return 0, domain.ErrAccessDenied
//		}
//		if errors.Is(err, sql.ErrNoRows) {
//			return 0, domain.ErrNotFound
//		}
//		return 0, err
//	}
//	err = tx.Commit()
//	if err != nil {
//		logrus.Error(err)
//		return 0, domain.ErrInternal
//	}
//	return user.Balance, nil
//}
//
//// TransferCurrency uses Serializable isolation level
//func (r localRepository) TransferCurrency(ctx context.Context, senderID uuid.UUID, receiverID uuid.UUID, amount float64) (newSenderBalance float64, err error) {
//	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
//	if err != nil {
//		return 0, domain.ErrInternal
//	}
//
//	// defer выполняется перед паникой
//	defer func(tx *sql.Tx) {
//		_ = tx.Rollback()
//	}(tx)
//
//	q := `
//		UPDATE users SET balance = balance - $1 WHERE id = $2 RETURNING balance;
//		`
//	logrus.Trace(formatQuery(q))
//	row := tx.QueryRow(q, amount, senderID)
//	err = row.Scan(&newSenderBalance)
//	if err != nil {
//		_ = tx.Rollback()
//		logrus.Error(err)
//		if err.Error() == pqBalanceNegativeError {
//			return 0, domain.ErrAccessDenied
//		}
//		if errors.Is(err, sql.ErrNoRows) {
//			return 0, domain.ErrNotFound
//		}
//		return 0, domain.ErrIncorrectBody
//	}
//
//	q = `
//		UPDATE users SET balance = balance + $1 WHERE id = $2 RETURNING balance;
//		`
//	logrus.Trace(formatQuery(q))
//	row = tx.QueryRow(q, amount, receiverID)
//	var receiverBalance float64
//	err = row.Scan(&receiverBalance)
//	if err != nil {
//		_ = tx.Rollback()
//		logrus.Error(err)
//		if err.Error() == pqBalanceNegativeError {
//			return 0, domain.ErrAccessDenied
//		}
//		if errors.Is(err, sql.ErrNoRows) {
//			return 0, domain.ErrNotFound
//		}
//		return 0, domain.ErrIncorrectBody
//	}
//	err = tx.Commit()
//	if err != nil {
//		logrus.Error(err)
//		return 0, domain.ErrInternal
//	}
//	return
//}
