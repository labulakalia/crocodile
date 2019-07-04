package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labulaka521/logging"
	"sync"
)

var (
	sqldb *sql.DB
	once  sync.Once
)

func New(mysqldsn string, maxilde, maxopen int) *sql.DB {
	var (
		err error
	)
	once.Do(func() {
		if sqldb, err = createMysqlClient(mysqldsn, maxilde, maxopen); err != nil {
			logging.Fatalf("[New CLient] createMysqlClient Err: %v", err)
			return
		}
	})
	return sqldb
}

func createMysqlClient(mysqldsn string, maxilde, maxopen int) (db *sql.DB, err error) {

	if db, err = sql.Open("mysql", mysqldsn); err != nil {
		return
	}
	db.SetMaxOpenConns(maxopen)
	db.SetMaxIdleConns(maxilde)
	if err = db.Ping(); err != nil {
		return
	}

	return
}
