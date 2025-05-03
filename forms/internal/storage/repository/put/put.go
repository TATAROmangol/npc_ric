package put

import (
	"context"
	"database/sql"
	"forms/internal/storage/repository/errors"
	"forms/pkg/logger"

	"github.com/lib/pq"
)

type Put struct {
	db *sql.DB
}
func NewPut(db *sql.DB) *Put {
	return &Put{db}
}

func (p *Put) PutInstitutionInfo(ctx context.Context, id int, name string, inn int) error {
	stmt, err := p.db.Prepare(`
	UPDATE institutions 
	SET name = $2, 
		inn = $3 
	WHERE id = $1
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrCreateStatement, err)
		return err 
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, name, inn)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrExecStatement, err)
		return err
	}
	return nil
}

func (p *Put) PutInstitutionColumns(ctx context.Context, id int, columns []string) error {
	stmt, err := p.db.Prepare(`
	UPDATE institutions 
	SET columns = $2
	WHERE id = $1
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrCreateStatement, err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, pq.Array(columns))
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrExecStatement, err)
		return err 
	}
	return nil
}

func (p *Put) PutMentor(ctx context.Context, id int, name string) error {
	stmt, err := p.db.Prepare(`
	UPDATE mentors 
	SET name = $2
	WHERE id = $1
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrCreateStatement, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, name)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrExecStatement, err)
		return err
	}
	return nil
}