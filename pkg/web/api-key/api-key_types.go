package apikey

type CredentialRequest struct {
	Zonid       string   `json:"zonid" validate:"required"`
	MemberId    string   `json:"member_id" validate:"required"`
	PackageId   string   `json:"package_id" validate:"required"`
	ClientName  string   `json:"client_name" validate:"required"`
	CallbackUrl string   `json:"callback_url" validate:"required"`
	WhitelistIp []string `json:"whitelist_ip" validate:"required"`
	RespData    bool     `json:"resp_data"`
}

type CredentialResponse struct {
	CallbackUrl string   `json:"callback_url"`
	WhitelistIp []string `json:"whitelist_ip"`
	RespData    bool     `json:"resp_data"`
	ApiKey      string   `json:"api_key"`
}
