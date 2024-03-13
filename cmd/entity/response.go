package entity

type Response struct {
	Data   any    `json:"data,omitempty" mask:"struct"`
	Result Result `json:"result"`
}

type ResponseWithList struct {
	Data   []interface{} `json:"data,omitempty" mask:"struct"`
	Result Result        `json:"result"`
}

type Result struct {
	Details []Detail `json:"details"`
	Source  string   `json:"source"`
}

type Detail struct {
	InternalCode string `json:"internalCode"`
	Message      string `json:"message"`
	Detail       string `json:"detail"`
}
