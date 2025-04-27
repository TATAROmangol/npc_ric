package repository

import (
	"context"
	"forms/pkg/logger"
)

func (s *Storage) DeleteInstitution(ctx context.Context, institutionId int) error {
	stmt, err := s.db.Prepare(`
	DELETE FROM institutions
	WHERE id = $1
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrCreateStatement, err)
		return err 
	}
	defer stmt.Close()

	_, err = stmt.Exec(institutionId)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrExecStatement, err)
		return err
	}
	return nil
}

func (s *Storage) DeleteMentor(ctx context.Context, mentorId int) error {
	stmt, err := s.db.Prepare(`
	DELETE FROM mentors
	WHERE id = $1
	`)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrCreateStatement, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(mentorId)
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, ErrExecStatement, err)
		return err
	}
	return nil
}