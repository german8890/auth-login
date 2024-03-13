package http_client

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type CustomTransport struct {
	RoundTripper http.RoundTripper
	propagator   propagation.TextMapPropagator
}

func NewCustomTransport(upstream *http.Transport) *CustomTransport {
	//upstream.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return &CustomTransport{RoundTripper: upstream, propagator: otel.GetTextMapPropagator()}
}
func (ct *CustomTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	var (
		childSpan trace.Span
		tracer2   trace.Tracer
	)
	ct.propagator.Inject(req.Context(), propagation.HeaderCarrier(req.Header))

	if parentSpan := trace.SpanFromContext(req.Context()); parentSpan.SpanContext().IsValid() {

		opts := []oteltrace.SpanStartOption{
			//oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", req)...),
			oteltrace.WithAttributes(semconv.HTTPClientAttributesFromHTTPRequest(req)...),
			oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(req)...),

			oteltrace.WithSpanKind(oteltrace.SpanKindClient),
		}
		spanName := req.URL.String()
		tracer2 = parentSpan.TracerProvider().Tracer(spanName)
		ctx := req.Context()
		_, childSpan = tracer2.Start(ctx, req.URL.String(), opts...)

		defer func() {
			childSpan.SetAttributes(semconv.HTTPServerAttributesFromHTTPRequest(req.URL.Host, req.URL.Path, req)...)
			if resp != nil {
				childSpan.SetAttributes(semconv.HTTPAttributesFromHTTPStatusCode(resp.StatusCode)...)
				childSpan.SetStatus(semconv.SpanStatusFromHTTPStatusCode(resp.StatusCode))

			}
			if err != nil {
				childSpan.RecordError(err)
				childSpan.SetStatus(codes.Error, err.Error())
				childSpan.End()
			} else {
				childSpan.End()
			}
		}()
	}

	resp, err = ct.RoundTripper.RoundTrip(req)

	return resp, err
}
