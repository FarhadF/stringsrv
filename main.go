package main

import (
	"net/http"
	httptransport "github.com/go-kit/kit/transport/http"
	"log"
	"context"
	"encoding/json"
)

func main(){
	svc := stringService{}
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
	log.Fatal(http.ListenAndServe(":8080",nil))
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