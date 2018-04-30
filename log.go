package main

import (
	"github.com/rs/zerolog"
	//"github.com/rs/zerolog/log"
	"context"
	"time"
)

//func loggingMiddleware(logger log.Logger) endpoint.Middleware {
//	return func (next endpoint.Endpoint) endpoint.Endpoint {
//		return func(ctx context.Context, request interface{}) (interface{}, error){
//			logger.Log("msg", "calling endpoint")
//			defer logger.Log("msg", "called endpoint")
//			return next(ctx, request)
//		}
//	}
//}

//struct passing the logger
type loggingMiddleware struct {
	logger zerolog.Logger
	next 	StringService
}

//each method will have its own logger for app logs
func (mw loggingMiddleware)Uppercase(ctx context.Context, s string) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Info().Str(
			"method", "uppercase").Str("input", s).Str("output", output).Err(err).Dur("took",
				time.Since(begin)).Msg("")

	}(time.Now())
	output, err = mw.next.Uppercase(ctx, s)
	return
}

//each method will have its own logger for app logs
func (mw loggingMiddleware)Count(ctx context.Context, s string) (output int) {
	defer func(begin time.Time) {
		mw.logger.Info().Str("method", "uppercase").Str("input", s).Int("output", output).Dur("took",
			time.Since(begin)).Msg("")
	}(time.Now())
	output = mw.next.Count(ctx, s)
	return
}