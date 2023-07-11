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
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	r := chi.NewRouter()
	apiRouter := chi.NewRouter()
	// mux := http.NewServeMux()
	r.Handle("/app/*", apiCfg.middlewareMetricsInc(handler))
	r.Handle("/app", apiCfg.middlewareMetricsInc(handler))
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/metrics", apiCfg.handlerMetrics)
	r.Mount("/api", apiRouter)
	corsMux := middlewareCors(r)

	server := &http.Server{
		Addr: ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

