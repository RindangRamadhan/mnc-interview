package h2hrequest

// ListOfCategory getting list of category request
type ListOfCategory struct {
	PackageId string `json:"package_id" validate:"required"`
}
