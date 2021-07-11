package h2hresponse

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"text/template"

	"gitlab.com/pt-mai/maihelper"
	"gitlab.com/pt-mai/maihelper/maigrpc"
	"gitlab.com/pt-mai/maihelper/mailog"
	jsonparse "gitlab.com/pt-mai/maihelper/utils/json-parse"
)

type ResponsePostpaidGenerate struct {
	Action   string
	Category string
	DataTpl  interface{}
}

type ResponseMessageTpl struct {
	CodeProduct     string
	CustomerNumber  string
	SN              string
	SisaSaldo       string
	SellingPrice    string
	CreatedAt       string
	SupplierMessage string
	PartnerTrxID    string
}

type ResponseMessagePostpaidTpl struct {
	CodeProduct     string
	SN              string
	SellingPrice    string
	CreatedAt       string
	PartnerTrxID    string
	Qty             string `json:"qty"`
	Tag             string `json:"tag"`
	Admin           string `json:"admin"`
	Total           string `json:"total"`
	Status          string `json:"status"`
	Periode         string `json:"periode,omitempty"`
	CustomerName    string `json:"customer_name"`
	ProductScheme   string `json:"product_scheme"`
	CustomerNumber  string `json:"customer_number"`
	SupplierMessage string `json:"supplier_message,omitempty"`
	CategoryId      string `json:"category_id"`
	TarifDaya       string `json:"tarif_daya,omitempty"`
	Kwh             string `json:"kwh,omitempty"`
	Msn             string `json:"msn,omitempty"`
	JumlahPeserta   string `json:"jumlah_peserta,omitempty"`
	JumlahBulan     string `json:"jumlah_bulan,omitempty"`
	Other           string `json:"other,omitempty"`
	SisaSaldo       string
	Raw             interface{} `json:"raw"`
}

type CategoryTpl struct {
	Label   string
	Message string
	Data    interface{}
}

// type ResponseMessagePLNTpl struct {
// 	Qty             string
// 	Tag             string
// 	Admin           string
// 	Total           string
// 	Status          string
// 	Periode         string
// 	CustomerName    string
// 	ProductScheme   string
// 	CustomerNumber  string
// 	SupplierMessage string
// }

// type ResponseMessageBPJSTpl struct {
// 	Qty             string
// 	Admin           string
// 	Periode         string
// 	CustomerName    string
// 	CustomerNumber  string
// 	SupplierMessage string
// }

// Prepaid
const (
	ResponseMsgSuccessPrepaid = iota

	ResponseMsgTransactionPendingPrepaid
	ResponseMsgTransactionProcessPrepaid
	ResponseMsgTransactionFailedPrepaid

	ResponseMsgTransactionCheckSuccessPrepaid
	ResponseMsgTransactionCheckPendingPrepaid
	ResponseMsgTransactionCheckFailedPrepaid
)

// Postpaid
const (
	ResponseMsgSuccessPostpaid = iota

	ResponseMsgInquirySuccessPostpaid
	ResponseMsgInquiryFailedPostpaid

	ResponseMsgPaymentProcessPostpaid
	ResponseMsgPaymentSuccessPostpaid
	ResponseMsgPaymentFailedPostpaid
)

