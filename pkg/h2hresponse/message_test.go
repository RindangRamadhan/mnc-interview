// +build unit

package h2hresponse

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestResponseMessageGenerate(t *testing.T) {
	var err error
	var msgRes string
	
	var tplData ResponseMessageTpl
	tplData = ResponseMessageTpl{
		CodeProduct:     "A1",
		CustomerNumber:  "085647047131",
		SN:              "232323232",
		SisaSaldo:       "10000000000",
		SellingPrice:    "21000",
		CreatedAt:       "",
		SupplierMessage: "sukses",
	}
	msgRes, err = ResponseMessageGenerate(ResponseMsgSuccessPrepaid, tplData)
	assert.Nilf(t, err, "`err` should be nil, got %v instead", err)
	assert.NotEmptyf(t, msgRes, "`msgRes` should be not be empty")
	log.Println(msgRes, "Test....")
}
