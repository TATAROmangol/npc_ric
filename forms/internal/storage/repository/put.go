package repository

import (
	"context"
	"forms/pkg/logger"
)

func (s *Storage) PutInstitutionInfo(ctx context.Context, id int, name string, inn int) error {
	stmt, err := s.db.Prepare(`
	UPDATE institutions 
	SET name = $2, 
		inn = $3 
	WHERE id = $1
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrCreateStatement, err)
		return err 
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, inn, id)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrExecStatement, err)
		return err
	}
	return nil
}

func (s *Storage) PutInstitutionColumns(ctx context.Context, id int, columns []string) error {
	stmt, err := s.db.Prepare(`
	UPDATE institutions 
	SET columns = $2
	WHERE id = $1
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrCreateStatement, err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(columns, id)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrExecStatement, err)
		return err 
	}
	return nil
}

func (s *Storage) PutMentor(ctx context.Context, id int, info string) error {
	stmt, err := s.db.Prepare(`
	UPDATE mentors 
	SET info = $2
	WHERE id = $1
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrCreateStatement, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(info, id)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrExecStatement, err)
		return err
	}
	return nil
}