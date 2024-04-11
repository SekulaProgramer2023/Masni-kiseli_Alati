package main

import (
	"alati_projekat/handlers"
	"alati_projekat/model"
	"alati_projekat/repositories"
	"alati_projekat/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

func main() {
	// Kreiranje kanala za kontrolu graceful shutdown-a
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Kreiranje repozitorijuma i servisa
	repo := repositories.NewConfigInMemRepository()
	service := services.NewConfigService(repo)

	// Dodavanje testne konfiguracije
	params := make(map[string]string)
	params["username"] = "pera"
	params["port"] = "5432"
	config := model.Config{
		Name:    "db_config",
		Version: 2,
		Params:  params,
	}
	service.Add(config)

	// Kreiranje rukovaoca
	handler := handlers.NewConfigHandler(service)

	// Kreiranje rutera
	router := mux.NewRouter()
	router.HandleFunc("/configs/{name}/{version}", handler.Get).Methods("GET")
	router.HandleFunc("/configs/", handler.Add).Methods("POST")
	router.HandleFunc("/configs/{name}/{version}", handler.Delete).Methods("DELETE")

	// Kreiranje HTTP servera
	server := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: router,
	}

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
		// Zatvaranje HTTP servera
		if err := server.Shutdown(nil); err != nil {
			log.Fatalf("Failed to gracefully shutdown server: %v", err)
		}
		close(stop)
	}()

	// Čekanje na zatvaranje servera ili prekid izvršavanja
	<-stop

	// Logovanje završetka procesa graceful shutdown-a
	log.Println("Server shutdown completed")

// Kreiranje HTTP servera
server := &http.Server{
	Addr:    "0.0.0.0:8000",
	Handler: router,
}

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
	// Zatvaranje HTTP servera
	if err := server.Shutdown(nil); err != nil {
		log.Fatalf("Failed to gracefully shutdown server: %v", err)
	}
	close(stop)
}()

// Čekanje na zatvaranje servera ili prekid izvršavanja
<-stop

// Logovanje završetka procesa graceful shutdown-a
log.Println("Server shutdown completed")
}
