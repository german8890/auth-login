package domain

const (
	ClientId      = "Client-Id"
	Country       = "Country"
	TransactionId = "Transaction-Id"
)

type CustomerHeader struct {
	Country          string
	Brand            string
	StoreRef         string
	ConsumerDateTime string
	ProcessRef       string
	ChannelRef       string
	ConsumerRef      string
	Environment      string
	TypeProduct      string
	TypeProcessRef   string
	UserTx           string
}
