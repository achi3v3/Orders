package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func CreateTables(ctx context.Context, conn *pgx.Conn) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("subs.CreateTables: %w", err)
	}
	defer tx.Rollback(ctx)

	queries := []string{
		`CREATE TABLE IF NOT EXISTS orders (
			order_uid VARCHAR(255) PRIMARY KEY NOT NULL,
			track_number VARCHAR(255) NOT NULL UNIQUE,
			entry VARCHAR(20) NOT NULL,
			locale VARCHAR(10) NOT NULL,
			internal_signature VARCHAR(255) DEFAULT '',
			customer_id VARCHAR(255) NOT NULL,
			delivery_service VARCHAR(100) NOT NULL,
			shardkey VARCHAR(20) NOT NULL,
			sm_id INTEGER NOT NULL,
			date_created TIMESTAMPTZ NOT NULL,
			oof_shard VARCHAR(20) NOT NULL
	);`,
		`CREATE TABLE IF NOT EXISTS deliveries (
			order_uid VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			phone VARCHAR(50) NOT NULL,
			zip VARCHAR(50) NOT NULL,
			city VARCHAR(255) NOT NULL,
			address TEXT NOT NULL,
			region VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL
	);`,
		`CREATE TABLE IF NOT EXISTS payments (
			transaction VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
			request_id VARCHAR(255) DEFAULT '',
			currency VARCHAR(20) NOT NULL,
			provider VARCHAR(150) NOT NULL,
			amount INTEGER NOT NULL,
			payment_dt BIGINT NOT NULL,
			bank VARCHAR(150) NOT NULL,
			delivery_cost INTEGER NOT NULL,
			goods_total INTEGER NOT NULL,
			custom_fee INTEGER DEFAULT 0
	);`,
		`CREATE TABLE IF NOT EXISTS items (
			id SERIAL PRIMARY KEY,
			chrt_id BIGINT NOT NULL,
			track_number VARCHAR(255) REFERENCES orders(track_number) ON DELETE CASCADE,
			price INTEGER NOT NULL,
			rid VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			sale INTEGER NOT NULL,
			size VARCHAR(20) DEFAULT '0',
			total_price INTEGER NOT NULL,
			nm_id INTEGER NOT NULL,
			brand VARCHAR(255) NOT NULL,
			status INTEGER NOT NULL
	);`,
	}
	for _, query := range queries {
		_, err := tx.Exec(ctx, query)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return tx.Commit(ctx)
}
