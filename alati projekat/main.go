package main

import (
	"alati_projekat/handlers"
	"alati_projekat/model"
	"alati_projekat/repositories"
	"alati_projekat/services"
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Kreiranje repozitorijuma i servisa
	repo := repositories.NewConfigInMemRepository()
	service := services.NewConfigService(repo)
	repoG := repositories.NewConfigGroupInMemRepository()
	servicesG := services.NewConfigGroupService(repoG)

	// Dodavanje testne konfiguracije
	params := make(map[string]string)
	params["username"] = "pera"
	params["port"] = "5432"

	config := model.Config{
		Name:    "db_config",
		Version: 2,
		Params:  params,
	}

	config2 := model.Config{
		Name:    "db_config2",
		Version: 3,
		Params:  params,
	}

	configMap := make(map[string]model.Config)
	configMap["conf1"] = config
	configMap["conf2"] = config2

	// Pretvaranje mape configMap u slice
	var configs []model.Config
	for _, conf := range configMap {
		configs = append(configs, conf)
	}

	group := model.ConfigGroup{
		Name:    "db_cg",
		Version: 2,
		Configs: configs,
	}

	service.Add(config2)
	service.Add(config)
	servicesG.Add(group)

	// Kreiranje rukovaoca
	handler := handlers.NewConfigHandler(service)
	handlerG := handlers.NewConfigGruopHandler(servicesG, service)

	// Kreiranje rutera
	router := mux.NewRouter()
	router.HandleFunc("/configs/{name}/{version}", handler.Get).Methods("GET")
	router.HandleFunc("/configGroups/{name}/{version}", handlerG.Get).Methods("GET")

	router.HandleFunc("/configGroups/", handlerG.Add).Methods("POST")
	router.HandleFunc("/configs/", handler.Add).Methods("POST")

	router.HandleFunc("/configGroups/{nameG}/{versionG}/config/{nameC}/{versionC}", handlerG.AddConfToGroup).Methods("PUT")
	router.HandleFunc("/configGroups/{nameG}/{versionG}/{nameC}/{versionC}", handlerG.RemoveConfFromGroup).Methods("PUT")

	// Kreiranje HTTP servera
	server := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: router,
	}

	// Kanal za hvatanje signala za zaustavljanje
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Pokretanje servera u posebnoj gorutini
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Čekanje na signal zaustavljanja
	<-shutdown

	// Logovanje početka procesa graceful shutdown-a
	log.Println("Shutting down server...")

	// Pravljenje kanala za oznaku zatvaranja servera
	stop := make(chan struct{})
	go func() {
		// Postavljanje timeout-a za graceful shutdown
		timeout := 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Zatvaranje HTTP servera
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Failed to gracefully shutdown server: %v", err)
		}
		close(stop)
	}()

	// Čekanje na zatvaranje servera ili prekid izvršavanja
	<-stop

	// Logovanje završetka procesa graceful shutdown-a
	log.Println("Server shutdown completed")
}
