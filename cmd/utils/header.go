package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	CountryHeader          = "X-country"
	ConsumerRefHeader      = "X-consumerRef"
	ChannelRefHeader       = "X-channelRef"
	BrandHeader            = "X-brand"
	StoreRefHeader         = "X-storeRef"
	UserTrxHeader          = "X-userTrx"
	EnvHeader              = "X-environment"
	ConsumerDateTimeHeader = "X-consumerDateTime"
	ProcessRefHeader       = "X-processRef"
	TypeProductHeader      = "X-typeProduct"
	TypeProcessRefHeader   = "X-typeProcessRef"

	IdempotencyKeyHeader = "X-Idempotency-key"

	B3TraceIdHeader  = "X-B3-TraceId"
	B3SpanIdHeader   = "X-B3-SpanId"
	B3ParentIdHeader = "X-B3-ParentSpanId"
	Apikey           = "X-Apikey"

	OriginRequested = "X-Origin-Requested"
	AcceptLanguage  = "X-Accept-Language"
)

var HeadersToMaskerIntegration = []string{Apikey, "Authorization"}

func GetCountryHeader(ctx *gin.Context) string {
	if country := ctx.Request.Header.Get(CountryHeader); country != "" {
		return country
	}
	return "JM"
}

func HeadersToMap(r *http.Request) map[string]string {
	var mapHeader = make(map[string]string)
	for name := range r.Header {
		mapHeader[name] = r.Header.Get(name)
	}
	return mapHeader
}
