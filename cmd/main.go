package main

import (
	"os"
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
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initial config: %s", err.Error())
	}

	if err:= godotenv.Load(); err != nil {
		logrus.Fatalf("error initial config: %s", err.Error())
	}

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

	repos:=repository.NewRepository(db)
	services:=service.NewService(repos)
	handlers:=handler.NewHandler(services)
	
	
	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"),handlers.InitRoutes()); err != nil {
		logrus.Fatalf("failed to run server: %v", err)
	}
}

func initConfig() error{
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}