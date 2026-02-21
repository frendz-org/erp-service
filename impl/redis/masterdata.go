package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"erp-service/masterdata/masterdatadto"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

const (
	masterdataCategoryByIDPrefix     = "md:cat:id:%s"
	masterdataCategoryByCodePrefix   = "md:cat:code:%s"
	masterdataCategoriesListPrefix   = "md:cats:list:%s"
	masterdataCategoryChildrenPrefix = "md:cat:children:%s"
	masterdataItemByIDPrefix         = "md:item:id:%s"
	masterdataItemByCodePrefix       = "md:item:code:%s:%s:%s"
	masterdataItemsListPrefix        = "md:items:list:%s"
	masterdataItemChildrenPrefix     = "md:item:children:%s"
	masterdataItemTreePrefix         = "md:item:tree:%s:%s"
	masterdataItemDefaultPrefix      = "md:item:default:%s:%s"

	masterdataCategoryPattern = "md:cat*"
	masterdataItemPattern     = "md:item*"
)

func (r *Redis) masterdataCategoryByIDKey(id uuid.UUID) string {
	return fmt.Sprintf(masterdataCategoryByIDPrefix, id.String())
}

func (r *Redis) masterdataCategoryByCodeKey(code string) string {
	return fmt.Sprintf(masterdataCategoryByCodePrefix, code)
}

func (r *Redis) masterdataCategoriesListKey(filterHash string) string {
	return fmt.Sprintf(masterdataCategoriesListPrefix, filterHash)
}

func (r *Redis) masterdataCategoryChildrenKey(parentID uuid.UUID) string {
	return fmt.Sprintf(masterdataCategoryChildrenPrefix, parentID.String())
}

func (r *Redis) masterdataItemByIDKey(id uuid.UUID) string {
	return fmt.Sprintf(masterdataItemByIDPrefix, id.String())
}

func (r *Redis) masterdataItemByCodeKey(categoryID uuid.UUID, tenantID *uuid.UUID, code string) string {
	tenantStr := "global"
	if tenantID != nil {
		tenantStr = tenantID.String()
	}
	return fmt.Sprintf(masterdataItemByCodePrefix, categoryID.String(), tenantStr, code)
}

func (r *Redis) masterdataItemsListKey(filterHash string) string {
	return fmt.Sprintf(masterdataItemsListPrefix, filterHash)
}

func (r *Redis) masterdataItemChildrenKey(parentID uuid.UUID) string {
	return fmt.Sprintf(masterdataItemChildrenPrefix, parentID.String())
}

func (r *Redis) masterdataItemTreeKey(categoryCode string, tenantID *uuid.UUID) string {
	tenantStr := "global"
	if tenantID != nil {
		tenantStr = tenantID.String()
	}
	return fmt.Sprintf(masterdataItemTreePrefix, categoryCode, tenantStr)
}

func (r *Redis) masterdataItemDefaultKey(categoryID uuid.UUID, tenantID *uuid.UUID) string {
	tenantStr := "global"
	if tenantID != nil {
		tenantStr = tenantID.String()
	}
	return fmt.Sprintf(masterdataItemDefaultPrefix, categoryID.String(), tenantStr)
}

