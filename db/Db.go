package db

import (
	"fmt"

	"xorm.io/xorm"
)

type Db struct {
	User       string
	Pass       string
	Host       string
	Port       int
	DbName     string
	DbProvider string
	Conn       *xorm.Engine
}

func (db *Db) Connect(user string, pass string, host string, port int, dbName string, dbProvider string) error {
	var err error
	db.User = user
	db.Pass = pass
	db.Host = host
	db.Port = port
	db.DbName = dbName
	db.DbProvider = dbProvider
	connectionString := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", host, port, user, pass, dbName)
	db.Conn, err = xorm.NewEngine(db.DbProvider, connectionString)
	return err
}
