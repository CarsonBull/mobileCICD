package pipelines

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/CarsonBull/mobileCICD/scheduler/support"
	"github.com/gorilla/mux"
	"net/http"
)

func MakeHandler(svc PipelineService) http.Handler  {

	postPipelineHandler := httptransport.NewServer(
		postPipelineEndpoint(svc),
		//auth.Authenticator{decodePostAccountRequest}.AuthenticateUser,
		decodePostPipelineRequest,
		support.EncodeResponse,
	)

	getPipelineHandler := httptransport.NewServer(
		getPipelineEndpoint(svc),
		//auth.Authenticator{decodePostAccountRequest}.AuthenticateUser,
		decodeGetPipelineRequest,
		support.EncodeResponse,
	)

	deletePipelineHandler := httptransport.NewServer(
		deletePipelineEndpoint(svc),
		//auth.Authenticator{decodePostAccountRequest}.AuthenticateUser,
		decodeGetPipelineRequest,
		support.EncodeDeleteResponse,
	)

	r := mux.NewRouter()

	r.Handle("/pipelines/v1/pipeline", postPipelineHandler).Methods("POST")
	r.Handle("/pipelines/v1/pipeline/{ID:[0-9]+}", getPipelineHandler).Methods("GET")
	r.Handle("/pipelines/v1/pipeline/{ID:[0-9]+}", deletePipelineHandler).Methods("DELETE")

	return r

}