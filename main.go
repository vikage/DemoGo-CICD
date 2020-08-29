package main

import (
	"go-cicd/app/authenticate"
	"go-cicd/app/database"
	"go-cicd/app/di"
	"go-cicd/app/domain/repository"
	"go-cicd/app/logger"
	"go-cicd/app/restful"

	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func waitSystemEndSignal() {
	EndChannel := make(chan os.Signal)
	signal.Notify(EndChannel, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	select {
	case output := <-EndChannel:
		log.Printf("end http server process by: %s", output)
	}

	close(EndChannel)
}

func loadEnvironment() {
	err := godotenv.Load()
	if err != nil {
		panic("Can not load .env")
	}
}

func configLogger() {
	logger.AddLogEngine(logger.ConsoleLog{})
	logger.AddLogEngine(logger.FileLog{})
}

func configDI() {
	database.RegisterDependencyInContainer(di.DefaultContainer)
	authenticate.RegisterDependencyInContainer(di.DefaultContainer)
	repository.RegisterDependencyInContainer(di.DefaultContainer)
}

func main() {
	loadEnvironment()
	configLogger()
	configDI()

	database.EnsureIndexes()

	restful.StartWebAPI()
	waitSystemEndSignal()
}
