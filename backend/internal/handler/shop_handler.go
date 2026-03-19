package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SX110903/match_app/backend/internal/auth"
	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/service"
	"github.com/SX110903/match_app/backend/pkg/response"
)

var _ = strconv.Atoi // ensure strconv used

type ShopHandler struct {
	shopSvc service.IShopService
}

func NewShopHandler(shopSvc service.IShopService) *ShopHandler {
	return &ShopHandler{shopSvc: shopSvc}
}

func (h *ShopHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.shopSvc.GetItems(r.Context())
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, items)
}

func (h *ShopHandler) Purchase(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	var req service.PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if err := h.shopSvc.Purchase(r.Context(), claims.Subject, req.ItemType, req.ItemValue); err != nil {
		switch err {
		case domain.ErrInvalidInput:
			response.BadRequest(w, "invalid purchase: check item_type, item_value, and credits")
		default:
			response.InternalError(w)
		}
		return
	}
	response.OK(w, map[string]string{"status": "purchased"})
}

func (h *ShopHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	txs, err := h.shopSvc.GetTransactions(r.Context(), claims.Subject, page, 20)
	if err != nil {
		response.InternalError(w)
		return
	}
	if txs == nil {
		txs = []domain.ShopTransaction{}
	}
	response.OK(w, txs)
}
