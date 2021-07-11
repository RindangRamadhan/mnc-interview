package h2hresponse

import (
	"log"
	"net/http"

	"gitlab.com/pt-mai/maihelper/mailog"

	jsoniter "github.com/json-iterator/go"
)

var Json = jsoniter.ConfigCompatibleWithStandardLibrary

type Tpl struct {
	HttpStatusCode int
	Body           interface{}
}

type Body struct {
	MAIResponseCode int         `json:"RC"`
	MAIStatus       string      `json:"status"`
	MAIMessage      string      `json:"message"`
	MAIError        interface{} `json:"error,omitempty"`
}

type BodyGeneral struct {
	Body
	MAIData interface{} `json:"data,omitempty"`
}

func Write(w http.ResponseWriter, res Tpl) {
	if body, ok := res.Body.(BodyTrx); ok {
		log.Println("WriwriteTrxte")
		writeTrx(w, res, body)
	} else if bodyG, okG := res.Body.(BodyGeneral); okG {
		log.Println("writeGeneral 1")
		writeGeneral(w, res, bodyG)
	} else {
		log.Println("writeGeneral 2")
		writeGeneral(w, res, BodyGeneral{})
	}
}

func writeTrx(w http.ResponseWriter, res Tpl, body BodyTrx) {
	var err error
	w.Header().Set("Content-Type", "application/json")
	if res.HttpStatusCode == 0 {
		res.HttpStatusCode = http.StatusOK
	}
	w.WriteHeader(res.HttpStatusCode)

	if body.MAIResponseCode == 0 {
		body.MAIResponseCode = res.HttpStatusCode
	}

	if body.MAIStatus == "" {
		body.MAIStatus = StatusName(body.MAIResponseCode)
	}

	if body.MAIMessage == "" {
		body.MAIMessage = StatusText(body.MAIResponseCode)
	}

	var bodyB []byte
	bodyB, err = Json.Marshal(body)
	if err != nil {
		bodyB = []byte(`{"RC":500, "status":"failed", "message":"Unknown error"}`)
		_, _ = w.Write(bodyB)
		return
	}
	_, err = w.Write(bodyB)

	if err != nil {
		log.Println(mailog.Error(err))
	}
}

func writeGeneral(w http.ResponseWriter, res Tpl, body BodyGeneral) {
	var err error
	w.Header().Set("Content-Type", "application/json")
	if res.HttpStatusCode == 0 {
		res.HttpStatusCode = http.StatusOK
	}
	w.WriteHeader(res.HttpStatusCode)

	if body.MAIResponseCode == 0 {
		body.MAIResponseCode = res.HttpStatusCode
	}

	if body.MAIStatus == "" {
		body.MAIStatus = StatusName(body.MAIResponseCode)
	}

	if body.MAIMessage == "" {
		body.MAIMessage = StatusText(body.MAIResponseCode)
	}

	var bodyB []byte
	bodyB, err = Json.Marshal(body)
	if err != nil {
		bodyB = []byte(`{"RC":500, "status":"failed", "message":"Unknown error"}`)
		_, _ = w.Write(bodyB)
		return
	}
	_, err = w.Write(bodyB)

	if err != nil {
		log.Println(mailog.Error(err))
	}
}
