package participant

import (
	"context"
	"fmt"
	"math"

)

func (uc *usecase) ListParticipants(ctx context.Context, req *ListParticipantsRequest) (*ListParticipantsResponse, error) {
	filter := &ParticipantFilter{
		TenantID:  req.TenantID,
		ProductID: req.ProductID,
		Status:    req.Status,
		Search:    req.Search,
		Page:      req.Page,
		PerPage:   req.PerPage,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	participants, total, err := uc.participantRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list participants: %w", err)
	}

	summaries := make([]ParticipantSummaryResponse, 0, len(participants))
	for _, p := range participants {
		summaries = append(summaries, ParticipantSummaryResponse{
			ID:             p.ID,
			FullName:       p.FullName,
			KTPNumber:      p.KTPNumber,
			EmployeeNumber: p.EmployeeNumber,
			PhoneNumber:    p.PhoneNumber,
			Status:         string(p.Status),
			SubmittedAt:    p.SubmittedAt,
			ApprovedAt:     p.ApprovedAt,
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.PerPage)))

	return &ListParticipantsResponse{
		Participants: summaries,
		Pagination: PaginationMeta{
			Page:       req.Page,
			PerPage:    req.PerPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
