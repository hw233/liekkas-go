package mysql

import (
	"context"
	"database/sql"
)

type Module interface {
	Table() string

	needInit() bool
	init(db *sql.DB, table *TableStruct)

	create(ctx context.Context, data interface{}) error
	load(ctx context.Context, data interface{}) error
	save(ctx context.Context, data interface{}) error
}
