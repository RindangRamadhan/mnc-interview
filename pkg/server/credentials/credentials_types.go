package credentials

// GetAllZonidRes struct
type GetAllZonidRes struct {
	Zonid string `json:"zonid"`
}

type GetByZonidResponse struct {
	AccessKey   string   `bson:"access_key" json:"access_key"`
	SecretKey   string   `bson:"secret_key" json:"secret_key"`
	CallbackUrl string   `bson:"callback_url" json:"callback_url"`
	WhitelistIp []string `bson:"whitelist_ip" json:"whitelist_ip"`
	RespData    bool     `bson:"resp_data" json:"resp_data"`
	ApiKey      string   `bson:"api_key" json:"api_key"`
}

type UpdateRequest struct {
	Zonid       string   `json:"zonid" validate:"required"`
	MemberId    string   `json:"member_id" validate:"required"`
	PackageId   string   `json:"package_id" validate:"required"`
	ClientName  string   `json:"client_name" validate:"required"`
	CallbackUrl string   `json:"callback_url" validate:"required"`
	WhitelistIp []string `json:"whitelist_ip" validate:"required"`
	RespData    bool     `json:"resp_data"`
}
