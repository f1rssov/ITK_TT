package repository

import(
	"context"
	"github.com/google/uuid"
	"os"
	"testing"
	"walletT/internal/model"
	"walletT/internal/storage"
	"walletT/internal/testutils"
	"github.com/stretchr/testify/assert"
)

var repo *WalletRepo
var testID uuid.UUID

func TestMain(m *testing.M) {
	ctx := context.Background()

	adminDSN := os.Getenv("DSN")
	testDBName := os.Getenv("DB_TEST_NAME")
	migrationsPath := "../../migrations"

	testDSN, err := testutils.CreateTestDatabase(ctx, adminDSN, testDBName)
	if err != nil {
		panic(err)
	}

	err = testutils.ApplyMigrations(migrationsPath, testDSN)
	if err != nil {
		panic(err)
	}

	pool, err := storage.NewPostgres(ctx, testDSN)
	if err != nil {
		panic(err)
	}

	repo = NewWalletRepository(pool)

	os.Exit(m.Run())
}

func TestCreateWallet(t *testing.T) {
	testID = uuid.New()
	wallet := &model.Wallet{
		Id:     testID,
		Amount: 100,
	}

	err := repo.CreateWallet(context.Background(), wallet)
	assert.NoError(t, err, "ошибка при создании кошелька")
}

func TestGetWalletAmount(t *testing.T) {
	amount, err := repo.GetWalletAmount(context.Background(), testID)
	assert.NoError(t, err)
	assert.Equal(t, 100, amount, "неправильный баланс")
}

func TestUpdateWalletAmount_Deposit(t *testing.T) {
	err := repo.UpdateWalletAmount(context.Background(), testID, "DEPOSIT", 50)
	assert.NoError(t, err)

	amount, err := repo.GetWalletAmount(context.Background(), testID)
	assert.NoError(t, err)
	assert.Equal(t, 150, amount)
}

func TestUpdateWalletAmount_Withdraw(t *testing.T) {
	err := repo.UpdateWalletAmount(context.Background(), testID, "WITHDRAW", 30)
	assert.NoError(t, err)

	amount, err := repo.GetWalletAmount(context.Background(), testID)
	assert.NoError(t, err)
	assert.Equal(t, 120, amount)
}

func TestUpdateWalletAmount_Overdraft(t *testing.T) {
	err := repo.UpdateWalletAmount(context.Background(), testID, "WITHDRAW", 9999)
	assert.Error(t, err)
	assert.EqualError(t, err, "недостаточно средств для снятия")
}