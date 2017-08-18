package database

import (
	"database/sql"

  // Imports the mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Database struct to connect to database
type Database struct {
	db *sql.DB
	connString string
}

// NewDatabaseConnection creates a new database connection
func NewDatabaseConnection(connString string) *Database {
	dtb := new(Database)
	dtb.connString = connString
	return dtb
}

// Connect function
func (datab *Database) Connect() error {

	if datab.db == nil {
		dbt, err := sql.Open("mysql", datab.connString)
		if err != nil {
			return err
		}

		err = dbt.Ping()

		if err != nil {
			return err
		}

		datab.db = dbt
	}

	return nil
}

//ExecuteNonQuery executes non query
func (datab *Database) ExecuteNonQuery(query string, args ...interface{}) error {
	if e := datab.Connect(); e != nil {
		return e
	}
	defer datab.Close()
	_, err := datab.db.Exec(query, args...)

	if err != nil {
		return err
	}

	return nil
}

// ExecuteQuery executes the query and returs the obtained rows.
func (datab *Database) ExecuteQuery(query string, args ...interface{}) (*sql.Rows,error) {
	if e := datab.Connect(); e != nil {
		return nil,e
	}
	defer datab.Close()
	return datab.db.Query(query, args...)
}

// Close function
func (datab *Database) Close() {
	if datab.db != nil {
		datab.db.Close()
		datab.db = nil
	}
}
