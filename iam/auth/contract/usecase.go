package contract

import (
	"context"

	"iam-service/masterdata/masterdatadto"
)

type MasterdataUsecase interface {
	ValidateItemCode(ctx context.Context, req *masterdatadto.ValidateCodeRequest) (*masterdatadto.ValidateCodeResponse, error)
}
