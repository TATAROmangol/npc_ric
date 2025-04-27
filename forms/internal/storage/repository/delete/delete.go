package delete

import (
	"context"
	"database/sql"
	"forms/internal/storage/repository/errors"
	"forms/pkg/logger"
)

type Delete struct{
	db *sql.DB
}

func NewDelete(db *sql.DB) *Delete {
	return &Delete{db}
}

func (d *Delete) DeleteInstitution(ctx context.Context, institutionId int) error {
	stmt, err := d.db.Prepare(`
	DELETE FROM institutions
	WHERE id = $1
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrCreateStatement, err)
		return err 
	}
	defer stmt.Close()

	_, err = stmt.Exec(institutionId)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrExecStatement, err)
		return err
	}
	return nil
}

func (d *Delete) DeleteMentor(ctx context.Context, mentorId int) error {
	stmt, err := d.db.Prepare(`
	DELETE FROM mentors
	WHERE id = $1
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrCreateStatement, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(mentorId)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrExecStatement, err)
		return err
	}
	return nil
}