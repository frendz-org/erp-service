package product

import (
	"erp-service/iam/product/contract"
	"erp-service/iam/product/internal"
)

type Usecase = contract.Usecase

func NewUsecase(repo contract.ProductRepository, cache contract.Cache) Usecase {
	return internal.NewUsecase(repo, cache)
}
