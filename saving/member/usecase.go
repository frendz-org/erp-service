package member

import "context"

type Usecase interface {
	RegisterMember(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	GetMyMember(ctx context.Context, req *GetMyMemberRequest) (*MyMemberResponse, error)
	ListMembers(ctx context.Context, req *ListRequest) (*ListResponse, error)
	GetMember(ctx context.Context, req *GetMemberRequest) (*MemberDetailResponse, error)
	ApproveMember(ctx context.Context, req *ApproveRequest) (*MemberDetailResponse, error)
	RejectMember(ctx context.Context, req *RejectRequest) (*MemberDetailResponse, error)
	ChangeRole(ctx context.Context, req *ChangeRoleRequest) (*MemberDetailResponse, error)
	DeactivateMember(ctx context.Context, req *DeactivateRequest) (*MemberDetailResponse, error)
}
