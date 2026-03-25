package iotdb_go

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnect_Ping(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, err := sql.Open("iotdb", "iotdb://root:root@127.0.0.1:6667")
	require.NoError(t, err)
	defer db.Close()

	err = db.PingContext(context.Background())
	require.NoError(t, err)
}

func TestConnect_PingTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, err := sql.Open("iotdb", "iotdb://root:root@127.0.0.1:6667")
	require.NoError(t, err)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	require.NoError(t, err)
}

func TestConnect_ConnectionPool(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, err := sql.Open("iotdb", "iotdb://root:root@127.0.0.1:6667")
	require.NoError(t, err)
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(10 * time.Minute)

	assert.Equal(t, 10, db.Stats().MaxOpenConnections)

	conn, err := db.Conn(context.Background())
	require.NoError(t, err)
	defer conn.Close()

	assert.GreaterOrEqual(t, db.Stats().OpenConnections, 1)
}

func TestConnect_Close(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, err := sql.Open("iotdb", "iotdb://root:root@127.0.0.1:6667")
	require.NoError(t, err)

	err = db.Close()
	require.NoError(t, err)
}

func TestConnect_QueryBasic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, err := sql.Open("iotdb", "iotdb://root:root@127.0.0.1:6667")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.ExecContext(context.Background(), "CREATE DATABASE root.conn_test")
	require.NoError(t, err)
	defer func() {
		_, _ = db.ExecContext(context.Background(), "DELETE DATABASE root.conn_test")
	}()

	_, err = db.ExecContext(context.Background(),
		"CREATE TIMESERIES root.conn_test.d1.value WITH DATATYPE=FLOAT, ENCODING=PLAIN")
	require.NoError(t, err)

	_, err = db.ExecContext(context.Background(),
		"INSERT INTO root.conn_test.d1(timestamp, value) VALUES (?, ?)", 1, 123.45)
	require.NoError(t, err)

	rows, err := db.QueryContext(context.Background(), "SELECT value FROM root.conn_test.d1")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.GreaterOrEqual(t, count, 1)
}

func TestConnect_ShowDatabases(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, err := sql.Open("iotdb", "iotdb://root:root@127.0.0.1:6667")
	require.NoError(t, err)
	defer db.Close()

	rows, err := db.QueryContext(context.Background(), "SHOW DATABASES")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.GreaterOrEqual(t, count, 1)
}
