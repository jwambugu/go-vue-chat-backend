package main

import (
	"chatapp/pkg/database"
	"chatapp/pkg/util"
	"chatapp/repository/mysql"
	"chatapp/services/chatroom"
	"chatapp/services/user"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"os/signal"
)

var (
	config *util.Config
)

// application provides dependency injection across the system
type application struct {
	config          *util.Config
	db              *sqlx.DB
	userService     user.Service
	chatroomService chatroom.Service
}

func init() {
	var err error

	config, err = util.ReadConfig(util.GetAbsolutePath())
	if err != nil {
		log.Fatal(err)
	}

}

func (app *application) initServices() {
	db, err := database.NewMySQLConnection(app.config.DBConfig.MySQL.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	app.db = db
	app.userService = user.NewService(mysql.NewUserRepository(app.db))
	app.chatroomService = chatroom.NewService(mysql.NewChatRoomRepository(app.db))
}

func main() {
	app := &application{
		config: config,
	}

	// Start up the services
	app.initServices()
	fiberApp := app.routes()

	osSigChan := make(chan os.Signal, 1)
	defer close(osSigChan)

	signal.Notify(osSigChan, os.Interrupt, os.Kill)

	go func() {
		_ = <-osSigChan

		fmt.Println("Gracefully shutting down the server...")
		if err := fiberApp.Shutdown(); err != nil {
			log.Fatalf("unexpected error shutting down the server:: %v", err)
		}
	}()

	addr := fmt.Sprintf(":%d", app.config.AppPort)
	if err := fiberApp.Listen(addr); err != nil {
		log.Fatalf("error starting server on port %s :: %v", addr, err)
	}
}
