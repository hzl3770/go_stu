package session

import (
	"go_stu/custom_orm/vainorm/clause"
	"reflect"
)

func (s *Session) Insert(values ...any) (int64, error) {
	if len(values) == 0 {
		return 0, nil
	}
	// 找出对象的表
	table := s.Model(values[0]).RefTable()
	// 构造insert语句
	s.clause.Set(clause.INSERT, table.Name, table.FieldNames)

	// 构造values
	// [[v1], [v2]]
	recordValues := make([]interface{}, 0)
	for _, v := range values {
		recordValues = append(recordValues, table.RecordValues(v))
	}

	// 设置values
	s.clause.Set(clause.VALUES, recordValues...)

	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	exec, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	return exec.RowsAffected()
}

func (s *Session) Find(values interface{}) error {
	// 获取slice的元素类型
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()

	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)

	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var vals []interface{}
		for _, name := range table.FieldNames {
			vals = append(vals, dest.FieldByName(name).Addr().Interface())
		}
		if err = rows.Scan(vals...); err != nil {
			return err
		}

		destSlice.Set(reflect.Append(destSlice, dest))
	}

	return rows.Close()
}
