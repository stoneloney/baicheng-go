package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"reflect"
)

var db *sql.DB

func GetDb() *sql.DB {
	if db != nil {
		return db
	}
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=disable", Cfg.Database.Host, Cfg.Database.User, Cfg.Database.Password, Cfg.Database.Port, Cfg.Database.Name)
	//建立连接
	db, err := sql.Open("mssql", connString)
	if err != nil {
		//log.Fatal("Open Connection failed:", err.Error())
		fmt.Println("Open Connection failed:", err.Error())
	}
	//defer db.Close()
	return db
}

func DbExec(queryStr string) (sql.Result, error) {
	if db == nil {
		db = GetDb()
	}
	stmt, err := db.Prepare(queryStr)
	if err != nil {
		Logger.Error(fmt.Sprintf("Prepare failed:%s", err.Error()))
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec()
	return result, err
}

func DbQuery(queryStr string, rowStruct interface{}) ([]interface{}, error) {
	if db == nil {
		db = GetDb()
	}
	stmt, err := db.Prepare(queryStr)
	if err != nil {
		//fmt.Println("Prepare failed:", err.Error())
		Logger.Error(fmt.Sprintf("Prepare failed:%s", err.Error()))
		return nil, err
	}
	defer stmt.Close()

	typ := reflect.TypeOf(rowStruct).Elem()

	var fieldNames []string
	for i := 0; i < typ.NumField(); i++ {
		if name, ok := typ.Field(i).Tag.Lookup("sql"); ok {
			fieldNames = append(fieldNames, name)
		}
	}
	if len(fieldNames) == 0 {
		Logger.Error(fmt.Sprintf("filed empty"))
		return nil, errors.New("field empty")
	}
	rows, err := stmt.Query()
	if err != nil {
		Logger.Error(fmt.Sprintf("query failed:%s", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var datas []interface{}
	for rows.Next() {
		row := reflect.New(typ)
		val := row.Elem()
		var fieldValues []interface{}
		for i := 0; i < typ.NumField(); i++ {
			if _, ok := typ.Field(i).Tag.Lookup("sql"); ok {
				field := val.Field(i)
				if field.Kind() != reflect.Ptr && field.CanAddr() {
					fieldValues = append(fieldValues, field.Addr().Interface())
				} else {
					fieldValues = append(fieldValues, field.Interface())
				}
			}
		}
		err = rows.Scan(fieldValues...)
		if err != nil {
			return nil, err
		}
		datas = append(datas, row.Interface())
	}
	return datas, nil
}
