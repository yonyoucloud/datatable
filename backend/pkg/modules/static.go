package modules

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func (ms *Modules) CreateConfigJs(c *gin.Context) {
	host := strings.TrimSpace(c.Query("host"))

	err := ms.createConfigJs(host)
	if err != nil {
		ms.response(1000, "创建config.js失败", err.Error(), c)
		return
	}

	ms.response(0, "", "OK", c)
}

func (ms *Modules) createConfigJs(host string) error {
	file, err := os.OpenFile(fmt.Sprintf("%s/config.js", ms.config.StaticDir), os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	content := fmt.Sprintf(`window.CONFIG = {
    apiHost: '%s',
}`, host)
	write.WriteString(content)

	write.Flush()

	return nil
}
