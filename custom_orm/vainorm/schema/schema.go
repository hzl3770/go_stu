package schema

import (
	"go/ast"
	"go_stu/custom_orm/vainorm/dialect"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	// 映射的对象
	Model interface{}
	// 表名
	Name string

	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {

	t := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     t.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < t.NumField(); i++ {
		p := t.Field(i)
		if p.Anonymous || !ast.IsExported(p.Name) {
			continue
		}

		field := &Field{
			Name: p.Name,
			Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
		}

		if v, ok := p.Tag.Lookup("vainorm"); ok {
			field.Tag = v
		}

		schema.Fields = append(schema.Fields, field)
		schema.FieldNames = append(schema.FieldNames, p.Name)
		schema.fieldMap[p.Name] = field
	}
	return schema
}

func (s *Schema) RecordValues(dest any) []any {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var values []any
	for _, name := range s.FieldNames {
		values = append(values, destValue.FieldByName(name).Interface())
	}
	return values
}
