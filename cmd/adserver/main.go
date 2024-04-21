package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	//mux.HandleFunc("GET /sources/{id}/campaigns", sourceCampaignsHandler)

	if err := http.ListenAndServe(":9999", mux); err != nil {
		log.Fatal(err)
	}
}

//func sourceCampaignsHandler(w http.ResponseWriter, r *http.Request) {
//	r.PathValue("id")
//}
