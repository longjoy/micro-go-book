package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	Name    string
	Habits []string
}

func write(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Custom-Header", "custom")
	w.WriteHeader(201)
	user := &User{
		Name:    "aoho",
		Habits: []string{"balls", "running", "hiking"},
	}
	json, _ := json.Marshal(user)
	w.Write(json)
}

func main() {
	http.HandleFunc("/write", write)         //设置访问的路由
	err := http.ListenAndServe(":8080", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
