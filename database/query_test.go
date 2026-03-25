package iotdb_go

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var queryConn *sql.DB

func init() {
	var err error
	queryConn, err = sql.Open("iotdb", "iotdb://root:root@127.0.0.1:6667")
	if err != nil {
		panic(err)
	}
}

func setupQueryTestData(t *testing.T, db *sql.DB, sg string) func() {
	_, err := db.ExecContext(context.Background(), "CREATE DATABASE "+sg)
	require.NoError(t, err)

	_, err = db.ExecContext(context.Background(),
		"CREATE TIMESERIES "+sg+".d1.temp WITH DATATYPE=FLOAT, ENCODING=PLAIN")
	require.NoError(t, err)

	_, err = db.ExecContext(context.Background(),
		"CREATE TIMESERIES "+sg+".d1.hum WITH DATATYPE=DOUBLE, ENCODING=PLAIN")
	require.NoError(t, err)

	_, err = db.ExecContext(context.Background(),
		"CREATE TIMESERIES "+sg+".d1.status WITH DATATYPE=BOOLEAN, ENCODING=PLAIN")
	require.NoError(t, err)

	_, err = db.ExecContext(context.Background(),
		"CREATE TIMESERIES "+sg+".d2.temp WITH DATATYPE=FLOAT, ENCODING=PLAIN")
	require.NoError(t, err)

	_, err = db.ExecContext(context.Background(),
		"INSERT INTO "+sg+".d1(timestamp, temp, hum, status) VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?)",
		1704067200000, 20.0, 60.0, true,
		1704067260000, 21.0, 61.0, false,
		1704067320000, 22.0, 62.0, true,
		1704067380000, 23.0, 63.0, false,
		1704067440000, 24.0, 64.0, true)
	require.NoError(t, err)

	_, err = db.ExecContext(context.Background(),
		"INSERT INTO "+sg+".d2(timestamp, temp) VALUES (?, ?), (?, ?), (?, ?)",
		1704067200000, 30.0,
		1704067260000, 31.0,
		1704067320000, 32.0)
	require.NoError(t, err)

	return func() {
		_, _ = db.ExecContext(context.Background(), "DELETE DATABASE "+sg)
	}
}

// ==================== SELECT 子句基础查询 ====================

func TestQuery_SelectSingleColumn(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.select_single")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT temp FROM root.select_single.d1")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 5, count)
}

func TestQuery_SelectMultipleColumns(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.select_multi")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT temp, hum, status FROM root.select_multi.d1")
	require.NoError(t, err)
	defer rows.Close()

	cols, err := rows.Columns()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(cols), 3)
}

func TestQuery_SelectAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.select_all")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT * FROM root.select_all.d1")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 5, count)
}

func TestQuery_SelectWithAlias(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.select_alias")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT temp AS temperature, hum AS humidity FROM root.select_alias.d1")
	require.NoError(t, err)
	defer rows.Close()

	cols, err := rows.Columns()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(cols), 2)
}

// ==================== WHERE 子句查询 ====================

func TestQuery_WhereTimeFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.where_time")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT temp FROM root.where_time.d1 WHERE time >= 1704067260000")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 4, count)
}

func TestQuery_WhereTimeRange(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.where_range")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT temp FROM root.where_range.d1 WHERE time >= 1704067260000 AND time < 1704067380000")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 2, count)
}

func TestQuery_WhereValueFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.where_value")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT temp FROM root.where_value.d1 WHERE temp > 22.0")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 2, count)
}

func TestQuery_WhereBooleanFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.where_bool")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT temp FROM root.where_bool.d1 WHERE status = true")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 3, count)
}

func TestQuery_WhereCombinedFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.where_combined")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT temp FROM root.where_combined.d1 WHERE time > 1704067260000 AND temp > 22.0 AND status = true")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 1, count)
}

// ==================== 聚合函数查询 ====================

func TestQuery_AggCount(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.agg_count")
	defer cleanup()

	row := queryConn.QueryRowContext(context.Background(),
		"SELECT COUNT(temp) FROM root.agg_count.d1")
	var count int64
	err := row.Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

func TestQuery_AggSum(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.agg_sum")
	defer cleanup()

	row := queryConn.QueryRowContext(context.Background(),
		"SELECT SUM(temp) FROM root.agg_sum.d1")
	var sum float64
	err := row.Scan(&sum)
	require.NoError(t, err)
	assert.InDelta(t, 110.0, sum, 0.1)
}

func TestQuery_AggAvg(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.agg_avg")
	defer cleanup()

	row := queryConn.QueryRowContext(context.Background(),
		"SELECT AVG(temp) FROM root.agg_avg.d1")
	var avg float64
	err := row.Scan(&avg)
	require.NoError(t, err)
	assert.InDelta(t, 22.0, avg, 0.1)
}

func TestQuery_AggMinMax(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.agg_minmax")
	defer cleanup()

	row := queryConn.QueryRowContext(context.Background(),
		"SELECT MIN_VALUE(temp), MAX_VALUE(temp) FROM root.agg_minmax.d1")
	var minTemp, maxTemp float64
	err := row.Scan(&minTemp, &maxTemp)
	require.NoError(t, err)
	assert.InDelta(t, 20.0, minTemp, 0.1)
	assert.InDelta(t, 24.0, maxTemp, 0.1)
}

func TestQuery_AggFirstLastValue(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.agg_firstlast")
	defer cleanup()

	row := queryConn.QueryRowContext(context.Background(),
		"SELECT FIRST_VALUE(temp), LAST_VALUE(temp) FROM root.agg_firstlast.d1")
	var first, last float64
	err := row.Scan(&first, &last)
	require.NoError(t, err)
	assert.InDelta(t, 20.0, first, 0.1)
	assert.InDelta(t, 24.0, last, 0.1)
}

func TestQuery_AggMinMaxTime(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.agg_mintime")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT MIN_TIME(temp), MAX_TIME(temp) FROM root.agg_mintime.d1")
	require.NoError(t, err)
	defer rows.Close()

	assert.True(t, rows.Next())
}

