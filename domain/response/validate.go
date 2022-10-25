package response

type ErrorValidation struct {
	FailedField string `json:"failed_field"`
	Rule        string `json:"rule"`
	Value       string `json:"value"`
}
