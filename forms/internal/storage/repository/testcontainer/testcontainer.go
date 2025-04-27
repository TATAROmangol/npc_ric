package testcontainer

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	testcontainersPostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	_ "github.com/lib/pq"
)

const (
	migrationPath = "file://../../migrations"
)

func SetupTestDB() (*sql.DB, func(), error) {
	ctx := context.Background()
	
	// Запускаем контейнер с PostgreSQL
	container, err := testcontainersPostgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		testcontainersPostgres.WithDatabase("testdb"),
		testcontainersPostgres.WithUsername("user"),
		testcontainersPostgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start container: %w", err)
	}

	// Получаем строку подключения
	connStr, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	// Добавляем параметры
	connStr = connStr + " sslmode=disable"

	// Подключаемся к БД
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Применяем миграции
	if err := applyMigrations(db); err != nil {
		return nil, nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Функция очистки
	cleanup := func() {
		_ = db.Close()
		if err := container.Terminate(ctx); err != nil {
			fmt.Printf("failed to terminate container: %s\n", err)
		}
	}

	return db, cleanup, nil
}

func applyMigrations(db *sql.DB) error {
	// Создаем драйвер для миграций
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Инициализируем миграции
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres", 
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	// Применяем все миграции вверх
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}