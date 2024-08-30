package mysql

import "database/sql"

type MysqlClient struct {
	*sql.DB
	driverName string
	dbName     string
}
