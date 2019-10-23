package db

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var (
	_db *sql.DB
)

type dbCfg struct {
	DriveName         string
	Dsn               string
	MaxIdleConnection int
	MaxOpenConnection int
	MaxQueryTime      int
}

func GetConn(ctx context.Context) (*sql.Conn, error) {
	return _db.Conn(ctx)
}

type Option func(*dbCfg)

func Drivename(drivename string) Option {
	return func(dbcfg *dbCfg) {
		dbcfg.DriveName = drivename
	}
}
func Dsn(dsn string) func(*dbCfg) {
	return func(dbcfg *dbCfg) {
		dbcfg.Dsn = dsn
	}
}
func MaxIdleConnection(idle int) Option {
	return func(dbcfg *dbCfg) {
		dbcfg.MaxIdleConnection = idle
	}
}
func MaxOpenConnection(open int) Option {
	return func(dbcfg *dbCfg) {
		dbcfg.MaxOpenConnection = open
	}
}

func MaxQueryTime(query int) Option {
	return func(dbcfg *dbCfg) {
		dbcfg.MaxQueryTime = query
	}
}

func defaultdbOption() *dbCfg {
	return &dbCfg{
		DriveName:         "sqlite3",
		Dsn:               "sqlite3.db",
		MaxIdleConnection: 10,
		MaxOpenConnection: 5,
		MaxQueryTime:      3,
	}
}

func NewDb(opts ...Option) (err error) {

	dbcfg := defaultdbOption()
	for _, opt := range opts {
		opt(dbcfg)
	}
	_db, err = sql.Open(dbcfg.DriveName, dbcfg.Dsn)
	if err != nil {
		return err
	}
	_db.SetMaxOpenConns(dbcfg.MaxOpenConnection)
	_db.SetMaxIdleConns(dbcfg.MaxIdleConnection)
	err = _db.Ping()
	return err
}
