package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

type (
	Table string
	Field struct {
		Field   string
		Type    string
		Null    string
		Key     string
		Default string
		Extra   string
	}
)

func (m *Model) ShowTable() ([]Table, error) {
	var tables []Table

	if result := m.db.Raw("SHOW TABLES").Scan(&tables); result.Error != nil {
		return tables, result.Error
	}

	return tables, nil
}

func (m *Model) DescTable(name string) ([]Field, error) {
	var fields []Field

	if result := m.db.Raw(fmt.Sprintf("DESCRIBE %s", name)).Scan(&fields); result.Error != nil {
		return fields, result.Error
	}

	return fields, nil
}

func (m *Model) ContentsTable(name string, pageSize, currentPage int, sortField, sortOrder, filter string) (interface{}, error) {
	var results []map[string]interface{}
	var total int64

	// 处理过滤
	filterSql := getFilterSql(filter)

	// 获取总数
	sql := fmt.Sprintf("SELECT COUNT(1) FROM `%s` WHERE %s LIMIT 1", name, filterSql)
	if result := m.db.Raw(sql).Scan(&total); result.Error != nil {
		return results, result.Error
	}

	// 处理查询
	sql = fmt.Sprintf("SELECT * FROM `%s` WHERE %s", name, filterSql)

	// 处理排序
	if sortOrder != "" && sortField != "" {
		sql = fmt.Sprintf("%s ORDER BY %s %s", sql, sortField, strings.ToUpper(sortOrder))
	}

	// 处理分页
	offset := (currentPage - 1) * pageSize
	sql = fmt.Sprintf("%s LIMIT %d, %d", sql, offset, pageSize)

	if result := m.db.Raw(sql).Scan(&results); result.Error != nil {
		return results, result.Error
	}

	type Result struct {
		List  []map[string]interface{} `json:"list"`
		Total int64                    `json:"total"`
	}

	r := Result{
		List:  results,
		Total: total,
	}

	return r, nil
}

func (m *Model) StatTable(name, property string, filter string) (interface{}, error) {
	var results []map[string]interface{}

	// 处理过滤
	filterSql := getFilterSql(filter)

	// 获取最大最小
	type Result struct {
		Min string
		Max string
	}
	var res Result
	sql := fmt.Sprintf("SELECT MIN(%s) AS min, MAX(%s) AS max FROM %s%s%s WHERE %s", property, property, "`", name, "`", filterSql)
	if result := m.db.Raw(sql).Scan(&res); result.Error != nil {
		return results, result.Error
	}

	// 处理查询
	sql = fmt.Sprintf(`SELECT
    date_format(t.interval_start, '%%Y-%%m-%%d %%H:%%i:%%s') AS start_time,
    date_format(t.interval_end, '%%Y-%%m-%%d %%H:%%i:%%s') AS end_time,
    COUNT(1) AS record_count
FROM (
    SELECT
        ('%s') +
        INTERVAL (seq - 1) * (TIMESTAMPDIFF(SECOND, ('%s'), ('%s')) / 30) SECOND AS interval_start,
        ('%s') +
        INTERVAL seq * (TIMESTAMPDIFF(SECOND, ('%s'), ('%s')) / 30) SECOND AS interval_end
    FROM (
        SELECT 1 AS seq
        UNION ALL SELECT 2
        UNION ALL SELECT 3
        UNION ALL SELECT 4
        UNION ALL SELECT 5
        UNION ALL SELECT 6
        UNION ALL SELECT 7
        UNION ALL SELECT 8
        UNION ALL SELECT 9
        UNION ALL SELECT 10
        UNION ALL SELECT 11
        UNION ALL SELECT 12
        UNION ALL SELECT 13
        UNION ALL SELECT 14
        UNION ALL SELECT 15
        UNION ALL SELECT 16
        UNION ALL SELECT 17
        UNION ALL SELECT 18
        UNION ALL SELECT 19
        UNION ALL SELECT 20
        UNION ALL SELECT 21
        UNION ALL SELECT 22
        UNION ALL SELECT 23
        UNION ALL SELECT 24
        UNION ALL SELECT 25
        UNION ALL SELECT 26
        UNION ALL SELECT 27
        UNION ALL SELECT 28
        UNION ALL SELECT 29
        UNION ALL SELECT 30
    ) AS sequence
) AS t
LEFT JOIN %s%s%s AS l ON l.%s >= t.interval_start AND l.%s < t.interval_end
WHERE %s
GROUP BY t.interval_start, t.interval_end`, res.Min, res.Min, res.Max, res.Min, res.Min, res.Max, "`", name, "`", property, property, filterSql)

	if result := m.db.Raw(sql).Scan(&results); result.Error != nil {
		return results, result.Error
	}

	return results, nil
}

func getFilterSql(filter string) string {
	pre := "1 = 1"

	var data interface{}
	if err := json.Unmarshal([]byte(filter), &data); err != nil {
		return pre
	}

	ops := []string{pre}
	for column, v := range data.(map[string]interface{}) {
		op, value1, value2 := "", "", ""
		for kk, vv := range v.(map[string]interface{}) {
			if kk == "value1" {
				value1 = vv.(string)
			} else if kk == "value2" {
				value2 = vv.(string)
			} else if kk == "type" {
				op = vv.(string)
			}
		}
		if op == "IN" {
			ops = append(ops, fmt.Sprintf("%s %s ('%s')", column, op, value1))
		} else if op == "RANGE" {
			ops = append(ops, fmt.Sprintf("%s BETWEEN '%s' AND '%s'", column, value1, value2))
		} else {
			ops = append(ops, fmt.Sprintf("%s %s '%s'", column, op, value1))
		}
	}

	return strings.Join(ops, " AND ")
}
