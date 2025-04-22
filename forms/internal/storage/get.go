package storage

import (
	"context"
	"forms/internal/entities"
	"forms/pkg/logger"
)

func (s *Storage) GetInstitutions(ctx context.Context) ([]entities.Institution, error){
	stmt, err := s.db.Prepare(`
	SELECT id, name, inn
	FROM institutions
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrCreateStatement, err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var institutions []entities.Institution
	for rows.Next() {
		var institution entities.Institution
		err = rows.Scan(&institution.Id, &institution.Name, &institution.INN)
		if err != nil {
			logger.GetFromCtx(ctx).ErrorContext(ctx, ErrScanRow, err)
			return nil, err
		}
		institutions = append(institutions, institution)
	}

	return institutions, nil
}

func (s *Storage) GetMentors(ctx context.Context) ([]entities.Mentor, error){
	stmt, err := s.db.Prepare(`
	SELECT id, info
	FROM mentors
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrCreateStatement, err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mentors []entities.Mentor
	for rows.Next() {
		var mentor entities.Mentor
		err = rows.Scan(&mentor.Id, &mentor.Info)
		if err != nil {
			logger.GetFromCtx(ctx).ErrorContext(ctx, ErrScanRow, err)
			return nil, err
		}
		mentors = append(mentors, mentor)
	}

	return mentors, nil
}

func (s *Storage) GetInstitutionFromINN(ctx context.Context, inn int) (entities.Institution, error){
	stmt, err := s.db.Prepare(`
	SELECT id, name, inn
	FROM institutions
	WHERE inn = $1
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrCreateStatement, err)
		return entities.Institution{}, err
	}
	defer stmt.Close()

	var institution entities.Institution
	err = stmt.QueryRow(inn).Scan(&institution.Id, &institution.Name, &institution.INN)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrScanRow, err)
		return entities.Institution{}, err
	}

	return institution, nil
}

func (s *Storage) GetFormColumns(ctx context.Context, id int) ([]string, error){
	stmt, err := s.db.Prepare(`
	SELECT columns
	FROM institutions
	WHERE id = $1
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrCreateStatement, err)
		return nil, err
	}
	defer stmt.Close()

	var columns []string
	err = stmt.QueryRow(id).Scan(&columns)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrScanRow, err)
		return nil, err
	}

	return columns, nil
}