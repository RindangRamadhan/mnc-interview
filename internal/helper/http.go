package helper

import (
	"log"
	"net/http"

	"gitlab.com/pt-mai/maihelper"
	jsonparse "gitlab.com/pt-mai/maihelper/utils/json-parse"
)

// HandleError ...
func HandleError(w http.ResponseWriter, err error) {
	maiErr := maihelper.GrpcClient.MaiErrorDetailParse(&err)

	log.Println("gRPC call error. Name: ", maiErr.Name, ", Detail:", maiErr.Detail)

	switch {
	case maiErr.Name == "validation":
		stackEntries, err := jsonparse.StringToJSON(maiErr.StackEntries)
		if err != nil {
			maihelper.GrpcClient.MaiHttpResponseHandler(w, http.StatusInternalServerError, "error", "Validation Error", maiErr.Detail)
			break
		}

		maihelper.GrpcClient.MaiHttpResponseHandler(w, http.StatusBadRequest, "error", maiErr.Detail, stackEntries)
		break
	case maiErr.Name == "blocked_deposit":
		stackEntries, err := jsonparse.StringToJSON(maiErr.StackEntries)
		if err != nil {
			maihelper.GrpcClient.MaiHttpResponseHandler(w, http.StatusInternalServerError, "error", "Blocked Deposit", maiErr.Detail)
			break
		}

		maihelper.GrpcClient.MaiHttpResponseHandler(w, http.StatusBadRequest, "error", maiErr.Detail, stackEntries)
		break
	case maiErr.Name == "dupplicate":
		stackEntries, err := jsonparse.StringToJSON(maiErr.StackEntries)
		if err != nil {
			maihelper.GrpcClient.MaiHttpResponseHandler(w, http.StatusConflict, "error", "Dupplicate Error", maiErr.Detail)
			break
		}

		maihelper.GrpcClient.MaiHttpResponseHandler(w, http.StatusConflict, "error", maiErr.Detail, stackEntries)
		break
	case maiErr.Name == "account-blocked":
		stackEntries, err := jsonparse.StringToJSON(maiErr.StackEntries)
		if err != nil {
			maihelper.GrpcClient.MaiHttpResponseHandler(w, http.StatusConflict, "error", "Account Blocked", maiErr.Detail)
			break
		}

		maihelper.GrpcClient.MaiHttpResponseHandler(w, http.StatusConflict, "error", maiErr.Detail, stackEntries)
		break
	case maiErr.Name == "deposit-blocked":
		stackEntries, err := jsonparse.StringToJSON(maiErr.StackEntries)
		if err != nil {
			maihelper.GrpcClient.MaiHttpResponseHandler(w, http.StatusConflict, "error", "Account Blocked", maiErr.Detail)
			break
		}

		maihelper.GrpcClient.MaiHttpResponseHandler(w, http.StatusConflict, "error", maiErr.Detail, stackEntries)
		break
	case maiErr.Name == "no-data":
		stackEntries, err := jsonparse.StringToJSON(maiErr.StackEntries)
		if err != nil {
			maihelper.GrpcClient.MaiHttpResponseHandler(w, http.StatusConflict, "error", "Data Not Found", maiErr.Detail)
			break
		}

		maihelper.GrpcClient.MaiHttpResponseHandler(w, http.StatusConflict, "error", maiErr.Detail, stackEntries)
		break
	default:
		maihelper.GrpcClient.MaiHttpResponseHandler(w, http.StatusInternalServerError, "error", "Internal Server Error", maiErr.Detail)
		break
	}
}
