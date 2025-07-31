package repository
import(
	"context"
	"github.com/google/uuid"
	"walletT/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"fmt"

)

type WalletRepo struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) *WalletRepo{
	return &WalletRepo{db: db}
}
// 5d5cce3c-1db4-43fb-946b-30375b7f5341 - uuid для теста
func (h *WalletRepo)CreateWallet(ctx context.Context, mod *model.Wallet) error{
	logrus.WithFields(logrus.Fields{	// debug
		"id":     mod.Id,				//
		"amount": mod.Amount,			//
	}).Debug("Создание кошелька")		//
	logrus.Info("Создание кошелька")

	query := `
        INSERT INTO wallets (id, balance)
        VALUES ($1, $2)
    `

	_, err := h.db.Exec(ctx, query, mod.Id, mod.Amount)
	if err != nil {
		logrus.WithError(err).WithField("id", mod.Id).Error("Ошибка при создании кошелька")
	}
	return err
}

func (h *WalletRepo)GetWalletAmount(ctx context.Context, id uuid.UUID) (int, error){
	var balance int
	logrus.WithFields(logrus.Fields{	// debug
		"id":     id,					//
	}).Debug("Получение баланса")		//
	logrus.Info("Получение баланса")

	query := `SELECT balance FROM wallets WHERE id = $1`

	err := h.db.QueryRow(ctx, query, id).Scan(&balance)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("Ошибка при получении баланса")
		return 0, err
	}
	return balance, nil
}

func (h *WalletRepo) UpdateWalletAmount(ctx context.Context, id uuid.UUID, oper string, amount int) error {
	logrus.Info("Начало транзакции")
	tx, err := h.db.Begin(ctx)
	if err != nil {
		logrus.WithError(err).Debug("Не удалось начать транзакцию") // debug
		return err
	}
	defer tx.Rollback(ctx)

	logrus.WithField("id", id).Debug("Блокировка строки кошелька (FOR UPDATE)")
	var currentBalance int
	err = tx.QueryRow(ctx, `SELECT balance FROM wallets WHERE id = $1 FOR UPDATE`, id).Scan(&currentBalance)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Debug("Ошибка при блокировке строки кошелька") // debug
		return fmt.Errorf("кошелек с таким id не найден")
	}
	
	switch oper {
	case "DEPOSIT":
		currentBalance += amount
		logrus.WithFields(logrus.Fields{	// debug
			"id":        id,				//
			"operation": oper,				//	
			"amount":    amount,			//
		}).Debug("Пополнение кошелька") 	//
		logrus.Info("Пополнение кошелька")
	case "WITHDRAW":
		if currentBalance < amount {
			logrus.WithFields(logrus.Fields{			// debug
				"id":             id,					//
				"operation":      oper,					//
				"requested":      amount,				//
				"currentBalance": currentBalance,		//
			}).Debug("Недостаточно средств для снятия")	//
			return fmt.Errorf("недостаточно средств для снятия")
		}
		currentBalance -= amount
		logrus.WithFields(logrus.Fields{	// debug
			"id":        id,				//
			"operation": oper,				//
			"amount":    amount,			//
		}).Debug("Снятие с кошелька")		//
	default:
		logrus.WithField("operation", oper).Debug("Неизвестный тип операции") // debug
		return fmt.Errorf("неизвестный тип операции")
	}

	_, err = tx.Exec(ctx, `UPDATE wallets SET balance = $1 WHERE id = $2`, currentBalance, id)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("Ошибка при обновлении баланса")
		return err
	}

	err = tx.Commit(ctx)
	if err != nil{
		logrus.WithError(err).Error("Ошибка при коммите транзакции")
		return  err
	}
	logrus.WithFields(logrus.Fields{	// debug
		"id":      id,					//
		"balance": currentBalance,		//
	}).Debug("Баланс успешно обновлён") //
	return nil
}