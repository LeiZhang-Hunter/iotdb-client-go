package iotdb_go

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var stdConn *sql.DB

func init() {
	var err error
	stdConn, err = sql.Open("iotdb", "iotdb://root:root@127.0.0.1:6667")
	if err != nil {
		panic(err)
	}
}

func TestStd_BasicInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, err := stdConn.ExecContext(context.Background(), "CREATE DATABASE root.std_test")
	require.NoError(t, err)
	defer func() { _, _ = stdConn.ExecContext(context.Background(), "DELETE DATABASE root.std_test") }()

	_, err = stdConn.ExecContext(context.Background(),
		"CREATE TIMESERIES root.std_test.d1.status WITH DATATYPE=BOOLEAN, ENCODING=PLAIN")
	require.NoError(t, err)

	_, err = stdConn.ExecContext(context.Background(),
		"CREATE TIMESERIES root.std_test.d1.temp WITH DATATYPE=FLOAT, ENCODING=PLAIN")
	require.NoError(t, err)

	_, err = stdConn.ExecContext(context.Background(),
		"INSERT INTO root.std_test.d1(timestamp, status, temp) VALUES (?, ?, ?)",
		1, true, 25.5)
	require.NoError(t, err)
}

func TestStd_MultiValueInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, err := stdConn.ExecContext(context.Background(), "CREATE DATABASE root.std_multi")
	require.NoError(t, err)
	defer func() { _, _ = stdConn.ExecContext(context.Background(), "DELETE DATABASE root.std_multi") }()

	_, err = stdConn.ExecContext(context.Background(),
		"CREATE TIMESERIES root.std_multi.d1.status WITH DATATYPE=BOOLEAN, ENCODING=PLAIN")
	require.NoError(t, err)
	_, err = stdConn.ExecContext(context.Background(),
		"CREATE TIMESERIES root.std_multi.d1.temp WITH DATATYPE=FLOAT, ENCODING=PLAIN")
	require.NoError(t, err)

	_, err = stdConn.ExecContext(context.Background(),
		"INSERT INTO root.std_multi.d1(timestamp, status, temp) VALUES (?, ?, ?), (?, ?, ?)",
		1, true, 25.0,
		2, false, 26.0)
	require.NoError(t, err)
}

func TestStd_InsertTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, err := stdConn.ExecContext(context.Background(), "CREATE DATABASE root.std_types")
	require.NoError(t, err)
	defer func() { _, _ = stdConn.ExecContext(context.Background(), "DELETE DATABASE root.std_types") }()

	_, err = stdConn.ExecContext(context.Background(),
		"CREATE TIMESERIES root.std_types.d1.int_val WITH DATATYPE=INT64, ENCODING=PLAIN")
	require.NoError(t, err)
	_, err = stdConn.ExecContext(context.Background(),
		"CREATE TIMESERIES root.std_types.d1.float_val WITH DATATYPE=FLOAT, ENCODING=PLAIN")
	require.NoError(t, err)

	_, err = stdConn.ExecContext(context.Background(),
		"INSERT INTO root.std_types.d1(timestamp, int_val, float_val) VALUES (?, ?, ?)",
		1, 100, 99.9)
	require.NoError(t, err)
}

func TestStd_QueryBasic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, err := stdConn.ExecContext(context.Background(), "CREATE DATABASE root.std_query")
	require.NoError(t, err)
	defer func() { _, _ = stdConn.ExecContext(context.Background(), "DELETE DATABASE root.std_query") }()

	_, err = stdConn.ExecContext(context.Background(),
		"CREATE TIMESERIES root.std_query.d1.temp WITH DATATYPE=FLOAT, ENCODING=PLAIN")
	require.NoError(t, err)

	_, err = stdConn.ExecContext(context.Background(),
		"INSERT INTO root.std_query.d1(timestamp, temp) VALUES (?, ?)",
		1704067200000, 25.5)
	require.NoError(t, err)

	rows, err := stdConn.QueryContext(context.Background(),
		"SELECT temp FROM root.std_query.d1")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.GreaterOrEqual(t, count, 1)
}

func TestStd_ScanStruct(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	type Reading struct {
		Time int64
		Temp float64
		Hum  float64
	}

	_, err := stdConn.ExecContext(context.Background(), "CREATE DATABASE root.std_scan")
	require.NoError(t, err)
	defer func() { _, _ = stdConn.ExecContext(context.Background(), "DELETE DATABASE root.std_scan") }()

	_, err = stdConn.ExecContext(context.Background(),
		"CREATE TIMESERIES root.std_scan.d1.temp WITH DATATYPE=FLOAT, ENCODING=PLAIN")
	require.NoError(t, err)
	_, err = stdConn.ExecContext(context.Background(),
		"CREATE TIMESERIES root.std_scan.d1.hum WITH DATATYPE=DOUBLE, ENCODING=PLAIN")
	require.NoError(t, err)

	_, err = stdConn.ExecContext(context.Background(),
		"INSERT INTO root.std_scan.d1(timestamp, temp, hum) VALUES (?, ?, ?)",
		1704067200000, 25.5, 60.5)
	require.NoError(t, err)

	rows, err := stdConn.QueryContext(context.Background(),
		"SELECT temp, hum FROM root.std_scan.d1")
	require.NoError(t, err)
	defer rows.Close()

	require.True(t, rows.Next())
	var r Reading
	err = rows.Scan(&r.Time, &r.Temp, &r.Hum)
	require.NoError(t, err)
	assert.InDelta(t, 25.5, r.Temp, 0.1)
	assert.InDelta(t, 60.5, r.Hum, 0.1)
}

func TestStd_ShowDatabases(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	rows, err := stdConn.QueryContext(context.Background(), "SHOW DATABASES")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.GreaterOrEqual(t, count, 1)
}
