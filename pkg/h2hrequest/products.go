package h2hrequest

// ListOfProductByCategoryCode getting list of category request
type ListOfProductByCategoryCode struct {
	CategoryCode string `json:"category_code" validate:"required"`
	PackageId    string `json:"package_id" validate:"required"`
}
