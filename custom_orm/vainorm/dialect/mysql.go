package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type mysql struct {
}

var _ Dialect = (*mysql)(nil)

func init() {
	RegisterDialect("mysql", &mysql{})
}

func (m *mysql) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return "int"
	case reflect.Int64:
		return "bigint"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "int unsigned"
	case reflect.Uint64:
		return "bigint unsigned"
	case reflect.Float32, reflect.Float64:
		return "double"
	case reflect.String:
		return "varchar(255)"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}

	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

func (m *mysql) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT table_name FROM information_schema.tables WHERE table_name = ?", args
}
