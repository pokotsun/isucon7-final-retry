package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	db     *sqlx.DB
	logger *zap.SugaredLogger

	PrivateIPs = []string{
		"10.174.0.11",
		"10.174.0.13",
		"10.174.0.15",
		"10.174.0.14",
	}
)

func initDB() {
	db_host := os.Getenv("ISU_DB_HOST")
	if db_host == "" {
		db_host = "127.0.0.1"
	}
	db_port := os.Getenv("ISU_DB_PORT")
	if db_port == "" {
		db_port = "3306"
	}
	db_user := os.Getenv("ISU_DB_USER")
	if db_user == "" {
		db_user = "root"
	}
	db_password := os.Getenv("ISU_DB_PASSWORD")
	if db_password != "" {
		db_password = ":" + db_password
	}

	dsn := fmt.Sprintf("%s%s@tcp(%s:%s)/isudb?parseTime=true&loc=Local&charset=utf8mb4",
		db_user, db_password, db_host, db_port)

	logger.Infof("Connecting to db: %q", dsn)
	db, _ = sqlx.Connect("mysql", dsn)
	for {
		err := db.Ping()
		if err == nil {
			break
		}
		logger.Info(err)
		time.Sleep(time.Second * 3)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(5 * time.Minute)
	logger.Info("Succeeded to connect db.")
}

func getInitializeHandler(w http.ResponseWriter, r *http.Request) {
	db.MustExec("TRUNCATE TABLE adding")
	db.MustExec("TRUNCATE TABLE buying")
	db.MustExec("TRUNCATE TABLE room_time")
	ConnMap = make(map[string]map[int]*WebSocket)
	for _, address := range PrivateIPs {
		uri := "http://" + address + "/other_initialize"
		res, err := http.Get(uri)
		defer res.Body.Close()
		if err != nil {
			logger.Info(err)
		}
	}
	w.WriteHeader(204)
}
func getOtherInitializeHandler(w http.ResponseWriter, r *http.Request) {
	ConnMap = make(map[string]map[int]*WebSocket)
	for _, v := range GoRouineFuncMap {
		v.Channel <- 1
	}
}

func getRoomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	roomName := vars["room_name"]
	path := "/ws/" + url.PathEscape(roomName)

	var i int
	err := db.QueryRow("SELECT 1 FROM room_time WHERE room_name = ?", roomName).Scan(&i)
	if err != nil {
		logger.Infow("SELECT ERROR", "err", err)
	}
	if i != 1 {
		db.Exec("INSERT INTO room_time(room_name, time) VALUES (?, 0)", roomName)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Host string `json:"host"`
		Path string `json:"path"`
	}{
		Host: "",
		Path: path,
	})
}

func wsGameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	roomName := vars["room_name"]

	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		logger.Info("Failed to upgrade", err)
		return
	}
	go serveGameConn(ws, roomName)
}

func main() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer l.Sync()
	logger = l.Sugar()

	// log.SetFlags(log.LstdFlags | log.Lshortfile)

	initDB()

	err = InitMItem()
	if err != nil {
		logger.Error(err)
	}
	for i := 0; i < 50; i++ {
		for _, item := range M_ITEM_DICT {
			item.GetPower(i)
			item.GetPrice(i)
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/initialize", getInitializeHandler)
	r.HandleFunc("/room/", getRoomHandler)
	r.HandleFunc("/room/{room_name}", getRoomHandler)
	r.HandleFunc("/ws/", wsGameHandler)
	r.HandleFunc("/ws/{room_name}", wsGameHandler)
	r.HandleFunc("/other_initialize", getOtherInitializeHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../public/")))

	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		for {
			<-ticker.C
			for key, _ := range POWER_DICT {
				logger.Infow("DEBUG",
					"POWER_DICT key", key,
				)
			}
		}
	}()

	log.Fatal(http.ListenAndServe(":5000", handlers.LoggingHandler(os.Stderr, r)))
}
