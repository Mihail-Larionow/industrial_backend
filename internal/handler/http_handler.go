package handler

import (
	"log"
	"net/http"
)

func Execute(w http.ResponseWriter, r *http.Request) {
	log.Printf("executing")
}
