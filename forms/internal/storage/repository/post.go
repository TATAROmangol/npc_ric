package repository

import (
	"context"
	"forms/pkg/logger"
)

func (s *Storage) PostInstitution(ctx context.Context, name string, inn int, columns []string) (int, error) {
	stmt, err := s.db.Prepare(`
	INSERT INTO institutions (name, inn) 
	VALUES ($1, $2) 
	RETURNING id
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrCreateStatement, err)
		return -1, err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRowContext(ctx, name, inn).Scan(&id)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrScanRow, err)
		return -1, err
	}
	return id, nil
}

func (s *Storage) PostMentor(ctx context.Context, name string) (int, error) {
	stmt, err := s.db.Prepare(`
	INSERT INTO mentors (info) 
	VALUES ($1) 
	RETURNING id
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrCreateStatement, err)
		return -1, err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRowContext(ctx, name).Scan(&id)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrScanRow, err)
		return -1, err
	}
	return id, nil
}

func (s *Storage) PostForm(ctx context.Context, institutionId int, info []string) (int, error) {
	stmt, err := s.db.Prepare(`
	INSERT INTO forms (institution_id, info) 
	VALUES ($1, $2) 
	RETURNING id
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrCreateStatement, err)
		return -1, err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRowContext(ctx, institutionId, info).Scan(&id)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrScanRow, err)
		return -1, err
	}
	return id, nil
}