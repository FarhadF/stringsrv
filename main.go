package main

import (
	"net/http"
	httptransport "github.com/go-kit/kit/transport/http"
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	"os"
)

func main(){
	//zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	//logger := log.NewLogfmtLogger(os.Stderr)
	var svc StringService
	svc = stringService{}
	svc = loggingMiddleware{logger, svc}
	/*var uppercase, count endpoint.Endpoint
	uppercase = makeUppercaseEndpoint(svc)
	uppercase = loggingMiddleware(log.With(logger, "method", "uppercase"))(uppercase)
	count = makeCountEndpoint(svc)
	count = loggingMiddleware(log.With(logger, "method", "count"))(count)*/
	//each service method needs a handler
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
