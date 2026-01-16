package transport

import (
	"fmt"
	"net/http"
	"strconv"
	"user-service/internal/models"
	"user-service/internal/services"
	"user-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	users  services.UserService
	wallet services.WalletService
}

func NewWalletHandler(
	users services.UserService, wallet services.WalletService,
) *WalletHandler {
	return &WalletHandler{users: users, wallet: wallet}
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	uid := c.GetUint("user_id")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	wallet, err := h.wallet.GetWallet(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

func (h *WalletHandler) WalletDeposit(c *gin.Context) {
	var req models.TransactionForRequest

	uid := c.GetUint("user_id")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive!"})
		return
	}

	if req.Description == "" {
		req.Description = utils.DefaultDescription
	}

	wallet, transaction, err := h.wallet.Deposit(uid, req.Amount, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wallet":         wallet,
		"transaction_id": transaction.ID},
	)
}

func (h *WalletHandler) WalletFreeze(c *gin.Context) {
	var req models.TransactionForRequest

	uid := c.GetUint("user_id")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive!"})
		return
	}

	if req.Description == "" {
		req.Description = utils.DefaultDescription
	}

	wallet, transaction, err := h.wallet.Freeze(uid, req.Amount, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wallet":         wallet,
		"transaction_id": transaction.ID},
	)
}

func (h *WalletHandler) WalletUnfreeze(c *gin.Context) {
	var req models.TransactionForRequest

	uid := c.GetUint("user_id")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive!"})
		return
	}

	if req.Description == "" {
		req.Description = utils.DefaultDescription
	}

	wallet, transaction, err := h.wallet.Unfreeze(uid, req.Amount, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wallet":         wallet,
		"transaction_id": transaction.ID},
	)
}

func (h *WalletHandler) WalletCharge(c *gin.Context) {
	var req models.TransactionForRequest

	uid := c.GetUint("user_id")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive!"})
		return
	}

	if req.Description == "" {
		req.Description = utils.DefaultDescription
	}

	wallet, transaction, err := h.wallet.Charge(uid, req.Amount, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wallet":         wallet,
		"transaction_id": transaction.ID},
	)
}

func (h *WalletHandler) ListTransactions(c *gin.Context) {
	uid := c.GetUint("user_id")
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit, err := parseQueryInt(c, "limit", 10, 1, 100)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offset, err := parseQueryInt(c, "offset", 0, 0, 10000)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactions, err := h.wallet.ListTransactions(uid, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

func parseQueryInt(c *gin.Context, key string, defaultVal, min, max int) (int, error) {
	val := c.Query(key)
	if val == "" {
		return defaultVal, nil
	}
	parsed, err := strconv.Atoi(val)
	if err != nil || parsed < min || parsed > max {
		return 0, fmt.Errorf("invalid %s", key)
	}
	return parsed, nil
}