func TestQuery_MultipleAggregations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.agg_multi")
	defer cleanup()

	row := queryConn.QueryRowContext(context.Background(),
		"SELECT COUNT(temp), SUM(temp), AVG(temp), MIN_VALUE(temp), MAX_VALUE(temp) FROM root.agg_multi.d1")
	var count int64
	var sum, avg, min, max float64
	err := row.Scan(&count, &sum, &avg, &min, &max)
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.InDelta(t, 110.0, sum, 0.1)
	assert.InDelta(t, 22.0, avg, 0.1)
	assert.InDelta(t, 20.0, min, 0.1)
	assert.InDelta(t, 24.0, max, 0.1)
}

// ==================== GROUP BY 时间区间聚合 ====================

func TestQuery_GroupByTimeInterval(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.group_time")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT COUNT(temp), AVG(temp) FROM root.group_time.d1 GROUP BY ([1704067200000, 1704067500000), 2m)")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Greater(t, count, 0)
}

func TestQuery_GroupByWithWhere(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.group_where")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT COUNT(temp) FROM root.group_where.d1 WHERE time < 1704067440000 GROUP BY ([1704067200000, 1704067500000), 2m)")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Greater(t, count, 0)
}

// ==================== ORDER BY 排序 ====================

func TestQuery_OrderByTimeDesc(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.order_time")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT temp FROM root.order_time.d1 ORDER BY time DESC")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 5, count)
}

func TestQuery_OrderByTimeAsc(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.order_asc")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT temp FROM root.order_asc.d1 ORDER BY time ASC")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 5, count)
}

// ==================== LIMIT 和 OFFSET ====================

func TestQuery_Limit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.limit_test")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT * FROM root.limit_test.d1 LIMIT 3")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 3, count)
}

func TestQuery_LimitOffset(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.limit_offset")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT temp FROM root.limit_offset.d1 ORDER BY time ASC LIMIT 2 OFFSET 2")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 2, count)
}

// ==================== ALIGN BY DEVICE ====================

func TestQuery_AlignByDevice(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.align_device")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT temp FROM root.align_device.d1, root.align_device.d2 ALIGN BY DEVICE")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.GreaterOrEqual(t, count, 1)
}

func TestQuery_AlignByTime(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.align_time")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT * FROM root.align_time.d1 ALIGN BY TIME")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 5, count)
}

// ==================== LAST 最新点查询 ====================

func TestQuery_LastSingleColumn(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.last_single")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT LAST temp FROM root.last_single.d1")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 1, count)
}

func TestQuery_LastMultipleColumns(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.last_multi")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT LAST temp, hum FROM root.last_multi.d1")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 2, count)
}

func TestQuery_LastAllColumns(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.last_all")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT LAST * FROM root.last_all.d1")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.GreaterOrEqual(t, count, 3)
}

func TestQuery_LastWithWhere(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.last_where")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT LAST temp FROM root.last_where.d1 WHERE time >= 1704067320000")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 1, count)
}

// ==================== 表达式查询 ====================

func TestQuery_ArithmeticExpression(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.expr_arith")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT temp + 10, temp * 2, temp / 2 FROM root.expr_arith.d1")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 5, count)
}

func TestQuery_AggregationExpression(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.expr_agg")
	defer cleanup()

	row := queryConn.QueryRowContext(context.Background(),
		"SELECT AVG(temp) + 1, SUM(temp) * 2 FROM root.expr_agg.d1")
	var avgPlus, sumTimes float64
	err := row.Scan(&avgPlus, &sumTimes)
	require.NoError(t, err)
	assert.InDelta(t, 23.0, avgPlus, 0.1)
	assert.InDelta(t, 220.0, sumTimes, 0.1)
}

// ==================== 数据类型扫描 ====================

func TestQuery_ScanTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.scan_types")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT temp, hum, status FROM root.scan_types.d1 WHERE time = 1704067200000")
	require.NoError(t, err)
	defer rows.Close()

	require.True(t, rows.Next())
	var ts int64
	var temp float64
	var hum float64
	var status bool
	err = rows.Scan(&ts, &temp, &hum, &status)
	require.NoError(t, err)
	assert.InDelta(t, 20.0, temp, 0.1)
	assert.InDelta(t, 60.0, hum, 0.1)
	assert.True(t, status)
}

func TestQuery_ScanNull(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.scan_null")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT temp FROM root.scan_null.d1 WHERE time = 999999999999999")
	require.NoError(t, err)
	defer rows.Close()

	assert.False(t, rows.Next())
}

// ==================== 元数据查询 ====================

func TestQuery_ShowDatabases(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	rows, err := queryConn.QueryContext(context.Background(), "SHOW DATABASES")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.GreaterOrEqual(t, count, 1)
}

func TestQuery_ShowTimeseries(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.show_ts")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SHOW TIMESERIES root.show_ts.**")
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	assert.Equal(t, 4, count)
}

func TestQuery_ColumnsMetadata(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	cleanup := setupQueryTestData(t, queryConn, "root.col_meta")
	defer cleanup()

	rows, err := queryConn.QueryContext(context.Background(),
		"SELECT * FROM root.col_meta.d1")
	require.NoError(t, err)
	defer rows.Close()

	cols, err := rows.Columns()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(cols), 3)
}
