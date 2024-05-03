package main

import (
	"log"
	"net/http"

	"github.com/gopay/internal"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/", internal.Index)
	router.GET("/accounts", internal.AccIndex)

	log.Fatal(http.ListenAndServe(":8080", router))
}
