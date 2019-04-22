package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"uapp.com/controllers"
	"uapp.com/models"

	"github.com/gorilla/mux"
)

var (
	host     = os.Getenv("PGHOST")
	port1    = os.Getenv("PGPORT")
	port, _  = strconv.Atoi(os.Getenv("PGPORT"))
	user     = os.Getenv("PGUSER")
	password = os.Getenv("PGPASSWORD")
	dbname   = os.Getenv("PGDB")
)

func main() {
	//  Create DB connection string and then use it to create
	//  our model services.
	fmt.Println("host: ", host)
	fmt.Println("port: ", port)
	fmt.Println("port1: ", port1)
	fmt.Println("user: ", user)
	fmt.Println("password: ", password)
	fmt.Println("dbname: ", dbname)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.AutoMigrate()

	usersC := controllers.NewUsers(us)

	r := mux.NewRouter()
	r.HandleFunc("/", usersC.List).Methods("Get")
	r.HandleFunc("/signup", usersC.New).Methods("Get")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	//Assets
	assetHandler := http.FileServer(http.Dir("./public"))
	assetHandler = http.StripPrefix("/public", assetHandler)
	r.PathPrefix("/public/").Handler(assetHandler)

	http.ListenAndServe(":3003", r)

}
