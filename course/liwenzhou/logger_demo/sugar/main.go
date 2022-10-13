package main

import (
	"net/http"

	"go.uber.org/zap"
)

var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

func InitLogger() {
	logger, _ = zap.NewDevelopment()
	sugarLogger = logger.Sugar()
}

func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Error(
			"Error fetching url..",
			zap.String("url", url),
			zap.Error(err),
		)
	} else {
		sugarLogger.Info(
			"Success...",
			zap.String("statusCode", resp.Status),
			zap.String("url", url),
		)
		resp.Body.Close()
	}
}

func main() {
	InitLogger()
	defer sugarLogger.Sync()

	simpleHttpGet("www.sogou.com")
	simpleHttpGet("http://www.sogou.com")
}
