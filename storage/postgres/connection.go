package postgres

import (
    "auth-service/config"
    "fmt"
    "log"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

func Connection() (*sqlx.DB, error) {
    cfg := config.Load()

    // Correcting the connection string format
    conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
        cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_NAME, cfg.DB_PASSWORD)
    fmt.Println(conn)
    db, err := sqlx.Open("postgres", conn)
    if err != nil {
        log.Println(err)
        return nil, err
    }
    if err := db.Ping(); err != nil {
        log.Println(err)
        return nil, err
    }
    return db, err
}
