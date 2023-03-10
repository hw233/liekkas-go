package mysql

import (
	"bytes"
)

type Generator struct {
	bytes.Buffer
	args      []interface{}
	table     *Table
	whereArgs *QueryArgs
}

func NewGenerator(table *Table) *Generator {
	return &Generator{
		Buffer:    bytes.Buffer{},
		args:      []interface{}{},
		table:     table,
		whereArgs: nil,
	}
}

func (g *Generator) Clear() {
	g.Buffer = bytes.Buffer{}
	g.args = []interface{}{}
}

func (g *Generator) AddArg(arg interface{}) {
	g.args = append(g.args, arg)
}

func (g *Generator) GenQueryArgs() *QueryArgs {
	return NewQueryArgs(g.String(), g.args)
}

func (g *Generator) Where(where string, args ...interface{}) {
	g.whereArgs = NewQueryArgs(where, args)
}

func (g *Generator) writeWhere() {
	g.WriteString(" WHERE ")

	if g.whereArgs == nil {
		// from major
		majorField := g.table.MajorField()

		g.WriteByte('`')
		g.WriteString(majorField.Name)
		g.WriteString("`=?")

		g.AddArg(majorField.Value)

		return
	}

	// from where
	g.WriteString(g.whereArgs.Query)
	g.AddArg(g.whereArgs.Args)
}

func (g *Generator) GenUpdateQueryArgs() (*QueryArgs, error) {
	defer g.Clear()

	g.WriteString("UPDATE `")
	g.WriteString(g.table.Name)
	g.WriteString("` SET ")

	comma := false

	// write fields
	for i := 0; i < g.table.FieldsNum(); i++ {
		field, err := g.table.Field(i)
		if err != nil {
			return nil, err
		}

		// no need update or be major key
		if !field.NeedUpdate || i == g.table.MajorIndex {
			continue
		}

		if comma {
			g.WriteByte(',')
		}

		g.WriteByte('`')
		g.WriteString(field.Name)
		g.WriteString("`=?")

		g.AddArg(field.Value)

		field.RefreshUpdate()

		comma = true
	}

	// write where
	g.writeWhere()

	return g.GenQueryArgs(), nil
}

func (g *Generator) GenSelectQueryArgs() (*QueryArgs, error) {
	defer g.Clear()

	g.WriteString("SELECT ")

	// write fields
	for i := 0; i < g.table.FieldsNum(); i++ {
		field, err := g.table.FieldStruct(i)
		if err != nil {
			return nil, err
		}

		if i > 0 {
			g.WriteByte(',')
		}

		g.WriteByte('`')
		g.WriteString(field.Name)
		g.WriteByte('`')
	}

	g.WriteString(" FROM ")
	g.WriteByte('`')
	g.WriteString(g.table.Name)
	g.WriteByte('`')

	// write where
	g.writeWhere()

	return g.GenQueryArgs(), nil
}

func (g *Generator) GenInsertQueryArgs() (*QueryArgs, error) {
	defer g.Clear()

	g.WriteString("INSERT INTO `")
	g.WriteString(g.table.Name)
	g.WriteString("` (")

	valuesBuf := bytes.Buffer{}

	for i := 0; i < g.table.FieldsNum(); i++ {
		field, err := g.table.Field(i)
		if err != nil {
			return nil, err
		}

		if i >= 1 {
			g.WriteByte(',')
			valuesBuf.WriteByte(',')
		}

		g.WriteByte('`')
		g.WriteString(field.Name)
		g.WriteByte('`')

		valuesBuf.WriteByte('?')

		g.AddArg(field.Value)

		field.RefreshUpdate()
	}

	g.WriteString(") VALUES (")
	g.Write(valuesBuf.Bytes())
	g.WriteByte(')')

	return g.GenQueryArgs(), nil
}
