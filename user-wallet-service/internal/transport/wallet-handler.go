package transport

import (
	"fmt"
	"log/slog"
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
	logger *slog.Logger
}

func NewWalletHandler(
	users services.UserService, wallet services.WalletService, logger *slog.Logger,
) *WalletHandler {
	return &WalletHandler{users: users, wallet: wallet, logger: logger}
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	uid := c.GetUint("user_id")
	if uid == 0 {
		h.logger.Warn("get wallet unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	h.logger.Info("get wallet", "user_id", uid)

	wallet, err := h.wallet.GetWallet(uid)
	if err != nil {
		h.logger.Error("get wallet failed", "user_id", uid, "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wallet)
}

func (h *WalletHandler) WalletDeposit(c *gin.Context) {
	var req models.TransactionForRequest

	uid := c.GetUint("user_id")
	if uid == 0 {
		h.logger.Warn("deposit unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("deposit bad request", "user_id", uid, "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Description == "" {
		req.Description = utils.DefaultDescription
	}

	h.logger.Info("deposit attempt", "user_id", uid, "amount", req.Amount)

	wallet, transaction, err := h.wallet.Deposit(uid, req.Amount, req.Description)
	if err != nil {
		h.logger.Error("deposit failed", "user_id", uid, "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("deposit success", "user_id", uid, "transaction_id", transaction.ID)

	c.JSON(http.StatusOK, gin.H{
		"wallet":         wallet,
		"transaction_id": transaction.ID},
	)
}

func (h *WalletHandler) WalletFreeze(c *gin.Context) {
	var req models.TransactionForRequest

	uid := c.GetUint("user_id")
	if uid == 0 {
		h.logger.Warn("freeze unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("freeze bad request", "user_id", uid, "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Description == "" {
		req.Description = utils.DefaultDescription
	}

	h.logger.Info("freeze attempt", "user_id", uid, "amount", req.Amount)

	wallet, transaction, err := h.wallet.Freeze(uid, req.Amount, req.Description)
	if err != nil {
		h.logger.Error("freeze failed", "user_id", uid, "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("freeze success", "user_id", uid, "transaction_id", transaction.ID)

	c.JSON(http.StatusOK, gin.H{
		"wallet":         wallet,
		"transaction_id": transaction.ID},
	)
}

func (h *WalletHandler) WalletUnfreeze(c *gin.Context) {
	var req models.TransactionForRequest

	uid := c.GetUint("user_id")
	if uid == 0 {
		h.logger.Warn("unfreeze unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("unfreeze bad request", "user_id", uid, "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Description == "" {
		req.Description = utils.DefaultDescription
	}

	h.logger.Info("unfreeze attempt", "user_id", uid, "amount", req.Amount)

	wallet, transaction, err := h.wallet.Unfreeze(uid, req.Amount, req.Description)
	if err != nil {
		h.logger.Error("unfreeze failed", "user_id", uid, "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("unfreeze success", "user_id", uid, "transaction_id", transaction.ID)

	c.JSON(http.StatusOK, gin.H{
		"wallet":         wallet,
		"transaction_id": transaction.ID},
	)
}

func (h *WalletHandler) WalletCharge(c *gin.Context) {
	var req models.TransactionForRequest

	uid := c.GetUint("user_id")
	if uid == 0 {
		h.logger.Warn("charge unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("charge bad request", "user_id", uid, "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Description == "" {
		req.Description = utils.DefaultDescription
	}

	h.logger.Info("charge attempt", "user_id", uid, "amount", req.Amount)

	wallet, transaction, err := h.wallet.Charge(uid, req.Amount, req.Description)
	if err != nil {
		h.logger.Error("charge failed", "user_id", uid, "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("charge success", "user_id", uid, "transaction_id", transaction.ID)

	c.JSON(http.StatusOK, gin.H{
		"wallet":         wallet,
		"transaction_id": transaction.ID},
	)
}

func (h *WalletHandler) ListTransactions(c *gin.Context) {
	uid := c.GetUint("user_id")
	if uid == 0 {
		h.logger.Warn("list transactions unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit, err := parseQueryInt(c, "limit", 10, 1, 100)
	if err != nil {
		h.logger.Warn("list transactions bad request", "user_id", uid, "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offset, err := parseQueryInt(c, "offset", 0, 0, 10000)
	if err != nil {
		h.logger.Warn("list transactions bad request", "user_id", uid, "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("list transactions", "user_id", uid, "limit", limit, "offset", offset)

	transactions, err := h.wallet.ListTransactions(uid, limit, offset)
	if err != nil {
		h.logger.Error("list transactions failed", "user_id", uid, "err", err.Error())
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
