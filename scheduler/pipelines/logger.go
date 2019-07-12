package pipelines

import (
	"github.com/go-kit/kit/log"
	"time"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next   PipelineService
}

func (mw LoggingMiddleware) CreatePipeline(name, owner string) (PipelineObject, error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "CreatePipeline",
			"took", time.Since(begin),
		)
	}(time.Now())

	pipeline, err := mw.Next.CreatePipeline(name, owner)

	return pipeline, err
}

func (mw LoggingMiddleware) GetPipeline(ID string) (PipelineObject, error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "GetPipeline",
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.Next.GetPipeline(ID)
}

func (mw LoggingMiddleware) DeletePipeline(ID string) error {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "DeletePipeline",
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.Next.DeletePipeline(ID)
}
