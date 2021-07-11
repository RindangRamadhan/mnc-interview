package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	"gitlab.com/pt-mai/maihelper"
)

type dbPostgres struct {
	Connection *sqlx.DB
	Error      error
}

type IPostgre interface {
	Init()
	Conn() *sqlx.DB
	Err() error
}

func NewPostgre() IPostgre {
	return &dbPostgres{}
}

//Init Initialization function for PostgreSQL
func (_d *dbPostgres) Init() {
	_d.Connection, _d.Error = maihelper.PostgreSql.ConnectSqlx()

	if _d.Error != nil {
		log.Println(_d.Error)
	}
}

func (_d *dbPostgres) Conn() *sqlx.DB {
	return _d.Connection
}

func (_d *dbPostgres) Err() error {
	return _d.Error
}

var Postgre = NewPostgre()
