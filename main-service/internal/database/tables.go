package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type HandlerDB struct {
	conn   *pgx.Conn
	logger *logrus.Logger
}

func NewHandlerDB(conn *pgx.Conn, logger *logrus.Logger) *HandlerDB {
	return &HandlerDB{
		conn:   conn,
		logger: logger,
	}
}

func (h *HandlerDB) CreateTables(ctx context.Context, conn *pgx.Conn) error {
	tx, err := h.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("HandlerDB.CreateTables: %w", err)
	}
	defer tx.Rollback(ctx)

	creator := NewTableCreator()

	queries := []string{
		creator.createOrders(),
		creator.createDeliveries(),
		creator.createPayments(),
		creator.createItems(),
	}

	for _, query := range queries {
		_, err := tx.Exec(ctx, query)
		if err != nil {
			return fmt.Errorf("HandlerDB.CreateTables: %w", err)
		}
	}

	return tx.Commit(ctx)
}
