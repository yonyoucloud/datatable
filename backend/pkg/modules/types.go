package modules

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/yonyoucloud/datatable/pkg/config"
	"github.com/yonyoucloud/datatable/pkg/model"
)

type (
	Modules struct {
		config *config.Config
		log    *logrus.Logger
		model  *model.Model
	}

	Response struct {
		ErrorCode    int         `json:"error_code"`
		ErrorMessage string      `json:"error_message"`
		Data         interface{} `json:"data"`
	}
)

func New(cfg *config.Config, log *logrus.Logger) (*Modules, error) {
	ms := &Modules{
		config: cfg,
		log:    log,
	}

	m, err := model.New(cfg)
	if err != nil {
		return ms, err
	}
	ms.model = m

	return ms, nil
}

func (ms *Modules) response(errorCode int, errorMessage string, data interface{}, c *gin.Context) {
	resp := Response{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
		Data:         data,
	}
	c.JSON(http.StatusOK, resp)
}
