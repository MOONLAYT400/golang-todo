package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"todo"
	"todo/pkg/handler"
	"todo/pkg/repository"
	"todo/pkg/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	// настройка кастомных логов 
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initial config: %s", err.Error())
	}


	// загрузка переменных окружения
	if err:= godotenv.Load(); err != nil {
		logrus.Fatalf("error initial config: %s", err.Error())
	}

	// создание подключения к базе данных
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %v", err)
	}

	// создание репозиториев для взаимодействия с базой данных
	repos:=repository.NewRepository(db)

	// создание сервисов для взаимодействия с репозиториями
	services:=service.NewService(repos)

	// создание хендлеров для взаимодействия с хендлерами
	handlers:=handler.NewHandler(services)
	
	
	srv := new(todo.Server)

	// го рутина, в которой запускаеться сервер, бесконечный цикл, слцушающий запросы 
	go func() {
		// запуск сервера и ини циализация роутов
	if err := srv.Run(viper.GetString("port"),handlers.InitRoutes()); err != nil {
		logrus.Fatalf("failed to run server: %v", err)
	}
	}()

	logrus.Print("App started")

	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("App shutdown")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("failed to shutdown server: %v", err)
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("failed to close database: %v", err)
	}
	}




func initConfig() error{
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}