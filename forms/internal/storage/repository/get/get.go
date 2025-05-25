package get

import (
	"context"
	"database/sql"
	"forms/internal/entities"
	"forms/internal/storage/repository/errors"
	"forms/pkg/logger"

	"github.com/lib/pq"
)

type Get struct {
	db *sql.DB
}

func NewGet(db *sql.DB) *Get {
	return &Get{db}
}

func (g *Get) GetInstitutions(ctx context.Context) ([]entities.Institution, error) {
	stmt, err := g.db.Prepare(`
	SELECT id, name, inn
	FROM institutions
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrCreateStatement, err)
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
			logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrScanRow, err)
			return nil, err
		}

		institutions = append(institutions, institution)
	}

	return institutions, nil
}

func (g *Get) GetMentors(ctx context.Context) ([]entities.Mentor, error) {
	stmt, err := g.db.Prepare(`
	SELECT *
	FROM mentors
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrCreateStatement, err)
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
		err = rows.Scan(&mentor.Id, &mentor.Name)
		if err != nil {
			logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrScanRow, err)
			return nil, err
		}
		mentors = append(mentors, mentor)
	}

	return mentors, nil
}

func (g *Get) GetInstitutionFromINN(ctx context.Context, inn int) (entities.Institution, error) {
	stmt, err := g.db.Prepare(`
	SELECT id, name, inn
	FROM institutions
	WHERE inn = $1
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrCreateStatement, err)
		return entities.Institution{}, err
	}
	defer stmt.Close()

	var institution entities.Institution
	err = stmt.QueryRow(inn).Scan(&institution.Id, &institution.Name, &institution.INN)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrScanRow, err)
		return entities.Institution{}, err
	}

	return institution, nil
}

func (g *Get) GetFormColumns(ctx context.Context, id int) ([]string, error) {
	stmt, err := g.db.Prepare(`
	SELECT columns
	FROM institutions
	WHERE id = $1
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrCreateStatement, err)
		return nil, err
	}
	defer stmt.Close()

	var columns []string
	err = stmt.QueryRow(id).Scan(pq.Array(&columns))
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrScanRow, err)
		return nil, err
	}

	return columns, nil
}

func (g *Get) GetFormRows(ctx context.Context, institution_id int) ([][]string, error) {
	stmt, err := g.db.Prepare(`
	SELECT info
	FROM forms
	WHERE institution_id = $1
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrCreateStatement, err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(institution_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([][]string, 0)
	for rows.Next() {
		var row []string
		err = rows.Scan(pq.Array(&row))
		if err != nil {
			logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrScanRow, err)
			return nil, err
		}

		logger.GetFromCtx(ctx).InfoContext(ctx, "row in storage", row)
		res = append(res, row)
	}

	return res, nil
}
