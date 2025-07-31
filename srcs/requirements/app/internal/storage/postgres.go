package storage

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/sirupsen/logrus"
    "time"
)

func NewPostgres(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
    const maxAttempts = 5                
    const retryInterval = 2 * time.Second

    var pool *pgxpool.Pool
    var err error

    cfg, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        return nil, fmt.Errorf("ошибка парсинга строки подключения: %w", err)
    }
    logrus.Info("Строка подключения успешно распознана")


    for attempt := 1; attempt <= maxAttempts; attempt++ {
        logrus.WithFields(logrus.Fields{
            "attempt": attempt,
            "max":     maxAttempts,
        }).Info("Попытка подключения к базе данных")


        pool, err = pgxpool.NewWithConfig(ctx, cfg)
        if err != nil {
            logrus.WithError(err).WithFields(logrus.Fields{
                "attempt": attempt,
            }).Warn("Не удалось создать пул подключений")
        } else {

            pingErr := pool.Ping(ctx)
            if pingErr != nil {
                logrus.WithError(pingErr).WithFields(logrus.Fields{
                    "attempt": attempt,
                }).Warn("Не удалось подключиться к базе данных")
                err = pingErr
                pool.Close()
            } else {
                logrus.WithFields(logrus.Fields{
                    "attempt": attempt,
                }).Info("Успешное подключение к базе данных")
                return pool, nil
            }
        }
        if attempt < maxAttempts {
            logrus.WithField("retry_interval", retryInterval).Info("Повторная попытка подключения через")
            time.Sleep(retryInterval)
        }
    }
    return nil, fmt.Errorf("не удалось подключиться к базе данных после %d попыток: %w", maxAttempts, err)
}