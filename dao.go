package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

const (
	driverName     = "mysql"
	dataSourceName = "root:root@tcp(127.0.0.1:3306)/xorm?charset=utf8"
)

func Conn() {
	engine, err := xorm.NewEngine(driverName, dataSourceName)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(engine.Tables)
}
