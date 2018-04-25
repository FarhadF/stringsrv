package main

import (
	"github.com/go-kit/kit/endpoint"
	"context"
)
//An endpoint represents a single RPC. That is, a single method in our service interface. We’ll write simple adapters to
//convert each of our service’s methods into an endpoint. Each adapter takes a StringService, and returns an endpoint
//that corresponds to one of the methods.
func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func (ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(ctx, req.Str)
		if err != nil {
			return uppercaseResponse{v, err}, nil
		}
		return uppercaseResponse{v, err}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func (ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		i:= svc.Count(ctx, req.Str)
		return countResponse{i}, nil
	}
}

//For each method, we define request and response structs, capturing all of the input and output parameters respectively.
type uppercaseRequest struct{
	Str string `json:"str"`
}

type uppercaseResponse struct{
	Str string `json:"str"`
	Err error `json:"err"`
}

type countRequest struct{
	Str string `json:"str"`
}

type countResponse struct{
	Count int `json:"count"`
}


