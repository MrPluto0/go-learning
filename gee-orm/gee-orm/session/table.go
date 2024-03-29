package session

import (
	"example/gee-orm/log"
	"example/gee-orm/schema"
	"fmt"
	"reflect"
	"strings"
)

func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		s := fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag)
		columns = append(columns, s)
	}
	desc := strings.Join(columns, ",")
	sql := fmt.Sprintf("CREATE TABLE %s (%s)", table.Name, desc)
	_, err := s.Raw(sql).Exec()
	return err
}

func (s *Session) DropTable() error {
	sql := fmt.Sprintf("Drop TABLE IF EXISTS %s", s.RefTable().Name)
	_, err := s.Raw(sql).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
}
