package testutils


import (
    "context"
    "fmt"
	"os"
    "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
	"walletT/internal/storage"
	"github.com/sirupsen/logrus"
    "log"
)

func CreateTestDatabase(ctx context.Context, adminDSN, testDBName string) (string, error) {

	log.Println("Попытка подключения к базе данных...")

	db, err := storage.NewPostgres(ctx, adminDSN)
	if err != nil{
		logrus.WithError(err).Fatal("Ошибка подключения к базе данных")
	}
	logrus.Info("Подключение к базе данных успешно")

    _, _ = db.Exec(ctx, fmt.Sprintf(`
        SELECT pg_terminate_backend(pid)
        FROM pg_stat_activity
        WHERE datname = '%s' AND pid <> pg_backend_pid()`, testDBName))

    _, _ = db.Exec(ctx, fmt.Sprintf(`DROP DATABASE IF EXISTS %s`, testDBName))

    _, err = db.Exec(ctx, fmt.Sprintf(`CREATE DATABASE %s`, testDBName))
    if err != nil {
        return "", fmt.Errorf("не удалось создать тестовую базу: %w", err)
    }
	user := os.Getenv("DB_USER")
    pass :=  os.Getenv("DB_PASS")
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")

    testDSN := "postgres://" + user + ":" + pass+"@"+host+":"+port+"/"+testDBName+"?sslmode=disable"
    logrus.Infof("Тестовая база %s создана", testDBName)
    return testDSN, nil
}

func ApplyMigrations(migrationsPath, dbURL string) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("ошибка при создании мигратора: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("ошибка применения миграций: %w", err)
	}

	logrus.Info("Миграции успешно применены")
	return nil
}