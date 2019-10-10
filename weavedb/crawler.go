package weavedb

type Crawler interface {
	Connect() error
	DriverName() string
	LoadTableNames() error
	ProcessTables() error
}
