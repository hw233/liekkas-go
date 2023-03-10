package mysql

import (
	"context"
	"database/sql"
	"sync"
)

type EmbedModule struct {
	sync.Mutex

	db        *sql.DB    // 数据库连接
	table     *Table     // 表结构和表数据，表结构是根据反射映射的，不可以修改
	generator *Generator // sql生成
	// tableName string     // 表名称
}

func (m *EmbedModule) SetTable(table string) {
	m.Lock()
	defer m.Unlock()

	m.table.Name = table
}

func (m *EmbedModule) Table() string {
	return m.table.Name
}

func (m *EmbedModule) SetWhere(where string, args ...interface{}) {
	m.generator.Where(where, args...)
}

func (m *EmbedModule) init(db *sql.DB, ts *TableStruct) {
	m.Lock()
	defer m.Unlock()

	m.db = db
	m.table = NewTable(ts)
	m.generator = NewGenerator(m.table)
}

func (m *EmbedModule) needInit() bool {
	m.Lock()
	defer m.Unlock()

	if m.db == nil || m.table == nil {
		return true
	}

	return false
}

func (m *EmbedModule) create(ctx context.Context, data interface{}) error {
	m.Lock()
	defer m.Unlock()

	if !m.table.AssignableTo(data) {
		return ErrNotAssignableType
	}

	err := m.table.ReflectValue(data)
	if err != nil {
		return err
	}

	queryArgs, err := m.generator.GenInsertQueryArgs()
	if err != nil {
		return err
	}

	dbLog.Debugf("intert sql: %v", queryArgs.Query)
	for i, v := range queryArgs.Args {
		dbLog.Debugf("intert i: %d args: %s", i, string(v.([]byte)))
	}

	_, err = m.db.ExecContext(ctx, queryArgs.Query, queryArgs.Args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *EmbedModule) load(ctx context.Context, data interface{}) error {
	m.Lock()
	defer m.Unlock()

	if !m.table.AssignableTo(data) {
		return ErrNotAssignableType
	}

	err := m.table.ReflectValue(data)
	if err != nil {
		return err
	}

	queryArgs, err := m.generator.GenSelectQueryArgs()
	if err != nil {
		return err
	}

	dbLog.Debugf("select sql: %v", queryArgs.Query)
	for i, v := range queryArgs.Args {
		dbLog.Debugf("select i: %d args: %s", i, string(v.([]byte)))
	}

	row := m.db.QueryRowContext(ctx, queryArgs.Query, queryArgs.Args...)

	fieldsNum := m.table.FieldsNum()

	values := make([]interface{}, 0, fieldsNum)

	for i := 0; i < fieldsNum; i++ {
		field, err := m.table.FieldValue(i)
		if err != nil {
			return err
		}

		values = append(values, field.ValuePtr())
	}

	err = row.Scan(values...)
	if err != nil {
		return err
	}

	err = m.table.ReflectSet(data)
	if err != nil {
		dbLog.Errorf("load error: %v", err)
		return err
	}

	return err
}

func (m *EmbedModule) save(ctx context.Context, data interface{}) error {
	m.Lock()
	defer m.Unlock()

	if !m.table.AssignableTo(data) {
		return ErrNotAssignableType
	}

	err := m.table.ReflectValue(data)
	if err != nil {
		return err
	}

	if !m.table.NeedUpdate() {
		return nil
	}

	queryArgs, err := m.generator.GenUpdateQueryArgs()
	if err != nil {
		return err
	}

	dbLog.Debugf("update sql: %v", queryArgs.Query)
	for i, v := range queryArgs.Args {
		dbLog.Debugf("update i: %d args: %s", i, string(v.([]byte)))
	}

	_, err = m.db.ExecContext(ctx, queryArgs.Query, queryArgs.Args...)
	if err != nil {
		return err
	}

	return nil
}
