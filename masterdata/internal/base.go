package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"erp-service/config"
	"erp-service/masterdata/contract"
)

type usecase struct {
	config       *config.Config
	categoryRepo contract.CategoryRepository
	itemRepo     contract.ItemRepository
	cache        contract.MasterdataCache
}

func NewUsecase(
	cfg *config.Config,
	categoryRepo contract.CategoryRepository,
	itemRepo contract.ItemRepository,
	cache contract.MasterdataCache,
) *usecase {
	return &usecase{
		config:       cfg,
		categoryRepo: categoryRepo,
		itemRepo:     itemRepo,
		cache:        cache,
	}
}

func hashFilter(filter any) string {
	data, _ := json.Marshal(filter)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:8])
}

func normalizePageParams(page, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	return page, perPage
}
