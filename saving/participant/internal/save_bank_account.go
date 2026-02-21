package internal

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/participant/participantdto"
)

func (uc *usecase) SaveBankAccount(ctx context.Context, req *participantdto.SaveBankAccountRequest) (*participantdto.BankAccountResponse, error) {
	var result *participantdto.BankAccountResponse

	err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		participant, err := uc.participantRepo.GetByID(txCtx, req.ParticipantID)
		if err != nil {
			return fmt.Errorf("get participant: %w", err)
		}

		if err := validateParticipantOwnership(participant, req.TenantID, req.ProductID); err != nil {
			return err
		}

		if err := validateEditableState(participant); err != nil {
			return err
		}

		if req.IsPrimary {
			if err := uc.bankAccountRepo.ClearPrimary(txCtx, req.ParticipantID); err != nil {
				return fmt.Errorf("clear primary bank accounts: %w", err)
			}
		}

		var account *entity.ParticipantBankAccount

		if req.ID != nil {
			account, err = uc.bankAccountRepo.GetByID(txCtx, *req.ID)
			if err != nil {
				return fmt.Errorf("get bank account: %w", err)
			}

			if account.ParticipantID != req.ParticipantID {
				return errors.ErrForbidden("bank account does not belong to this participant")
			}

			account.BankCode = req.BankCode
			account.AccountNumber = req.AccountNumber
			account.AccountHolderName = req.AccountHolderName
			account.AccountType = req.AccountType
			account.CurrencyCode = req.CurrencyCode
			account.IsPrimary = req.IsPrimary
			account.IssueDate = req.IssueDate
			account.ExpiryDate = req.ExpiryDate

			if err := uc.bankAccountRepo.Update(txCtx, account); err != nil {
				return fmt.Errorf("update bank account: %w", err)
			}
		} else {
			now := time.Now()
			account = &entity.ParticipantBankAccount{
				ParticipantID:     req.ParticipantID,
				BankCode:          req.BankCode,
				AccountNumber:     req.AccountNumber,
				AccountHolderName: req.AccountHolderName,
				AccountType:       req.AccountType,
				CurrencyCode:      req.CurrencyCode,
				IsPrimary:         req.IsPrimary,
				IssueDate:         req.IssueDate,
				ExpiryDate:        req.ExpiryDate,
				Version:           1,
				CreatedAt:         now,
				UpdatedAt:         now,
			}

			if err := uc.bankAccountRepo.Create(txCtx, account); err != nil {
				return fmt.Errorf("create bank account: %w", err)
			}
		}

		resp := mapBankAccountToResponse(account)
		result = &resp
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
