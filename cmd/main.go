package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/gorm"
	"log"
	"net/http"
	"tf_ocg/cmd/app/router"
	"tf_ocg/pkg/database_manager"
)

type Server struct {
	Db     *gorm.DB
	Router *mux.Router
}

func Init() {
	var server = Server{}
	server.Db = database_manager.InitDb()
	//server.Db.AutoMigrate(&models.User{})
	server.Router = mux.NewRouter()
	router.InitializeRoutes(server.Router)
	server.Run(":8000")
}
func (server *Server) Run(addr string) {
	fmt.Println("Listening to port " + addr)
	log.Fatal(http.ListenAndServe(addr, cors.AllowAll().Handler(server.Router)))
}
func main() {
	// init server
	Init()
}
