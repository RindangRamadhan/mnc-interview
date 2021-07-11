package h2hrequest

// CreateTrxPrepaid transaction prepaid request
type TransactionPrepaid struct {
	CustomerId   string `json:"customer_id" validate:"required"`
	ProductCode  string `json:"product_code" validate:"required"`
	PartnerTrxId string `json:"partner_trx_id" validate:"required"`
}

// CreateTrxPostpaid transaction postpaid request
type TransactionPostpaid struct {
	Act          string `json:"act" validate:"required"` // Action for "inquiry" or "payment"
	Method       string `json:"method"`                  // Method is "sync" or "async"
	Qty          string `json:"qty"`
	CustomerId   string `json:"customer_id" validate:"required"`
	ProductCode  string `json:"product_code" validate:"required"`
	PartnerTrxId string `json:"partner_trx_id" validate:"required"`
}

type DeviceInfo struct {
	DeviceName string `json:"device_name"`
}

type H2HPayload struct {
	PartnerTrxID string `json:"partner_trx_id"`
	ProductCode  string `json:"product_code"`
	CallbackURL  string `json:"callback_url"`
	Act          string `json:"act"`
	Method       string `json:"method"`
	RespData     string `json:"resp_data"`
}

// CreateTrxPrepaid for creating invoice prepaid
type CreateTrxPrepaid struct {
	Qty               int         `json:"qty"`
	Zonid             string      `json:"zonid"`
	GroupId           int         `json:"group_id"`
	MemberId          int         `json:"member_id"`
	Reference         string      `json:"reference"`
	TotalPrice        float64     `json:"total_price"`
	ProductId         string      `json:"product_id"`
	PackageId         string      `json:"package_id"`
	CustomerId        string      `json:"customer_id"`
	H2HPayload        *H2HPayload `json:"h2h_payload"`
	DeviceInfo        DeviceInfo  `json:"device_info"`
	WaitingTime       int64       `json:"waiting_time"`
	ProductHppId      string      `json:"product_hpp_id"`
	PaymentMethodId   int         `json:"payment_method_id"`
	PaymentMethodType string      `json:"payment_method_type"`
}

// InquiryTrxPostpaid for creating inquiry postpaid
type InquiryTrxPostpaid struct {
	Qty               int         `json:"qty"`
	Zonid             string      `json:"zonid"`
	GroupId           int         `json:"group_id"`
	MemberId          int         `json:"member_id"`
	Reference         string      `json:"reference"`
	TotalPrice        float64     `json:"total_price"`
	ProductId         string      `json:"product_id"`
	PackageId         string      `json:"package_id"`
	CustomerId        string      `json:"customer_id"`
	H2HPayload        *H2HPayload `json:"h2h_payload"`
	DeviceInfo        DeviceInfo  `json:"device_info"`
	WaitingTime       int64       `json:"waiting_time"`
	ProductHppId      string      `json:"product_hpp_id"`
	PaymentMethodId   int         `json:"payment_method_id"`
	PaymentMethodType string      `json:"payment_method_type"`
}

// CreateTrxPostpaid for creating invoice postpaid
type CreateTrxPostpaid struct {
	Qty               int         `json:"qty"`
	Zonid             string      `json:"zonid"`
	GroupId           int         `json:"group_id"`
	MemberId          int         `json:"member_id"`
	Reference         string      `json:"reference"`
	TotalPrice        float64     `json:"total_price"`
	ProductId         string      `json:"product_id"`
	PackageId         string      `json:"package_id"`
	CustomerId        string      `json:"customer_id"`
	H2HPayload        *H2HPayload `json:"h2h_payload"`
	DeviceInfo        DeviceInfo  `json:"device_info"`
	WaitingTime       int64       `json:"waiting_time"`
	PaymentMethodId   int         `json:"payment_method_id"`
	PaymentMethodType string      `json:"payment_method_type"`
}

// DirectPaymentPostpaid for creating inquiry & invoice postpaid
type DirectPaymentPostpaid struct {
	Qty               int         `json:"qty"`
	Zonid             string      `json:"zonid"`
	GroupId           int         `json:"group_id"`
	MemberId          int         `json:"member_id"`
	Reference         string      `json:"reference"`
	TotalPrice        float64     `json:"total_price"`
	ProductId         string      `json:"product_id"`
	PackageId         string      `json:"package_id"`
	CustomerId        string      `json:"customer_id"`
	H2HPayload        *H2HPayload `json:"h2h_payload"`
	DeviceInfo        DeviceInfo  `json:"device_info"`
	WaitingTime       int64       `json:"waiting_time"`
	ProductHppId      string      `json:"product_hpp_id"`
	PaymentMethodId   int         `json:"payment_method_id"`
	PaymentMethodType string      `json:"payment_method_type"`
}
