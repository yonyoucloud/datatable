package routes

import "github.com/gin-gonic/gin"

func (rs *Routes) AddList(rg *gin.RouterGroup) {
	listGroup := rg.Group("/list")

	listGroup.GET("/tables", rs.modules.ShowTables)
	listGroup.GET("/table/fields/:name", rs.modules.DescTable)
	listGroup.POST("/table/contents/:name", rs.modules.ContentsTable)
}
