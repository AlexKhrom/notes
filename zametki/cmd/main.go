package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"net/http"
	"zametki/pkg/handlers"
	"zametki/pkg/middleware"

	// "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	if err := initConfig(); err != nil {
		fmt.Printf("error initializing config : %s\n", err.Error())
		return
	}

	zapLogger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("can't zap.NewProduction()")
		return
	}
	defer func() {
		err = zapLogger.Sync()
		fmt.Println("can't  zapLogger.Sync()")
	}()

	r := mux.NewRouter()

	var db *sql.DB

	dbPort := viper.Get("db_port").(string)
	dsn := "root:root@tcp(localhost:" + dbPort + ")/zametki?"
	dsn += "&charset=utf8mb4"
	dsn += "&interpolateParams=true"

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("error: can't sql.Open:", err)
		return
	}

	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("failed connection to db :%w", err))
	}

	NoteHandler := handlers.NewNoteHandler(db)

	r.HandleFunc("/api/getNotes", NoteHandler.GetNotes).Methods("GET")
	r.HandleFunc("/api/newNote", NoteHandler.NewNote).Methods("POST")
	r.HandleFunc("/api/deleteNote", NoteHandler.DeleteNotes).Methods("POST")

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(".././static"))))

	handler := middleware.LofInfo(r)
	port := viper.Get("port").(string)

	fmt.Println("start serv on port " + port)
	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		fmt.Println("can't Listen and server")
		return
	}
}

func initConfig() error {
	viper.AddConfigPath("../configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
