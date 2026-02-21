package auth

import (
	"context"

	"erp-service/masterdata"
)

type MasterdataUsecase interface {
	ValidateItemCode(ctx context.Context, req *masterdata.ValidateCodeRequest) (*masterdata.ValidateCodeResponse, error)
}
