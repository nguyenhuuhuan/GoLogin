package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/http"
	"os"
	"time"
	"web-golang/logger"
	"web-golang/middlewares"
	"web-golang/models"
)

import "github.com/jinzhu/gorm"

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func openFile(pathFile string) *os.File {
	f, err := os.OpenFile(pathFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return f
}
func (server *Server) Initialize(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	_ = logger.New(
		log.New(openFile("/home/nhhuan/Desktop/info.log"), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	if DbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		server.DB, err = gorm.Open(DbDriver, DBURL)
		if err != nil {
			//newLogger.Error(err.Error())
			fmt.Printf("Cannot connect to %s database", DbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connect to the %s database", DbDriver)
		}
	}

	server.DB.Debug().AutoMigrate(
		&models.Roles{},
		&models.User{},
		&models.Beverage{},
		&models.Topping{},
		&models.OrderDetail{},
		&models.Order{},
		&models.ToppingDetail{},
	)
	server.Router = mux.NewRouter()
	server.initializeRoutes()
}

func (server *Server) initializeRoutes() {

	//LoginUser
	server.Router.HandleFunc("/auth/login", middlewares.SetMiddlewareJSON(server.Login)).Methods("POST")
	server.Router.HandleFunc("/auth/logout", middlewares.SetMiddlewareJSON(server.LogoutHandler)).Methods("GET")

	//User
	server.Router.HandleFunc("/user/search", middlewares.SetMiddlewareJSON(server.SearchByUsername)).Methods("GET")

	server.Router.HandleFunc("/user/", middlewares.SetMiddlewareJSON(server.GetAllUsers)).Methods("GET")
	server.Router.HandleFunc("/auth/register", middlewares.SetMiddlewareJSON(server.Register)).Methods("POST")
	server.Router.HandleFunc("/user/{id}", middlewares.SetMiddlewareJSON(server.GetUser)).Methods("GET")
	server.Router.HandleFunc("/user/{id}", middlewares.SetMiddlewareAuthentication(server.UpdateUser)).Methods("PUT")
	//Role
	server.Router.HandleFunc("/role", middlewares.SetMiddlewareJSON(server.CreateRole)).Methods("POST")

	//Beverage
	server.Router.HandleFunc("/beverage", middlewares.SetMiddlewareJSON(server.createBeverage)).Methods("POST")
	server.Router.HandleFunc("/beverage", middlewares.SetMiddlewareJSON(server.GetAllBeverage)).Methods("GET")
	server.Router.HandleFunc("/beverage/type", middlewares.SetMiddlewareJSON(server.GetBeveragesByType)).Methods("GET")

	//Cart
	server.Router.HandleFunc("/cart", middlewares.SetMiddlewareJSON(server.GetAllCart)).Methods("GET")
	server.Router.HandleFunc("/cart/addItem/{id}", middlewares.SetMiddlewareJSON(server.addBeverageToCart)).Methods("GET")
	server.Router.HandleFunc("/cart/removeItem/{id}", middlewares.SetMiddlewareJSON(server.RemoveItem)).Methods("GET")

	//Topping
	server.Router.HandleFunc("/topping", middlewares.SetMiddlewareJSON(server.CreateTopping)).Methods("POST")

	//OrderDetail
	server.Router.HandleFunc("/orderDetail", middlewares.SetMiddlewareJSON(server.SaveCart)).Methods("GET")
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8000")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
