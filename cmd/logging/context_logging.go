package logging

import (
	"context"
	"net/http"
	"strings"

	"autenticacion-ms/cmd/config/model"
	"autenticacion-ms/cmd/entity"
	"autenticacion-ms/cmd/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

const (
	requestIDKey contextKey = iota
	correlationIDKey
	country          = "country"
	consumerRef      = "consumerRef"
	channelRef       = "channelRef"
	brand            = "brand"
	storeRef         = "storeRef"
	userTrx          = "userTrx"
	traceId          = "traceId"
	spanId           = "spanId"
	env              = "env"
	consumerDateTime = "consumerDateTime"
	processRef       = "processRef"
	groupId          = "groupId"
	operationId      = "operationId"
	typeProduct      = "typeProduct"
	typeProcessRef   = "typeProcessRef"
	idempotencyKey   = "idempotencyKey"
)

var (
	Service entity.Service
)

func AddOperationInContext(c *gin.Context, artifact *model.ArtifactResources) *gin.Context {
	c.Set(groupId, artifact.GroupName)
	for _, resource := range artifact.Resources {
		if strings.EqualFold(c.FullPath(), resource.Path) && strings.EqualFold(c.Request.Method, resource.Method) {
			c.Set(operationId, resource.Operation)
		}
	}
	return c
}

// WithRequest returns a context which knows the request ID and correlation ID in the given request.
func WithRequest(ctx *gin.Context, req *http.Request) *gin.Context {

	//if val, _ := ctx.Get(TraceId); val == nil || val == "" {
	ctx.Set(country, getCountry(req))
	ctx.Set(consumerRef, getConsumerRef(req))
	ctx.Set(channelRef, getChannelRef(req))
	ctx.Set(brand, getBrand(req))
	ctx.Set(storeRef, getStoreRef(req))
	ctx.Set(userTrx, getUserTrx(req))
	ctx.Set(env, getEnv(req))
	ctx.Set(consumerDateTime, getConsumerDateTime(req))
	ctx.Set(processRef, getProcessRef(req))
	ctx.Set(traceId, getTraceId(req))
	ctx.Set(spanId, getSpanId(req))
	if val, valid := ctx.Get(utils.IdempotencyKey); !valid && val == nil {
		ctx.Set(utils.IdempotencyKey, getIdempotencyKeyHeader(req))
	}
	ctx.Set(typeProduct, getTypeProductHeader(req))
	ctx.Set(typeProcessRef, getTypeProcessRefHeader(req))
	//}
	return ctx
}

// WithRequestV2 returns a context which knows the request ID and correlation ID in the given request.
func WithRequestV2(ctx context.Context, req *http.Request) context.Context {
	//if val := ctx.Value(TraceId); val == nil || val == "" {
	ctx = context.WithValue(ctx, country, getCountry(req))
	ctx = context.WithValue(ctx, consumerRef, getConsumerRef(req))
	ctx = context.WithValue(ctx, channelRef, getChannelRef(req))
	ctx = context.WithValue(ctx, brand, getBrand(req))
	ctx = context.WithValue(ctx, storeRef, getStoreRef(req))
	ctx = context.WithValue(ctx, userTrx, getUserTrx(req))
	ctx = context.WithValue(ctx, env, getEnv(req))
	ctx = context.WithValue(ctx, consumerDateTime, getConsumerDateTime(req))
	ctx = context.WithValue(ctx, processRef, getProcessRef(req))
	ctx = context.WithValue(ctx, traceId, getTraceId(req))
	ctx = context.WithValue(ctx, spanId, getSpanId(req))
	if val := ctx.Value(utils.IdempotencyKey); val == nil {
		ctx = context.WithValue(ctx, idempotencyKey, getIdempotencyKeyHeader(req))
	}
	ctx = context.WithValue(ctx, typeProduct, getTypeProductHeader(req))
	ctx = context.WithValue(ctx, typeProcessRef, getTypeProcessRefHeader(req))
	//}
	return ctx
}

func getCountry(req *http.Request) string {
	return req.Header.Get("X-Country")
}

func getConsumerRef(req *http.Request) string {
	return req.Header.Get("X-consumerRef")
}

func getChannelRef(req *http.Request) string {
	return req.Header.Get("X-channelRef")
}

func getBrand(req *http.Request) string {
	return req.Header.Get("X-brand")
}

func getStoreRef(req *http.Request) string {
	return req.Header.Get("X-storeRef")
}

func getUserTrx(req *http.Request) string {
	return req.Header.Get("X-userTrx")
}

func getTypeProductHeader(req *http.Request) string {
	return req.Header.Get(utils.TypeProductHeader)
}

func getTypeProcessRefHeader(req *http.Request) string {
	return req.Header.Get(utils.TypeProcessRefHeader)
}

func getTraceId(req *http.Request) string {
	if span := trace.SpanFromContext(req.Context()); span.SpanContext().IsValid() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}

func getSpanId(req *http.Request) string {
	if span := trace.SpanFromContext(req.Context()); span.SpanContext().IsValid() {
		return span.SpanContext().SpanID().String()
	}
	return ""
}

func getEnv(req *http.Request) string {
	return req.Header.Get("X-env")
}

func getConsumerDateTime(req *http.Request) string {
	return req.Header.Get("X-consumerDateTime")
}

func getProcessRef(req *http.Request) string {
	return req.Header.Get("X-processRef")
}

func getIdempotencyKeyHeader(req *http.Request) string {
	if req.Header.Get(utils.IdempotencyKeyHeader) != "" {
		return req.Header.Get(utils.IdempotencyKeyHeader)
	}
	return uuid.New().String()
}
