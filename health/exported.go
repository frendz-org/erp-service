package health

import "iam-service/health/internal"

type Usecase interface {
	CheckHealth() error
}

func NewUsecase() Usecase {
	return internal.NewUsecase()
}
