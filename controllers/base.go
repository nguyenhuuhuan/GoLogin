package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/http"
	"web-golang/middlewares"
	"web-golang/models"
)

import "github.com/jinzhu/gorm"

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error

	if DbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		server.DB, err = gorm.Open(DbDriver, DBURL)

		if err != nil {
			fmt.Printf("Cannot connect to %s database", DbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connect to the %s database", DbDriver)
		}
	}

	server.DB.Debug().AutoMigrate(&models.User{})
	server.DB.Debug().AutoMigrate(&models.Roles{})
	server.Router = mux.NewRouter()
	server.initializeRoutes()
}

func (server *Server) initializeRoutes() {

	//LoginUser
	server.Router.HandleFunc("/auth/login", middlewares.SetMiddlewareJSON(server.Login)).Methods("POST")

	//User
	server.Router.HandleFunc("/auth/register", middlewares.SetMiddlewareJSON(server.Register)).Methods("POST")
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8000")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
