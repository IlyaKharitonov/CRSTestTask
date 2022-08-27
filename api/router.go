package api

import (
	"net/http"
)

func RegisterHandlers(){
	http.HandleFunc("/gasSummary", gasSummary)
}



