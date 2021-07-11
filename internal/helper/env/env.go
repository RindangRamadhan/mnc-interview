package env

import (
	"os"
)

//Env Helper regarding Environment variables
type Env struct {
}

//IsDev get variable based on GO_ENV
func (*Env) IsDev() bool {
	return os.Getenv("GO_ENV") == "development"
}

//IsDev get variable based on IS_DEBUG
func (*Env) IsDebug() bool {
	return os.Getenv("IS_DEBUG") == "true"
}

//IsTest get variable based on GO_ENV
func (*Env) IsTest() bool {
	return os.Getenv("GO_ENV") == "test"
}

//Get Singleton of Env
var Get = &Env{}
