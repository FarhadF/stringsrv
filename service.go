package main

import (
	"context"
	"errors"
	"strings"
)
//each service is an interface with service methods
type StringService interface {
	Uppercase(context.Context, string) (string, error)
	Count(context.Context, string) int
}

//implement service as an empty struct
type stringService struct {}

var errEmpty = errors.New("input string is empty")
//implement methods
func (stringService) Uppercase(_ context.Context, s string) (string, error) {
	if s == "" {
		return "", errEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(_ context.Context, s string) int {
	return len(s)
}

