package h2hresponse

import (
	"encoding/json"

	invoicedetail "gitlab.com/pt-mai/transaction-service/pkg/server/invoice-detail"
)

type BodyTrx struct {
	PartnerTrxID *string `json:"partner_trx_id,omitempty"`
	Body
	ProductCode  *string     `json:"product_code,omitempty"`
	CustomerID   *string     `json:"customer_id,omitempty"`
	ACT          *string     `json:"act,omitempty"`
	SN           *string     `json:"SN,omitempty"`
	Admin        *float64    `json:"admin,omitempty"`
	TotalPayment *float64    `json:"total_payment,omitempty"`
	PotSaldo     *float64    `json:"pot_saldo,omitempty"`
	RefundSaldo  *float64    `json:"refund_saldo,omitempty"`
	LastBalance  *float64    `json:"last_balance,omitempty"`
	MAIData      interface{} `json:"data,omitempty"`
}

type BodyTrxData struct {
	ZonID           string               `json:"zonid"`
	InvNumber       string               `json:"inv_number"`
	PaymentMethod   int                  `json:"payment_method"`
	InvDetailStatus invoicedetail.Status `json:"inv_detail_status"`
	Price           int                  `json:"std_customer_price"`
	Discount        int                  `json:"discount"`
	DiscountType    uint32               `json:"discount_type"`
	TotalPayment    int                  `json:"total_payment"`
	Product         json.RawMessage      `json:"product"`
	Receipt         *BodyTrxDataReceipt  `json:"receipt"`
	CreatedAt       string               `json:"created_at"`
	UpdatedAt       string               `json:"updated_at"`
}

type BodyTrxDataReceipt struct {
	Json struct {
		Footer  []interface{}            `json:"footer"`
		Header  string                   `json:"header"`
		Content []map[string]interface{} `json:"content"`
	} `json:"json"`
	Text interface{} `json:"text"`
}
