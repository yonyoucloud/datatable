package main

import (
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/yonyoucloud/datatable/pkg/config"
	"github.com/yonyoucloud/datatable/pkg/modules"
	"github.com/yonyoucloud/datatable/pkg/routes"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var log = logrus.New()
	// log.SetFormatter(&logrus.JSONFormatter{})
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)

	execPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("执行目录获取失败: %s", err.Error())
	}

	configFile := execPath + "/etc/config.yaml"
	cfg, err := config.New(configFile)
	if err != nil {
		log.Fatalf("配置文件加载失败: %s", err.Error())
	}

	// 设置日志级别
	switch cfg.ErrorLevel {
	case "Panic":
		log.SetLevel(logrus.PanicLevel)
		break
	case "Fatal":
		log.SetLevel(logrus.FatalLevel)
		break
	case "Error":
		log.SetLevel(logrus.ErrorLevel)
		break
	case "Warn":
		log.SetLevel(logrus.WarnLevel)
		break
	case "Info":
		log.SetLevel(logrus.InfoLevel)
		break
	case "Debug":
		log.SetLevel(logrus.DebugLevel)
		break
	case "Trace":
		log.SetLevel(logrus.TraceLevel)
		break
	default:
		log.SetLevel(logrus.WarnLevel)
	}

	ms, err := modules.New(cfg, log)
	if err != nil {
		log.Fatalf("初始化模块失败: %s", err.Error())
	}

	var router = gin.Default()
	//router.Use(ms.Logger(log), gin.Recovery())
	router.Use(ms.Cors(cfg.WebRoot))

	rs := routes.New(ms)

	webRoot := router.Group(cfg.WebRoot)
	rs.AddStatic(webRoot, cfg.Host, cfg.StaticDir)

	apiV1 := router.Group("/api/v1")
	rs.AddList(apiV1)

	router.Run(cfg.Host)
}
