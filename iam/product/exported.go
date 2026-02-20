package product

import (
	"iam-service/iam/product/contract"
	"iam-service/iam/product/internal"
)

type Usecase = contract.Usecase

func NewUsecase(repo contract.ProductRepository, cache contract.Cache) Usecase {
	return internal.NewUsecase(repo, cache)
}
