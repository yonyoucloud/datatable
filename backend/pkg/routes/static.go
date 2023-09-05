package routes

import (
	"github.com/gin-gonic/gin"
)

func (rs *Routes) AddStatic(rg *gin.RouterGroup, staticDir string) {
	rg.Static("/", staticDir)
}
