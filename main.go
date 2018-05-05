package main

import (
	"net/http"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	"os"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main(){
	//zerolog.SetGlobalLevel(zerolog.DebugLevel)
	//create zerolog logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	//logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}) //for pretty printing but inefficient


	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here


	//logger := log.NewLogfmtLogger(os.Stderr)
	//svc is the interface for service
	var svc StringService
	svc = stringService{}
	//logging middleware
	svc = loggingMiddleware{logger, svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, countResult, svc}
	/*var uppercase, count endpoint.Endpoint
	uppercase = makeUppercaseEndpoint(svc)
	uppercase = loggingMiddleware(log.With(logger, "method", "uppercase"))(uppercase)
	count = makeCountEndpoint(svc)
	count = loggingMiddleware(log.With(logger, "method", "count"))(count)*/
	//each service method needs a handler
	//todo: change the httprouter
	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)
	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)
	//POST /uppercase -d '{"str":"string"}'
	http.Handle("/uppercase", uppercaseHandler)
	//POST /count -d '{"str":"string"}'
	http.Handle("/count", countHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Info().Err(http.ListenAndServe(":8081",nil)).Msg("server failed to start")
}
//request and response decoder/encoders
func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {

	return json.NewEncoder(w).Encode(response)
}
