package query

import (
	"database/sql"
	"fmt"
	suffix "github.com/doytowin/doyto-query-go-sql/field"
	"reflect"
	"strings"
)

type EntityMetadata[E any] struct {
	Type      reflect.Type
	TableName string
	Columns   []string
	ColStr    string
}

func IntPtr(o int) *int {
	return &o
}

func BuildConditions(query interface{}) (string, []any) {
	refType := reflect.TypeOf(query)
	rv := reflect.ValueOf(query)
	cnt, argCnt := 0, 0
	conditions := make([]string, refType.NumField())
	args := make([]any, refType.NumField(), 2*refType.NumField())
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		value := rv.FieldByName(field.Name)
		if IsValidValue(value) {
			conditions[cnt] = suffix.Process(field.Name)
			cnt++
			if value.Type().String() == "*int" {
				args[argCnt] = reflect.Indirect(value).Int()
				argCnt++
			}
		}
	}
	return strings.Join(conditions[0:cnt], " AND "), args[0:argCnt]
}

func IsValidValue(value reflect.Value) bool {
	if value.Type().Name() == "bool" {
		return value.Bool()
	} else {
		return !value.IsNil()
	}
}

func BuildEntityMetadata[E any](entity interface{}) EntityMetadata[E] {
	refType := reflect.TypeOf(entity)
	conditions := make([]string, refType.NumField())
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		conditions[i] = field.Name
	}
	return EntityMetadata[E]{
		Type:      refType,
		TableName: strings.TrimSuffix(refType.Name(), "Entity"),
		Columns:   conditions,
		ColStr:    strings.Join(conditions, ", "),
	}
}

func (em *EntityMetadata[E]) BuildSelect(query interface{}) (string, []any) {
	conditions, args := BuildConditions(query)
	s := "SELECT " + em.ColStr + " FROM " + em.TableName
	if len(conditions) > 0 {
		s += " WHERE " + conditions
	}
	fmt.Println("SQL: " + s)
	return s, args
}

func (em *EntityMetadata[E]) Query(db *sql.DB, query interface{}) ([]E, error) {
	sqlStr, args := em.BuildSelect(query)
	stmt, err := db.Prepare(sqlStr)
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	var result []E

	length := len(em.Columns)
	pointers := make([]interface{}, length)
	entity := reflect.New(em.Type).Elem().Interface().(E)
	elem := reflect.ValueOf(&entity).Elem()
	for i := range pointers {
		pointers[i] = elem.Field(i).Addr().Interface()
	}

	rows, _ := stmt.Query(args...)
	for rows.Next() {
		err := rows.Scan(pointers...)
		if err != nil {
			break
		}
		result = append(result, entity)
	}

	return result, err
}
