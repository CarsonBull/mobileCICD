package pipelines

import (
	"context"
	"encoding/json"
	"github.com/CarsonBull/mobileCICD/scheduler/support"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"net/http"
)

// postPipeline godoc
// @Summary Create a pipeline
// @Description Create a pipeline. Pipeline will be returned after
// @Tags Pipeline
// @ID postPipeline
// @Accept json
// @produce json
// @Param Authorization header string true "Authentication header"
// @Param account body pipelines.PipelineObject true "Add Pipeline"
// @Success 201 {object} pipelines.PipelineObject
// @Router /pipelines/v1/pipeline [post]
func postPipelineEndpoint(svc PipelineService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {

		pipelineRequest := request.(map[string]interface{})

		pipeline, err := svc.CreatePipeline(pipelineRequest["name"].(string), pipelineRequest["owner"].(string))

		if err != nil {
			return support.Errorer{err}, nil
		}

		return pipeline, nil
	}
}

// getPipeline godoc
// @Summary Get a pipeline
// @Description Get a pipeline. Pipeline will be returned after
// @Tags Pipeline
// @ID getPipeline
// @produce json
// @Param Authorization header string true "Authentication header"
// @Param ID path string true "Pipeline ID"
// @Success 200
// @Router /pipelines/v1/pipeline/{ID} [get]
func getPipelineEndpoint(svc PipelineService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {

		pipelineID := request.(string)

		pipeline, err := svc.GetPipeline(pipelineID)

		if err != nil {
			return support.Errorer{err}, nil
		}

		return pipeline, nil
	}
}

// deletePipeline godoc
// @Summary Delete a pipeline
// @Description Delete a pipeline.
// @Tags Pipeline
// @ID deletePipeline
// @produce json
// @Param Authorization header string true "Authentication header"
// @Param ID path string true "Pipeline ID"
// @Success 200
// @Router /pipelines/v1/pipeline/{ID} [delete]
func deletePipelineEndpoint(svc PipelineService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {

		pipelineID := request.(string)

		err := svc.DeletePipeline(pipelineID)

		if err != nil {
			return support.Errorer{err}, nil
		}
		return struct {Message string `json:"msg"`}{"Success"}, nil
	}
}

func decodePostPipelineRequest(_ context.Context, r *http.Request) (interface{}, error) {


	var pipeline map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&pipeline); err != nil {
		return err, nil
	}
	//
	//account.Joined = time.Now()
	// TODO: add auth to this request
	//request := userRequest{ids: ids{}, object:account}

	return pipeline, nil
}

func decodeGetPipelineRequest(_ context.Context, r *http.Request) (interface{}, error) {

	urlParams := mux.Vars(r)
	ID := urlParams["ID"]

	return ID, nil
}
