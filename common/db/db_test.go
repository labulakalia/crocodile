package db

import (
	"context"
	"os"
	"testing"
)

func TestNewDb(t *testing.T) {

	err := NewDb(Drivename("sqlite3"),
		Dsn("sqlite3.db"),
		MaxIdleConnection(10),
		MaxQueryTime(3),
		MaxQueryTime(3),
		MaxOpenConnection(3),
	)
	if err != nil {
		t.Fatalf("NewDb Err: %v", err)
	}
	conn, err := GetConn(context.Background())
	if err != nil {
		t.Fatalf("Get Conn Err: %v", err)
	}
	conn.Close()
	_ = os.Remove("sqlite3.db")
}
