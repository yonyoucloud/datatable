package routes

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func (rs *Routes) AddStatic(rg *gin.RouterGroup, host, staticDir string) {
	rg.Static("/", staticDir)

	_ = createConfigJs(host, staticDir)
}

func createConfigJs(host, staticDir string) error {
	file, err := os.OpenFile(fmt.Sprintf("%s/config.js", staticDir), os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	content := fmt.Sprintf(`window.CONFIG = {
    apiHost: 'http://%s/',
}`, host)
	write.WriteString(content)

	write.Flush()

	return nil
}
