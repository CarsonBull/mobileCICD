package main

import (
	"github.com/CarsonBull/mobileCICD/scheduler/pipelines"
	"github.com/go-kit/kit/log"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"

	_ "github.com/CarsonBull/mobileCICD/scheduler/docs"
)

// @title Swagger Example API
// @version 1.0
// @description This is the heart cards account service

// @contact.name API Support
// @contact.url http://www.blah.com/support
// @contact.email support@blah.com

// @host localhost:8090

func main() {

	logger := log.NewLogfmtLogger(os.Stderr)

	var svc pipelines.PipelineService
	svc = pipelines.Service{}
	svc = pipelines.LoggingMiddleware{logger, svc}

	muxServer := http.NewServeMux()

	muxServer.Handle("/pipelines/v1/", pipelines.MakeHandler(svc))

	muxServer.Handle("/docs/", httpSwagger.WrapHandler)

	http.Handle("/", muxServer)

	_ = logger.Log("msg", "HTTP", "addr", ":8090")
	_ = logger.Log("err", http.ListenAndServe(":8090", nil))

}




