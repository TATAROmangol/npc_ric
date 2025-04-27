package storage

import (
	"context"
	"database/sql"
	"forms/pkg/logger"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestStorage_Institution(t *testing.T){
	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)
	
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer pgContainer.Terminate(ctx)

	dsn, err := pgContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

}

func TestStorage_Mentor(t *testing.T){
	
}

func TestStorage_Form(t *testing.T){
	
}