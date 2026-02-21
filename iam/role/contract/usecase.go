package contract

import (
	"context"
	"erp-service/iam/role/roledto"
)

type Usecase interface {
	Create(ctx context.Context, req *roledto.CreateRequest) (*roledto.CreateResponse, error)
}
