package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"

    _ "github.com/ClickHouse/clickhouse-go/v2"
)

var DB *sql.DB

func InitDB() {
    dsn := fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s",
        os.Getenv("CLICKHOUSE_USER"),
        os.Getenv("CLICKHOUSE_PASSWORD"),
        os.Getenv("CLICKHOUSE_HOST"),
        os.Getenv("CLICKHOUSE_PORT"),
        os.Getenv("CLICKHOUSE_DATABASE"),
    )

    var err error
    DB, err = sql.Open("clickhouse", dsn)
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := DB.PingContext(ctx); err != nil {
        log.Fatalf("DB ping failed: %v", err)
    }

    log.Println("✅ ClickHouse connected")
}
