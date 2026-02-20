package contract

import (
	"context"

	"iam-service/saving/member/memberdto"
)

type Usecase interface {
	RegisterMember(ctx context.Context, req *memberdto.RegisterRequest) (*memberdto.RegisterResponse, error)
	ListMembers(ctx context.Context, req *memberdto.ListRequest) (*memberdto.ListResponse, error)
	GetMember(ctx context.Context, req *memberdto.GetMemberRequest) (*memberdto.MemberDetailResponse, error)
	ApproveMember(ctx context.Context, req *memberdto.ApproveRequest) (*memberdto.MemberDetailResponse, error)
	RejectMember(ctx context.Context, req *memberdto.RejectRequest) (*memberdto.MemberDetailResponse, error)
	ChangeRole(ctx context.Context, req *memberdto.ChangeRoleRequest) (*memberdto.MemberDetailResponse, error)
	DeactivateMember(ctx context.Context, req *memberdto.DeactivateRequest) (*memberdto.MemberDetailResponse, error)
}
