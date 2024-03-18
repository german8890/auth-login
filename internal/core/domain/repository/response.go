package repository_models

type ProcessEvent struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type ResponseTopicGCPHeaders struct {
	XAcceptLanguage string `json:"X-Accept-Language"`
	XB3SpanID       string `json:"X-B3-SpanId"`
	XB3TraceID      string `json:"X-B3-TraceId"`
	XIdempotencyKey string `json:"X-Idempotency-Key"`
	XBrand          string `json:"X-brand"`
	XChannelRef     string `json:"X-channelRef"`
	XConsumerRef    string `json:"X-consumerRef"`
	XEnviorment     string `json:"X-enviroment"`
	XProcessRef     string `json:"X-processRef"`
	XStoreRef       string `json:"X-storeRef"`
	XTypeProduct    string `json:"X-typeProduct"`
	Country         string `json:"country"`
}

type ResponseTopicGCP struct {
	Capability        string `json:"capability"`
	ConsumerDateTime  string `json:"consumerDateTime"`
	DataSpecVersion   string `json:"dataSpecVersion"`
	DataContentType   string `json:"datacontenttype"`
	Domain            string `json:"domain"`
	EventName         string `json:"eventName"`
	LoanAccountId     string `json:"loanAccountId"`
	LoanTransactionId string `json:"loanTransactionId"`
	Producer          string `json:"producer"`
	TimestampEvent    string `json:"timestampEvent"`
}
