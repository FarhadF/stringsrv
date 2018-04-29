package main

import (
	"net/http"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/kit/log"
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"os"
)

func main(){
	svc := stringService{}
	logger := log.NewLogfmtLogger(os.Stderr)
	var uppercase endpoint.Endpoint
	uppercase = makeUppercaseEndpoint(svc)
	uppercase = loggingMiddleware(log.With(logger, "method", "uppercase"))(uppercase)
	//each service method needs a handler
	uppercaseHandler := httptransport.NewServer(
		uppercase,
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
	logger.Log("err",http.ListenAndServe(":8081",nil))
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
