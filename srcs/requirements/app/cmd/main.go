package main
import(
	"context"
	"os"
	"walletT/internal/storage"
	"walletT/internal/repository"
	"walletT/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/swaggo/files"
    "github.com/swaggo/gin-swagger"
    "walletT/docs"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}
	
	parsLevel, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		logrus.Warnf("Некорректный LOG_LEVEL='%s', используется Info по умолчанию", level)
		parsLevel = logrus.InfoLevel
	}

	logrus.SetLevel(parsLevel)
	logrus.Infof("Уровень логирования установлен: %s", level)
}


func  main(){
	ctx := context.Background()
	dsn  := os.Getenv("DSN")

	logrus.Info("Попытка подключения к базе данных...")

	db, err := storage.NewPostgres(ctx, dsn)
	if err != nil{
		logrus.Fatalf("Ошибка подключения к базе данных: %v", err)
		return
	}
	defer func() {
        logrus.Info("Закрытие подключения к базе данных")
        db.Close()
    }()
	logrus.Info("Подключение к базе данных успешно")

	walletRepo := repository.NewWalletRepository(db)
	handlerRepo := handler.NewHandlerRepository(walletRepo)

	router := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("api/v1")
	v1.POST("/wallets", handlerRepo.UpdateWalletAmount)                    
	v1.GET("/wallets/:wallet_id", handlerRepo.GetWalletAmount) 
	v1.PUT("/wallets", handlerRepo.CreateWalletTempl) //для ручного теста хендлеров


	logrus.Info("Запуск сервера на порту :8080")
	if err := router.Run(":8080"); err != nil {
		logrus.Fatalf("Ошибка при запуске сервера: %v", err)
	}

}