package post

import (
	"context"
	"database/sql"
	"forms/internal/storage/repository/errors"
	"forms/pkg/logger"

	"github.com/lib/pq"
)

type Post struct {
	db *sql.DB
}

func NewPost(db *sql.DB) *Post {
	return &Post{db}
}

func (p *Post) PostInstitution(ctx context.Context, name string, inn int, columns []string) (int, error) {
	stmt, err := p.db.Prepare(`
	INSERT INTO institutions (name, inn, columns)
	VALUES ($1, $2, $3) 
	RETURNING id
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrCreateStatement, err)
		return -1, err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRowContext(ctx, name, inn, pq.Array(columns)).Scan(&id)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrScanRow, err)
		return -1, err
	}
	return id, nil
}

func (p *Post) PostMentor(ctx context.Context, name string) (int, error) {
	stmt, err := p.db.Prepare(`
	INSERT INTO mentors (name) 
	VALUES ($1) 
	RETURNING id
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrCreateStatement, err)
		return -1, err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRowContext(ctx, name).Scan(&id)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrScanRow, err)
		return -1, err
	}
	return id, nil
}

func (p *Post) PostForm(ctx context.Context, info []string, institutionId int, mentorId int) (int, error) {
	stmt, err := p.db.Prepare(`
	INSERT INTO forms (info, institution_id, mentor_id) 
	VALUES ($1, $2, $3) 
	RETURNING id
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrCreateStatement, err)
		return -1, err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRowContext(ctx, pq.Array(info), institutionId, mentorId).Scan(&id)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, errors.ErrScanRow, err)
		return -1, err
	}
	return id, nil
}