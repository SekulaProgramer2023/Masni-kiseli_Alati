package main

import (
	"alati_projekat/handlers"
	"alati_projekat/model"
	"alati_projekat/repositories"
	"alati_projekat/services"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	repo := repositories.NewConfigInMemRepository()
	service := services.NewConfigService(repo)
	params := make(map[string]string)
	params["username"] = "pera"
	params["port"] = "5432"
	config := model.Config{
		Name:    "db_config",
		Version: 2,
		Params:  params,
	}
	service.Add(config)
	handler := handlers.NewConfigHandler(service)

	router := mux.NewRouter()

	router.HandleFunc("/configs/{name}/{version}", handler.Get).Methods("GET")

	http.ListenAndServe("0.0.0.0:8000", router)
}
