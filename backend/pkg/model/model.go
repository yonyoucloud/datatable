package model

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"github.com/yonyoucloud/datatable/pkg/config"
)

type Model struct {
	config *config.Config
	db     *gorm.DB
}

func New(cfg *config.Config) (*Model, error) {
	m := &Model{config: cfg}
	err := m.initDB()

	return m, err
}

func (m *Model) initDB() error {
	var err error
	db, err := gorm.Open(mysql.Open(m.config.Mysql.Master), &gorm.Config{})
	if err != nil {
		return err
	}
	m.db = db

	var sources []gorm.Dialector
	for _, dsn := range m.config.Mysql.Sources {
		sources = append(sources, mysql.Open(dsn))
	}

	var replicas []gorm.Dialector
	for _, dsn := range m.config.Mysql.Replicas {
		replicas = append(replicas, mysql.Open(dsn))
	}

	err = m.db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  sources,
		Replicas: replicas,
		// sources/replicas 负载均衡策略
		Policy: dbresolver.RandomPolicy{},
	}).
		//设置了连接可复用的最大时间
		SetConnMaxIdleTime(time.Hour).
		//设置了连接可复用的最大时间
		SetConnMaxLifetime(m.config.Mysql.ConnMaxLifetime).
		// 设置空闲连接池中连接的最大数量
		SetMaxIdleConns(m.config.Mysql.MaxIdleConns).
		// 设置打开数据库连接的最大数量
		SetMaxOpenConns(m.config.Mysql.MaxOpenConns))

	return err
}
