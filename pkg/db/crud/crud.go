package crud

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type crudRepository struct {
	db        *sqlx.DB
	tableName string
}

func NewCrudRepository(db *sqlx.DB, tableName string) CRUD {
	return &crudRepository{
		db:        db,
		tableName: tableName,
	}
}

func (c *crudRepository) create(obj reflect.Value) error {
	columns := ""
	values := ""
	var args []interface{}
	k := int(0)
	for i := 0; i < obj.Type().NumField(); i++ {
		if val, ok := obj.Type().Field(i).Tag.Lookup("db"); ok {
			valueField := obj.Field(i)
			if valueField.Kind() == reflect.Interface && !valueField.IsNil() {
				elm := valueField.Elem()
				if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
					valueField = elm
				}
			}
			if valueField.Kind() == reflect.Ptr {
				valueField = valueField.Elem()
			}
			var valU interface{}
			var currentValue *interface{}
			if valueField.IsValid() {
				valU = valueField.Interface()
				currentValue = &valU
			} else {
				currentValue = nil
				continue
			}
			columns += fmt.Sprintf("%s, ", val)
			args = append(args, currentValue)
			values += "$" + strconv.Itoa(k+1) + ", "
			k += 1
		}
	}
	columns = columns[:len(columns)-2]
	values = values[:len(values)-2]
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", c.tableName, columns, values)
	_, err := c.db.Exec(query, args...)
	return err
}

func (c *crudRepository) Create(data interface{}) error {
	return c.create(reflect.Indirect(reflect.ValueOf(data)))
}

func (c *crudRepository) Get(filters map[string]interface{}) (interface{}, error) {
	var data []map[string]interface{}
	var query string
	query += " WHERE "
	for key, value := range filters {
		switch value.(type) {
		case string:
			query += fmt.Sprintf("%s='%s' AND ", key, value)
		default:
			query += fmt.Sprintf("%s=%v AND ", key, value)
		}
	}
	query = query[:len(query)-5]
	query = fmt.Sprintf("SELECT * FROM %s %s limit 1", c.tableName, query)
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	colNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	cols := make([]interface{}, len(colNames))
	colPtrs := make([]interface{}, len(colNames))
	for i := 0; i < len(colNames); i++ {
		colPtrs[i] = &cols[i]
	}
	k := 0
	for rows.Next() {
		data = append(data, make(map[string]interface{}))
		err = rows.Scan(colPtrs...)
		if err != nil {
			return nil, err
		}
		for i, col := range cols {
			data[k][colNames[i]] = col
		}
		k++
	}
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	return data[0], nil

}

func (c *crudRepository) Fetch(filters map[string]interface{}, order *string, orderType *string, limit *int64) ([]interface{}, error) {
	var data []map[string]interface{}
	var dataResponse []interface{}
	var query string
	query += " WHERE "
	for key, value := range filters {
		switch value.(type) {
		case string:
			query += fmt.Sprintf("%s='%s' AND ", key, value)
		default:
			query += fmt.Sprintf("%s=%v AND ", key, value)
		}
	}
	query = query[:len(query)-5]
	query = fmt.Sprintf("SELECT * FROM %s %s", c.tableName, query)
	if order != nil {
		query += query + " order by " + *order
	}
	if orderType != nil {
		query += query + " " + *orderType
	}
	if limit != nil {
		query += query + " limit " + strconv.Itoa(int(*limit))
	}
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	colNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	cols := make([]interface{}, len(colNames))
	colPtrs := make([]interface{}, len(colNames))
	for i := 0; i < len(colNames); i++ {
		colPtrs[i] = &cols[i]
	}
	k := 0
	for rows.Next() {
		data = append(data, make(map[string]interface{}))
		err = rows.Scan(colPtrs...)
		if err != nil {
			return nil, err
		}
		for i, col := range cols {
			data[k][colNames[i]] = col
		}
		dataResponse = append(dataResponse, data[k])
		k++
	}
	if err != nil {
		return nil, err
	}
	if len(dataResponse) == 0 {
		return nil, nil
	}
	return dataResponse, nil
}

func (c *crudRepository) Update(filters map[string]interface{}, units map[string]interface{}) error {
	var query string
	query += fmt.Sprintf("update %s set ", c.tableName)
	for key, value := range units {
		switch value.(type) {
		case string:
			query += fmt.Sprintf("%s='%s', ", key, value)
		default:
			query += fmt.Sprintf("%s=%v, ", key, value)
		}
	}
	query = query[:len(query)-2]
	query += " WHERE "
	for key, value := range filters {
		switch value.(type) {
		case string:
			query += fmt.Sprintf("%s='%s' AND ", key, value)
		default:
			query += fmt.Sprintf("%s=%v AND ", key, value)
		}
	}
	query = query[:len(query)-5]
	_, err := c.db.Exec(query)
	return err
}

func (c *crudRepository) Delete(filters map[string]interface{}) error {
	var query string
	query += fmt.Sprintf("delete from %s ", c.tableName)
	query += " WHERE "
	for key, value := range filters {
		switch value.(type) {
		case string:
			query += fmt.Sprintf("%s='%s' AND ", key, value)
		default:
			query += fmt.Sprintf("%s=%v AND ", key, value)
		}
	}
	query = query[:len(query)-5]
	_, err := c.db.Exec(query)
	return err

}
