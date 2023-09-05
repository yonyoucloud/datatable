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
	var results []map[string]interface{} // 存储查询结果的切片
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

func getFilterSql(filter string) string {
	pre := "1 = 1"

	var data interface{}
	if err := json.Unmarshal([]byte(filter), &data); err != nil {
		return pre
	}

	ops := []string{pre}
	for column, v := range data.(map[string]interface{}) {
		op, value := "", ""
		for kk, vv := range v.(map[string]interface{}) {
			if kk == "name" {
				value = vv.(string)
			} else if kk == "type" {
				op = vv.(string)
			}
		}
		if op == "IN" {
			ops = append(ops, fmt.Sprintf("%s %s ('%s')", column, op, value))
		} else {
			ops = append(ops, fmt.Sprintf("%s %s '%s'", column, op, value))
		}
	}

	return strings.Join(ops, " AND ")
}
