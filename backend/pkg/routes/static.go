package routes

import (
	"github.com/gin-gonic/gin"
)

func (rs *Routes) AddStatic(rg *gin.RouterGroup, staticDir string) {
	rg.Static("/", staticDir)
}

func (rs *Routes) CreateConfigJs(rg *gin.RouterGroup) {
	rg.GET("/create/config/js", rs.modules.CreateConfigJs)
}