func (r *Redis) GetCategoryByID(ctx context.Context, id uuid.UUID) (*masterdatadto.CategoryResponse, error) {
	key := r.masterdataCategoryByIDKey(id)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var category masterdatadto.CategoryResponse
	if err := json.Unmarshal(data, &category); err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *Redis) SetCategoryByID(ctx context.Context, id uuid.UUID, category *masterdatadto.CategoryResponse, ttl time.Duration) error {
	key := r.masterdataCategoryByIDKey(id)
	data, err := json.Marshal(category)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *Redis) GetCategoryByCode(ctx context.Context, code string) (*masterdatadto.CategoryResponse, error) {
	key := r.masterdataCategoryByCodeKey(code)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var category masterdatadto.CategoryResponse
	if err := json.Unmarshal(data, &category); err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *Redis) SetCategoryByCode(ctx context.Context, code string, category *masterdatadto.CategoryResponse, ttl time.Duration) error {
	key := r.masterdataCategoryByCodeKey(code)
	data, err := json.Marshal(category)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *Redis) GetCategoriesList(ctx context.Context, filterHash string) (*masterdatadto.ListCategoriesResponse, error) {
	key := r.masterdataCategoriesListKey(filterHash)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var response masterdatadto.ListCategoriesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *Redis) SetCategoriesList(ctx context.Context, filterHash string, response *masterdatadto.ListCategoriesResponse, ttl time.Duration) error {
	key := r.masterdataCategoriesListKey(filterHash)
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *Redis) GetCategoryChildren(ctx context.Context, parentID uuid.UUID) ([]*masterdatadto.CategoryResponse, error) {
	key := r.masterdataCategoryChildrenKey(parentID)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var categories []*masterdatadto.CategoryResponse
	if err := json.Unmarshal(data, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *Redis) SetCategoryChildren(ctx context.Context, parentID uuid.UUID, categories []*masterdatadto.CategoryResponse, ttl time.Duration) error {
	key := r.masterdataCategoryChildrenKey(parentID)
	data, err := json.Marshal(categories)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *Redis) InvalidateCategories(ctx context.Context) error {
	return r.DeleteByPattern(ctx, masterdataCategoryPattern)
}

func (r *Redis) GetItemByID(ctx context.Context, id uuid.UUID) (*masterdatadto.ItemResponse, error) {
	key := r.masterdataItemByIDKey(id)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var item masterdatadto.ItemResponse
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Redis) SetItemByID(ctx context.Context, id uuid.UUID, item *masterdatadto.ItemResponse, ttl time.Duration) error {
	key := r.masterdataItemByIDKey(id)
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *Redis) GetItemByCode(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, code string) (*masterdatadto.ItemResponse, error) {
	key := r.masterdataItemByCodeKey(categoryID, tenantID, code)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var item masterdatadto.ItemResponse
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Redis) SetItemByCode(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, code string, item *masterdatadto.ItemResponse, ttl time.Duration) error {
	key := r.masterdataItemByCodeKey(categoryID, tenantID, code)
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *Redis) GetItemsList(ctx context.Context, filterHash string) (*masterdatadto.ListItemsResponse, error) {
	key := r.masterdataItemsListKey(filterHash)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var response masterdatadto.ListItemsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *Redis) SetItemsList(ctx context.Context, filterHash string, response *masterdatadto.ListItemsResponse, ttl time.Duration) error {
	key := r.masterdataItemsListKey(filterHash)
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *Redis) GetItemChildren(ctx context.Context, parentID uuid.UUID) ([]*masterdatadto.ItemResponse, error) {
	key := r.masterdataItemChildrenKey(parentID)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var items []*masterdatadto.ItemResponse
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Redis) SetItemChildren(ctx context.Context, parentID uuid.UUID, items []*masterdatadto.ItemResponse, ttl time.Duration) error {
	key := r.masterdataItemChildrenKey(parentID)
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *Redis) GetItemTree(ctx context.Context, categoryCode string, tenantID *uuid.UUID) ([]*masterdatadto.ItemResponse, error) {
	key := r.masterdataItemTreeKey(categoryCode, tenantID)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var items []*masterdatadto.ItemResponse
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Redis) SetItemTree(ctx context.Context, categoryCode string, tenantID *uuid.UUID, items []*masterdatadto.ItemResponse, ttl time.Duration) error {
	key := r.masterdataItemTreeKey(categoryCode, tenantID)
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *Redis) GetItemDefault(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID) (*masterdatadto.ItemResponse, error) {
	key := r.masterdataItemDefaultKey(categoryID, tenantID)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var item masterdatadto.ItemResponse
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Redis) SetItemDefault(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, item *masterdatadto.ItemResponse, ttl time.Duration) error {
	key := r.masterdataItemDefaultKey(categoryID, tenantID)
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *Redis) InvalidateItems(ctx context.Context) error {
	return r.DeleteByPattern(ctx, masterdataItemPattern)
}
