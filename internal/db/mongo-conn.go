package db

import (
	"os"

	"gitlab.com/pt-mai/maihelper"
	"go.mongodb.org/mongo-driver/mongo"
)

type dbMongo struct {
	Connection *mongo.Database
	Database   *mongo.Database
	Error      error
}

type IMongo interface {
	Init()
	Conn() *mongo.Database
	DB() *mongo.Database
	Err() error
}

func NewMongo() IMongo {
	return &dbMongo{}
}

func (_d *dbMongo) Init() {
	os.Getenv("")
	_d.Connection, _d.Error = maihelper.Mongo.Connect()
	if _d.Error == nil {
		_d.Database = _d.Connection.Client().Database(os.Getenv("MONGO_DBNAME"))
	}
}

func (_d *dbMongo) Conn() *mongo.Database {
	return _d.Connection
}

func (_d *dbMongo) Err() error {
	return _d.Error
}

func (_d *dbMongo) DB() *mongo.Database {
	return _d.Database
}

var Mongo = NewMongo()
