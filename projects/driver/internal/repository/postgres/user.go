package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/models"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
)

var _ adapters.UserRepository = userRepository{}

const (
	pqBalanceNegativeError = "pq: new row for relation \"users\" violates check constraint \"balance_non_negative\""
)

/*
ВАЖНО! Для операций с изменением баланса необходимо использовать изоляцию транзакций,
для корректной работы CONSTRAINT balance_non_negative CHECK (balance >= 0)
*/
type userRepository struct {
	db *sqlx.DB
}

func NewDriverRepository(repos *sqlx.DB) adapters.UserRepository {
	return &userRepository{
		db: repos,
	}
}

func (r userRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	q := `
	INSERT INTO users (username, balance)
	VALUES ($1, $2)
	RETURNING *
	`
	logrus.Trace(formatQuery(q))
	var resUser models.User
	err := r.db.Get(&resUser, q, user.Name, 0)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrInternal
		}
		return domain.User{}, err
	}
	return resUser.ToDomain(), nil
}

func (r userRepository) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	q := `
	SELECT * FROM users WHERE id = $1
	`
	logrus.Trace(formatQuery(q))
	var user models.User
	err := r.db.Get(&user, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
	}
	return user.ToDomain(), nil
}

func (r userRepository) GetBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	q := `
	SELECT * FROM users WHERE id = $1
	`
	logrus.Trace(formatQuery(q))
	var user models.User
	err := r.db.Get(&user, q, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.ErrNotFound
		}
		return 0, err
	}
	return user.Balance, nil
}

// AddToBalance uses Serializable isolation level
func (r userRepository) AddToBalance(ctx context.Context, userID uuid.UUID, amount float64) (float64, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return 0, domain.ErrInternal
	}

	// Defer a rollback in case anything fails. (from go.dev)
	// defer выполняется перед паникой
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	var user models.User
	q := `
	UPDATE users
	SET balance = balance + $1
	WHERE id = $2
	RETURNING id, username, balance
	`
	logrus.Trace(formatQuery(q))
	row := tx.QueryRow(q, amount, userID)
	err = row.Scan(&user.Id, &user.Username, &user.Balance)
	if err != nil {
		_ = tx.Rollback()
		logrus.Error(err)
		if err.Error() == pqBalanceNegativeError {
			return 0, domain.ErrAccessDenied
		}
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.ErrNotFound
		}
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return 0, domain.ErrInternal
	}
	return user.Balance, nil
}

// TransferCurrency uses Serializable isolation level
func (r userRepository) TransferCurrency(ctx context.Context, senderID uuid.UUID, receiverID uuid.UUID, amount float64) (newSenderBalance float64, err error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return 0, domain.ErrInternal
	}

	// defer выполняется перед паникой
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	q := `
		UPDATE users SET balance = balance - $1 WHERE id = $2 RETURNING balance;
		`
	logrus.Trace(formatQuery(q))
	row := tx.QueryRow(q, amount, senderID)
	err = row.Scan(&newSenderBalance)
	if err != nil {
		_ = tx.Rollback()
		logrus.Error(err)
		if err.Error() == pqBalanceNegativeError {
			return 0, domain.ErrAccessDenied
		}
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.ErrNotFound
		}
		return 0, domain.ErrIncorrectBody
	}

	q = `
		UPDATE users SET balance = balance + $1 WHERE id = $2 RETURNING balance;
		`
	logrus.Trace(formatQuery(q))
	row = tx.QueryRow(q, amount, receiverID)
	var receiverBalance float64
	err = row.Scan(&receiverBalance)
	if err != nil {
		_ = tx.Rollback()
		logrus.Error(err)
		if err.Error() == pqBalanceNegativeError {
			return 0, domain.ErrAccessDenied
		}
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.ErrNotFound
		}
		return 0, domain.ErrIncorrectBody
	}
	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return 0, domain.ErrInternal
	}
	return
}
