package v1

type Example struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type SimpleSuccessResponse struct {
	Success bool `json:"success"`
}
