package service

import (
	"context"
	"fmt"

	"github.com/SX110903/match_app/backend/internal/config"
	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/repository"
)

type shopService struct {
	shopRepo repository.IShopRepository
}

func NewShopService(shopRepo repository.IShopRepository) IShopService {
	return &shopService{shopRepo: shopRepo}
}

func (s *shopService) GetItems(ctx context.Context) ([]config.ShopItem, error) {
	return config.VIPItems, nil
}

func (s *shopService) Purchase(ctx context.Context, userID string, itemType string, itemValue int) error {
	// Accept "vip" as alias for "vip_upgrade"
	if itemType == "vip" {
		itemType = "vip_upgrade"
	}
	if itemType != "vip_upgrade" {
		return domain.ErrInvalidInput
	}
	if itemValue < 1 || itemValue > 5 {
		return domain.ErrInvalidInput
	}

	var cost int
	for _, item := range config.VIPItems {
		if item.ItemType == itemType && item.ItemValue == itemValue {
			cost = item.Cost
			break
		}
	}
	if cost == 0 {
		return domain.ErrInvalidInput
	}

	if err := s.shopRepo.PurchaseVIP(ctx, userID, itemValue, cost); err != nil {
		return fmt.Errorf("purchase: %w", err)
	}
	return nil
}

func (s *shopService) GetTransactions(ctx context.Context, userID string, page, limit int) ([]domain.ShopTransaction, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	return s.shopRepo.GetTransactionsByUser(ctx, userID, limit, offset)
}
