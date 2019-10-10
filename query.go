package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"playground/asciitable"
	"strings"

	"time"

	"strconv"

	//_ "github.com/alexbrainman/odbc"
	//_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

type params struct {
	ConnectionType   string `json:"connection_type"`
	ConnectionString string `json:"connection_string"`
}

func newParams() *params {
	return new(params)
}

func main() {
	table := asciitable.New()

	table.SetTitle("Query")
	params := readParameters()
	sql := readSql()
	db := connectDb(params)
	rows := query(db, sql)
	defer rows.Close()
	setHeaders(table, rows)
	setData(table, rows)
	fmt.Println(table.String())
}

func readParameters() *params {
	var params = newParams()
	raw, err := ioutil.ReadFile("/Users/dsearle/go/src/playground/querytool/params.json")
	if err != nil {
		fmt.Println("error reading params.json: " + err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(raw, params)
	if err != nil {
		fmt.Println("error while unmarshalling: " + err.Error())
		os.Exit(1)
	}
	fmt.Printf("params: %v\n", params)
	return params
}

func readSql() string {
	raw, err := ioutil.ReadFile("/Users/dsearle/go/src/playground/querytool/sql.txt")
	if err != nil {
		fmt.Println("error reading sql.txt: " + err.Error())
		os.Exit(1)
	}
	return string(raw)
}

func connectDb(p *params) *sql.DB {
	db, err := sql.Open(p.ConnectionType, p.ConnectionString)
	if err != nil {
		fmt.Println("error creating db connection: " + err.Error())
		os.Exit(1)
	}
	return db
}

func query(db *sql.DB, sql string) *sql.Rows {
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println("error executing query: " + err.Error())
		os.Exit(1)
	}
	return rows
}

func setHeaders(table *asciitable.Table, rows *sql.Rows) {
	cTypes, err := rows.ColumnTypes()
	if err != nil {
		fmt.Println("error getting column types: " + err.Error())
		os.Exit(1)
	}

	headers := make([]string, 0)
	for _, col := range cTypes {
		headers = append(headers, strings.ToUpper(col.Name()+" ("+col.DatabaseTypeName()+")"))
	}
	table.SetHeaders(headers)
}

func setData(table *asciitable.Table, rows *sql.Rows) {
	vals := createDestinationSlices(rows)
	for rows.Next() {
		err := rows.Scan(vals...)
		if err != nil {
			fmt.Println("error scanning data: " + err.Error())
			os.Exit(1)
		}
		data := make([]string, 0)
		for _, val := range vals {
			switch val.(type) {
			case *time.Time:
				v := val.(*time.Time)
				data = append(data, v.Format(time.RFC3339))
			case *sql.NullString:
				v := val.(*sql.NullString)
				if v.Valid {
					data = append(data, v.String)
				} else {
					data = append(data, "NULL")
				}
			case *sql.NullInt64:
				d := int64(0)
				v := val.(*sql.NullInt64)
				if v.Valid {
					d = v.Int64
				}
				data = append(data, strconv.FormatInt(d, 10))
			default:
				fmt.Printf("unknown value: %v\n", val)
			}
		}
		table.AddRow(data)
	}
}

func createDestinationSlices(rows *sql.Rows) []interface{} {
	ctypes, err := rows.ColumnTypes()
	rval := make([]interface{}, 0)
	if err != nil {
		fmt.Println("error getting column types: " + err.Error())
		os.Exit(1)
	}
	for _, ctype := range ctypes {
		switch ctype.DatabaseTypeName() {
		case "DATETIME":
			rval = append(rval, new(time.Time))
		case "TINYINT", "SMALLINT", "INT":
			rval = append(rval, new(sql.NullInt64))
		case "CHAR", "VARCHAR", "NCHAR", "NVARCHAR":
			rval = append(rval, new(sql.NullString))
		default:
			rval = append(rval, new(interface{}))
			fmt.Println("unknown database type: ", ctype)
			os.Exit(1)
		}
	}
	return rval
}
