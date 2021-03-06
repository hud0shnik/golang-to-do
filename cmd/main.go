package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"todo-app"
	"todo-app/pkg/handler"
	"todo-app/pkg/repository"
	"todo-app/pkg/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// Функция инициализации сонфигурационного файла
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {

	// Задаёт для логгера формат JSON
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// Инициализация конфига
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	// Инициализация базы данных
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	// Проверка инициализации БД
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	// Слой работы с БД
	repos := repository.NewRepository(db)
	// Слой бизнес логики
	services := service.NewService(repos)
	// Слой работы с http
	handlers := handler.NewHandler(services)

	// Создание сервера
	srv := new(todo.Server)

	// Запуск сервера в горутине
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()

	// Вывод информации о запуске приложения
	logrus.Print("App started")

	// Канал для получения сигнала о закрытии приложения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// Чтение из канала, которое блокирует выполнение мейна
	<-quit

	// Вывод информации о закрытии приложения
	logrus.Print("App shutting down")

	// Остановка сервера
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occurred on server shutting down")
	}

	// Закрытие всех соединений с БД
	if err := db.Close(); err != nil {
		logrus.Errorf("error occurred on db connection close")
	}

}
