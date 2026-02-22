package role

import "context"

type Usecase interface {
	Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error)
}
