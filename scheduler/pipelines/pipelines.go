package pipelines

import (
	"time"
)

type PipelineService interface {
	CreatePipeline(name string, owner string) (PipelineObject, error)
	GetPipeline(ID string) (PipelineObject, error)
	DeletePipeline(ID string) error
}

type Service struct {}

type PipelineObject struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
	Name string `json:"name"`
	// TODO: Owner should be its own datatype
	Owner string `json:"owner"`
	//// TODO: steps should be its own datatype
	//Steps []string `json:"steps"`
}

func (s Service) CreatePipeline(name string, owner string) (PipelineObject, error) {

	var pipeObject = PipelineObject{Name:name, Owner:owner}

	err := DB.Create(&pipeObject).Error

	if err != nil {
		return PipelineObject{}, err
	}

	return pipeObject, nil
}

func (s Service) GetPipeline(ID string) (PipelineObject, error) {

	var pipelineObject PipelineObject

	err := DB.Where("id = ?", ID).Find(&pipelineObject).Error

	if err != nil {
		return PipelineObject{}, err
	}

	return pipelineObject, nil

}

func (s Service) DeletePipeline(ID string) error {

	var pipelineObject PipelineObject

	err := DB.Where("id = ?", ID).Delete(&pipelineObject).Error

	if err != nil {
		return err
	}

	return nil

}