var mapResponse = map[string]map[string]CategoryTpl{
	"5f0579902e6b8601c108c874": {
		"inquiry": {
			Label:   "PLN",
			Message: "{{if eq .ProductScheme \"P\"}}{{if eq .Status \"success\"}}NAMA:{{.CustomerName}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}{{- else}}{{if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TD:{{.TarifDaya}}#KWH:{{.Kwh}}#MSN:{{.Msn}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}{{- end}}",
		},
		"payment": {
			Label:   "PLN",
			Message: "{{if eq .ProductScheme \"P\"}}{{if eq .Status \"process\"}}STATUS:PROSES{{- else if eq .Status \"success\"}}NAMA:{{.CustomerName}}{{.Other}}#STATUS:BERHASIL{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}{{- else}}{{if eq .Status \"process\"}}STATUS:PROSES{{- else if eq .Status \"success\"}}NAMA:{{.CustomerName}}#DENOM:{{.Denom}}#PPN:{{.Ppn}}#PPJU:{{.Ppj}}#KWH:{{.Kwh}}#RPTOKEN:{{.RpToken}}#TOKEN:{{.Token}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}{{- end}}",
		},
	},
	"5f0579902e6b8601c108c875": {
		"inquiry": {
			Label:   "BPJS",
			Message: "{{- if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
		"payment": {
			Label:   "BPJS",
			Message: "{{if eq .Status \"process\"}}STATUS:PROSES{{- else if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
	},
	"5f0579902e6b8601c108c876": {
		"inquiry": {
			Label:   "TV BERLAGGANAN",
			Message: "{{- if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
		"payment": {
			Label:   "TV BERLAGGANAN",
			Message: "{{if eq .Status \"process\"}}STATUS:PROSES{{- else if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
	},
	"5f0579902e6b8601c108c877": {
		"inquiry": {
			Label:   "GAME ONLINE",
			Message: "{{- if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
		"payment": {
			Label:   "GAME ONLINE",
			Message: "{{if eq .Status \"process\"}}STATUS:PROSES{{- else if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
	},
	"5f0579902e6b8601c108c87c": {
		"inquiry": {
			Label:   "PDAM",
			Message: "{{- if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
		"payment": {
			Label:   "PDAM",
			Message: "{{if eq .Status \"process\"}}STATUS:PROSES{{- else if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
	},
	"5f0579902e6b8601c108c87d": {
		"inquiry": {
			Label:   "TELKOM",
			Message: "{{- if eq .Status \"success\"}}NAMA:{{.CustomerName}}#PERIOD:{{.Periode}}#TAGIHAN:{{.Tag}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
		"payment": {
			Label:   "TELKOM",
			Message: "{{if eq .Status \"process\"}}STATUS:PROSES{{- else if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
	},
	"5f0579902e6b8601c108c87f": {
		"inquiry": {
			Label:   "PERTAGAS",
			Message: "{{- if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
		"payment": {
			Label:   "PERTAGAS",
			Message: "{{if eq .Status \"process\"}}STATUS:PROSES{{- else if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
	},
	"5f0579902e6b8601c108c880": {
		"inquiry": {
			Label:   "MULTIFINANCE",
			Message: "{{- if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
		"payment": {
			Label:   "MULTIFINANCE",
			Message: "{{if eq .Status \"process\"}}STATUS:PROSES{{- else if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
	},
	"5f0579902e6b8601c108c882": {
		"inquiry": {
			Label:   "TELEPON PASCABAYAR",
			Message: "{{- if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
		"payment": {
			Label:   "TELEPON PASCABAYAR",
			Message: "{{if eq .Status \"process\"}}STATUS:PROSES{{- else if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
	},
	"5f0579902e6b8601c108c885": {
		"inquiry": {
			Label:   "PBB",
			Message: "{{- if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
		"payment": {
			Label:   "PBB",
			Message: "{{if eq .Status \"process\"}}STATUS:PROSES{{- else if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
	},
	"5f0579902e6b8601c108c886": {
		"inquiry": {
			Label:   "PGN",
			Message: "{{- if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
		"payment": {
			Label:   "PGN",
			Message: "{{if eq .Status \"process\"}}STATUS:PROSES{{- else if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
	},
	"5f0579902e6b8601c108c888": {
		"inquiry": {
			Label:   "SAMSAT ONLINE",
			Message: "{{- if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
		"payment": {
			Label:   "SAMSAT ONLINE",
			Message: "{{if eq .Status \"process\"}}STATUS:PROSES{{- else if eq .Status \"success\"}}NAMA:{{.CustomerName}}#TAGIHAN:{{.Tag}}{{.Other}}{{- else}}STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}{{- end}}",
		},
	},
}

var responseMsgPostpaidTpl = map[int]interface{}{
	ResponseMsgInquirySuccessPostpaid: mapResponse,
	ResponseMsgInquiryFailedPostpaid:  mapResponse,

	ResponseMsgPaymentProcessPostpaid: mapResponse,
	ResponseMsgPaymentSuccessPostpaid: mapResponse,
	ResponseMsgPaymentFailedPostpaid:  mapResponse,
}

var responseMsgTpl = map[int]string{
	ResponseMsgSuccessPrepaid:            "PRODUK:{{.CodeProduct}}#TUJUAN:{{.CustomerNumber}}#STATUS:BERHASIL#HARGA:{{.SellingPrice}}",
	ResponseMsgTransactionProcessPrepaid: "PRODUK:{{.CodeProduct}}#TUJUAN:{{.CustomerNumber}}#STATUS:PROSES#HARGA:{{.SellingPrice}}",
	ResponseMsgTransactionPendingPrepaid: "PRODUK:{{.CodeProduct}}#TUJUAN:{{.CustomerNumber}}#STATUS:PENDING#HARGA:{{.SellingPrice}}",
	ResponseMsgTransactionFailedPrepaid:  "PRODUK:{{.CodeProduct}}#TUJUAN:{{.CustomerNumber}}#STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}#REFUND:{{.SellingPrice}}",

	ResponseMsgTransactionCheckSuccessPrepaid: "TRANSAKSI:{{.PartnerTrxID}}#PRODUK:{{.CodeProduct}}#TUJUAN:{{.CustomerNumber}}#SUDAH PERNAH TERJADI PADA:{{.CreatedAt}}#STATUS:BERHASIL",
	ResponseMsgTransactionCheckPendingPrepaid: "TRANSAKSI:{{.PartnerTrxID}}#PRODUK:{{.CodeProduct}}#TUJUAN:{{.CustomerNumber}}#SUDAH PERNAH TERJADI PADA:{{.CreatedAt}}#STATUS:PENDING",
	ResponseMsgTransactionCheckFailedPrepaid:  "TRANSAKSI:{{.PartnerTrxID}}#PRODUK:{{.CodeProduct}}#TUJUAN:{{.CustomerNumber}}#SUDAH PERNAH TERJADI PADA:{{.CreatedAt}}#STATUS:GAGAL#KETERANGAN:{{.SupplierMessage}}",
}

var mapMaiErrorName = map[string]interface{}{
	"product": map[string]interface{}{
		"rc": RCTransactionFailed,
		"includes": []string{
			"category-unactive",
			"group-unactive",
			"product-unactive",
			"product-not-found",
			"max-queue-reached",
		},
	},
	"customer": map[string]interface{}{
		"rc": RCTransactionInputInvalid,
		"includes": []string{
			"bad-request",
		},
	},
	"mitra": map[string]interface{}{
		"rc": http.StatusBadRequest,
		"includes": []string{
			"validation",
		},
	},
	"account": map[string]interface{}{
		"rc": RCTransactionFailed,
		"includes": []string{
			"insufficient-balance",
			"balance-not-match",
			"system-maintenance",
			"account-blocked",
		},
	},
}

func ResponseMessageGenerate(responseMsg int, data ResponseMessageTpl) (msg string, err error) {
	var tpl *template.Template
	var tplStr string
	if v, ok := responseMsgTpl[responseMsg]; ok {
		tplStr = v
	} else {
		err = errors.New("template string should not empty")
		log.Println(mailog.Error(err))
		return "", err
	}

	tpl, err = template.New("statusTpl").Parse(tplStr)
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	err = tpl.Execute(&buff, data)
	if err != nil {
		return "", err
	}

	msg = buff.String()

	return
}

func ResponseMessagePostpaidGenerate(responseMsg int, req ResponsePostpaidGenerate) (msg string, err error) {
	var tpl *template.Template
	var tplStr string
	if v, ok := responseMsgPostpaidTpl[responseMsg]; ok {
		log.Println("req.Category", req.Category)
		log.Println("req.Action", req.Action)
		tplStr = v.(map[string]map[string]CategoryTpl)[req.Category][req.Action].Message
	} else {
		err = errors.New("template string should not empty")
		log.Println(mailog.Error(err))
		return "", err
	}

	tpl, err = template.New("statusTpl").Parse(tplStr)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	var buff bytes.Buffer
	err = tpl.Execute(&buff, req.DataTpl)
	if err != nil {
		return "", err
	}

	msg = buff.String()

	return
}

func ResponseMessageError(w http.ResponseWriter, svHost string, err error, pid *string, product *string, act *string) {
	maiErr := maihelper.GrpcClient.MaiErrorDetailParse(&err)
	log.Println("Error Host: ", svHost, " Name:", maiErr.Name, ", Detail:", maiErr.Detail)

	switch maiErr.Name {
	case "grpc_deadline-exceed":
		Write(w, Tpl{
			HttpStatusCode: http.StatusBadRequest,
			Body: BodyTrx{
				Body: Body{
					MAIResponseCode: http.StatusInternalServerError,
					MAIMessage:      http.StatusText(http.StatusInternalServerError),
					MAIError: map[string]interface{}{
						"msg": err.Error(),
						"raw": err,
					},
				},
				PartnerTrxID: pid,
				ACT:          act,
				ProductCode:  product,
			},
		})
	default:
		res, err := MappingResponseFailed(maiErr, svHost, pid, product, act)
		if err != nil {
			return
		}

		Write(w, res)
	}
}

func MappingResponseFailed(maiErr maigrpc.MaiErrorDetail, host string, pid *string, product *string, act *string) (res Tpl, err error) {
	se, err := jsonparse.StringToJSON(maiErr.StackEntries)
	if err != nil {
		err = maihelper.GrpcServer.MaiErrorDetail("bad-request", "Error parse stack entries : "+err.Error(), nil)
		return
	}

	var found bool

	for i := range mapMaiErrorName {
		x := mapMaiErrorName[i]
		RC := x.(map[string]interface{})["rc"].(int)
		includes := x.(map[string]interface{})["includes"].([]string)
		_, found = InSlice(includes, maiErr.Name)

		if found {

			if len(se.([]map[string]interface{})) > 0 {
				res = Tpl{
					HttpStatusCode: http.StatusBadRequest,
					Body: BodyTrx{
						Body: Body{
							MAIResponseCode: RC,
							MAIStatus:       "failed",
							MAIMessage:      http.StatusText(http.StatusBadRequest),
						},
						MAIData:      se,
						PartnerTrxID: pid,
						ACT:          act,
						ProductCode:  product,
					},
				}
			} else {
				res = Tpl{
					HttpStatusCode: http.StatusOK,
					Body: BodyTrx{
						Body: Body{
							MAIResponseCode: RC,
							MAIMessage:      maiErr.Detail,
						},
						PartnerTrxID: pid,
						ACT:          act,
						ProductCode:  product,
					},
				}
			}

			return
		}
	}

	if !found {
		res = Tpl{
			HttpStatusCode: http.StatusInternalServerError,
			Body: BodyTrx{
				Body: Body{
					MAIResponseCode: http.StatusInternalServerError,
					MAIMessage:      maiErr.Detail,
				},
			},
		}
	}

	return
}

func InSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
