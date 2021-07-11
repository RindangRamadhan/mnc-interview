package h2hresponse

const (
	RCSuccess    = 200
	RCBadRequest = 400
	RCFailed     = 500

	RCTransactionProcess      = 100
	RCTransactionPending      = 130
	RCTransactionSuccess      = 230
	RCTransactionFailed       = 530
	RCTransactionInputInvalid = 430
)

var statusText = map[int]string{
	RCSuccess: "Berhasil",
	RCFailed:  "Internal Server Error",

	RCTransactionInputInvalid: "Invalid Request",

	RCTransactionProcess: "Transaksi Sedang Di Proses",

	RCTransactionPending: "Transaksi Pending",

	RCTransactionSuccess: "Transaksi Berhasil",

	RCTransactionFailed: "Transaksi Gagal",
}

func StatusText(code int) string {
	return statusText[code]
}

func StatusName(code int) string {
	switch {
	case code == 100:
		return "PROCESS"
	case code > 100 && code < 200:
		return "PENDING"
	case code >= 200 && code < 300:
		return "SUCCESS"
	case code >= 400 && code < 499:
		return "BAD REQUEST"
	default:
		return "FAILED"
	}
}
