package weavedb

import (
	"database/sql"
)

type MssqlCrawler struct {
	connectionString string
	db               *sql.DB
	tables           []string
}

func NewMssqlCrawler(connectionString string) *MssqlCrawler {
	return &MssqlCrawler{
		connectionString: connectionString,
	}
}

func (c *MssqlCrawler) Connect() error {
	db, err := sql.Open(c.DriverName(), c.connectionString)
	if err != nil {
		return err
	}
	c.db = db
	return nil
}

func (c *MssqlCrawler) DriverName() string {
	return "mssql"
}

func (c *MssqlCrawler) LoadTableNames() error {
	sql := ""
}

func (c *MssqlCrawler) ProcessTables() error {
	return nil
}
