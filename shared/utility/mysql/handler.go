package mysql

import (
	"context"
	"database/sql"
	"reflect"
	"sync"

	_ "github.com/go-sql-driver/mysql"

	"shared/utility/naming"
)

// 复用TableStruct，key是Table.Name
var (
	tables = &sync.Map{}
	mutex  sync.Mutex
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	// db, err := sql.Open("mysql", config.Addr)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// db.SetMaxIdleConns(config.MaxIdleConn)
	// db.SetMaxOpenConns(config.MaxOpenConn)
	// db.SetConnMaxLifetime(config.ConnMaxLifetime)
	return &Handler{
		db: db,
	}
}

// func NewHandler(config *Config) (*Handler, error) {
// 	db, err := sql.Open("mysql", config.Addr)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	db.SetMaxIdleConns(config.MaxIdleConn)
// 	db.SetMaxOpenConns(config.MaxOpenConn)
// 	db.SetConnMaxLifetime(config.ConnMaxLifetime)
//
// 	return &Handler{
// 		db: db,
// 	}, db.Ping()
// }

// type initOpts struct {
// 	Table string
// 	Data  Module
// }

func (h *Handler) init(data Module) error {
	if data.needInit() {
		// name := data.Table()
		// if name == "" {
		name := naming.UnderlineNaming(reflect.TypeOf(data).Elem().Name())
		// }

		var table *TableStruct

		iTable, ok := tables.Load(name)
		if !ok {
			mutex.Lock()

			// 取2次，防止重复初始化
			iTable, ok = tables.Load(name)
			if !ok {
				table = NewTableStruct(name)

				err := table.ReflectTable(data)
				if err != nil {
					mutex.Unlock()
					return err
				}

				tables.Store(table.Name, table)
			} else {
				table, _ = iTable.(*TableStruct)
			}

			mutex.Unlock()
		} else {
			table, _ = iTable.(*TableStruct)
		}

		data.init(h.db, table)
	}

	return nil
}

//
// func (m *Handler) BatchLoad(data []Module) error {
// 	rows, err := m.db.Query("SELECT * FROM user", 1)
// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()
//
// 	for rows.Next() {
// 		rows.Scan()
// 	}
//
// 	err := m.checkInit(data)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

func (h *Handler) Preload(data Module) error {
	name := data.Table()
	if name == "" {
		name = naming.UnderlineNaming(reflect.TypeOf(data).Elem().Name())
	}

	_, ok := tables.Load(name)
	if !ok {
		mutex.Lock()

		// 取2次，防止重复初始化
		_, ok = tables.Load(name)
		if !ok {
			table := NewTableStruct(name)

			err := table.ReflectTable(data)
			if err != nil {
				return err
			}

			tables.Store(table.Name, table)
		}

		mutex.Unlock()
	}

	return nil
}

func (h *Handler) Init(data Module) error {
	err := h.init(data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) Create(ctx context.Context, data Module) error {
	err := h.init(data)
	if err != nil {
		return err
	}

	return data.create(ctx, data)
}

func (h *Handler) Load(ctx context.Context, data Module) error {
	err := h.init(data)
	if err != nil {
		return err
	}

	return data.load(ctx, data)
}

func (h *Handler) Save(ctx context.Context, data Module) error {
	err := h.init(data)
	if err != nil {
		return err
	}

	return data.save(ctx, data)
}

func (h *Handler) Close() error {
	return h.db.Close()
}
