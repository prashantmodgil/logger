package main

import (
	"getLogs/api"
	"getLogs/bucket/minio"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)
func init(){
	minio.InitializeMinio()
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/seed", api.LogRequest(api.SeedHandler)).Methods("GET")
	router.HandleFunc("/logs", api.SearchLogs).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}




