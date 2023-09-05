package routes

import "github.com/gin-gonic/gin"

func (rs *Routes) AddV1(rg *gin.RouterGroup) {
	v1Group := rg.Group("/v1")

	v1Group.GET("/create/config/js", rs.modules.CreateConfigJs)
	v1Group.GET("/get/tables", rs.modules.ShowTables)
	v1Group.GET("/get/table/fields/:name", rs.modules.DescTable)
	v1Group.POST("/get/table/contents/:name", rs.modules.ContentsTable)
}
