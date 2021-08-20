package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	XEngine *xorm.Engine
)

func Conn() {
	cfg := GlobalCfg.DataBaseConfig
	driverName := cfg.Driver
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.DBName)
	var err error
	XEngine, err = xorm.NewEngine(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("数据库连接失败: %v \n", err)
	}
}
