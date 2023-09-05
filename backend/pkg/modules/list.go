package modules

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/yonyoucloud/datatable/pkg/model"
)

func (ms *Modules) DescTable(c *gin.Context) {
	name := c.Param("name")

	if name == "" {
		ms.response(1000, "表名不能为空", nil, c)
		return
	}

	if !ms.validateTable(name) {
		ms.response(1002, "该表不允许查询", nil, c)
		return
	}

	data, err := ms.model.DescTable(name)
	if err != nil {
		ms.response(100, "数据库查询报错", err.Error(), c)
		return
	}

	ms.response(0, "", data, c)
}

func (ms *Modules) ShowTables(c *gin.Context) {
	data, err := ms.model.ShowTable()
	if err != nil {
		ms.response(100, "数据库查询报错", err.Error(), c)
		return
	}

	var tables []model.Table
	for _, table := range data {
		if ms.validateTable(string(table)) {
			tables = append(tables, table)
		}
	}

	ms.response(0, "", tables, c)
}

func (ms *Modules) ContentsTable(c *gin.Context) {
	name := c.Param("name")

	if name == "" {
		ms.response(1000, "表名不能为空", nil, c)
		return
	}

	if !ms.validateTable(name) {
		ms.response(1002, "该表不允许查询", nil, c)
		return
	}

	type Payload struct {
		PageSize    int    `json:"page_size"`
		CurrentPage int    `json:"current_page"`
		SortField   string `json:"sort_field"`
		SortOrder   string `json:"sort_order"`
		Filter      string `json:"filter"`
	}
	var payload Payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		ms.response(1001, "参数解析错误", err.Error(), c)
		return
	}

	data, err := ms.model.ContentsTable(name, payload.PageSize, payload.CurrentPage, payload.SortField, payload.SortOrder, payload.Filter)
	if err != nil {
		ms.response(100, "数据库查询报错", err.Error(), c)
		return
	}

	ms.response(0, "", data, c)
}

func (ms *Modules) validateTable(name string) bool {
	if name == "" {
		return false
	}

	for _, allowTable := range ms.config.Mysql.AllowTable {
		str := fmt.Sprintf(`(?i)%s`, allowTable)
		if regexp.MustCompile(str).MatchString(name) {
			return true
		}
	}

	return false
}
