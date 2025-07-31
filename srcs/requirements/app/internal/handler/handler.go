package handler

import (
	"walletT/internal/repository"
	"walletT/internal/model"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/sirupsen/logrus"
)

type WalletUpdateInput struct {
    Id      uuid.UUID `json:"valletId" binding:"required"`
    Operation     string    `json:"operationType" binding:"required"`
    Amount        int       `json:"amount" binding:"required"`
}

type HandlerRepo struct {
	repo *repository.WalletRepo
}

func NewHandlerRepository(repo *repository.WalletRepo) *HandlerRepo{
	return &HandlerRepo{repo: repo}
}

// нужна только для теста хендлеров
func (h *HandlerRepo)CreateWalletTempl(c *gin.Context){
	templateID := uuid.New()
	wallet := &model.Wallet{
		Id:     templateID,
		Amount: 100,
	}
	err := h.repo.CreateWallet(c.Request.Context(), wallet)
	if err!=nil{
		logrus.WithError(err).Error("handler: ошибка создания кошелька")
	}
}

// UpdateWalletAmount
// @Summary Обновить сумму на кошельке
// @Description Выполняет операцию (пополнение или списание) на балансе кошелька
// @Tags Wallet
// @Accept json
// @Produce json
// @Param input body WalletUpdateInput true "Данные для обновления баланса"
// @Success 200 {string} string "Операция успешно выполнена"
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Failure 404 {object} map[string]string "Кошелек не найден"
// @Router /wallet [put]
func (h *HandlerRepo)UpdateWalletAmount(c *gin.Context){
	var input WalletUpdateInput
	if err := c.BindJSON(&input); err !=nil{
		logrus.WithError(err).Error("handler: Неверный формат запроса")
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
		return
	}
	logrus.Debug("handler: Обновление баланса") // debug
	if err := h.repo.UpdateWalletAmount(c.Request.Context(), input.Id, input.Operation, input.Amount); err!=nil{
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Операция успешно  выполнена")
}

// GetWalletAmount
// @Summary Получить сумму на кошельке
// @Description Возвращает текущий баланс кошелька по ID
// @Tags Wallet
// @Accept json
// @Produce json
// @Param wallet_id path string true "ID кошелька в формате UUID"
// @Success 200 {object} model.Wallet "Информация о кошельке с балансом"
// @Failure 400 {object} map[string]string "Неверный формат wallet_id"
// @Failure 404 {object} map[string]string "Кошелек не найден"
// @Router /wallet/{wallet_id} [get]
func (h *HandlerRepo)GetWalletAmount(c *gin.Context){
	wallet_id := c.Param("wallet_id")

	logrus.WithFields(logrus.Fields{		// debug
		"id":        wallet_id,				//
	}).Debug("handler: Парсинг wallet_id") 	//

	w_id, err := uuid.Parse(wallet_id)
	if err != nil {
		logrus.WithError(err).Error("handler: Неверный формат wallet_id в URL")
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат wallet_id"})
		return
	}

	logrus.WithFields(logrus.Fields{		// debug
		"id":        w_id,					//
	}).Debug("handler: Получение баланса") 	//

	wal, err := h.repo.GetWalletAmount(c.Request.Context(), w_id)
	if err != nil{
		logrus.WithError(err).Error("handler: Ошибка получения баланса")
		c.JSON(http.StatusNotFound, gin.H{"error": "Кошелька с таким  id не существует "})
		return
	}
	c.JSON(http.StatusOK, wal)
}