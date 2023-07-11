package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{
		fileserverHits: 0,
	}
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	router := chi.NewRouter()
	apiRouter := chi.NewRouter()
	// mux := http.NewServeMux()
	router.Handle("/app/*", fsHandler)
	router.Handle("/app", fsHandler)
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/metrics", apiCfg.handlerMetrics)
	router.Mount("/api", apiRouter)
	corsMux := middlewareCors(router)

	server := &http.Server{
		Addr: ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

