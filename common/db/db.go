package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql" // registry sqlite3 deive
	_ "github.com/mattn/go-sqlite3"    // registry  mysql drive
)

var (
	_db *sql.DB
)

type dbCfg struct {
	DriveName         string
	Dsn               string
	MaxIdleConnection int
	MaxOpenConnection int
	MaxQueryTime      time.Duration
}

// GetConn from db conn pool
func GetConn(ctx context.Context) (*sql.Conn, error) {
	return _db.Conn(ctx)
}

// Option is function option
type Option func(*dbCfg)

// Drivename  Set mysql or sqlite
func Drivename(drivename string) Option {
	return func(dbcfg *dbCfg) {
		dbcfg.DriveName = drivename
	}
}

// Dsn set dbCfg conn addr
func Dsn(dsn string) func(*dbCfg) {
	return func(dbcfg *dbCfg) {
		dbcfg.Dsn = dsn
	}
}

// MaxIdleConnection set sql db max idle conn
func MaxIdleConnection(idle int) Option {
	return func(dbcfg *dbCfg) {
		dbcfg.MaxIdleConnection = idle
	}
}

// MaxOpenConnection  set sql db  max open conn
func MaxOpenConnection(open int) Option {
	return func(dbcfg *dbCfg) {
		dbcfg.MaxOpenConnection = open
	}
}

// MaxQueryTime set sql conn exec max timeout
func MaxQueryTime(query time.Duration) Option {
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

// NewDb create new db
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
