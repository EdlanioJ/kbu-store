package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EdlanioJ/kbu-store/domain"
)

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *accountRepository {
	return &accountRepository{
		db: db,
	}
}

func (r *accountRepository) Store(ctx context.Context, a *domain.Account) (err error) {
	query := `INSERT INTO accounts (id,created_at,updated_at,balance) VALUES ($1,$2,$3,$4)`
	res, err := r.db.ExecContext(ctx, query, a.ID, a.CreatedAt, a.UpdatedAt, a.Balance)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}

	if affect != 1 {
		err = fmt.Errorf("Weird Behavior. Total Affected: %d", affect)
		return
	}
	return
}

func (r *accountRepository) FindByID(ctx context.Context, id string) (res *domain.Account, err error) {
	query := `SELECT * FROM accounts WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	a := new(domain.Account)
	err = row.Scan(
		&a.ID,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.Balance,
	)
	if err != nil {
		return
	}

	res = a
	return
}

func (r *accountRepository) Update(ctx context.Context, a *domain.Account) (err error) {
	query := `UPDATE accounts SET created_at=$1,updated_at=$2,balance=$3 WHERE id = $4`
	res, err := r.db.ExecContext(ctx, query, a.CreatedAt, a.UpdatedAt, a.Balance, a.ID)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}

	if affect != 1 {
		err = fmt.Errorf("Weird Behavior. Total Affected: %d", affect)
		return
	}

	return
}

func (r *accountRepository) Delete(ctx context.Context, id string) (err error) {
	query := `DELETE FROM accounts WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}

	if affect != 1 {
		err = fmt.Errorf("Weird Behavior. Total Affected: %d", affect)
		return
	}

	return
}
