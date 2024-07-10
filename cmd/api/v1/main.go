package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	inMemoryDB "receipt-processor-api/pkg/in-memory-db"
	receipt "receipt-processor-api/pkg/receipt"
	"time"

	"github.com/google/uuid"
)

// Run the application
func main() {
	app := Application{}
	app.Initialize()
	app.Run()
}

// Middleware
func jsonValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Executing jsonValidationMiddleware")
		if r.Header.Get("Content-Type") != "application/json" {
			slog.Info("http request must be json")
			http.Error(w, "Unsupported media type, please use application/json", http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Application contains the application state and dependencies.
type Application struct {
	Server     *http.Server
	Router     *http.ServeMux
	InMemoryDB inMemoryDB.Client
}

// Response types
type processedReceiptResponse struct {
	ID string `json:"id"`
}

type pointsResp struct {
	Points int `json:"points"`
}

// Helper functions
func castToReceipt(a any) (receipt.Receipt, bool) {
	rec, ok := a.(receipt.Receipt)
	return rec, ok
}

func (app *Application) findByID(id string) (any, bool) {
	return app.InMemoryDB.Get(id)
}

// Handlers
func (app *Application) ProcessReciept(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handling POST request for handler ProcessReciept")

	// No validation is being done on the rec
	var rec receipt.Receipt
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, http.StatusText(http.StatusBadRequest))
		slog.Error("Error decoding payload in handler ProcessReciept", "error", err.Error())
		return
	}

	// Generate a new uuid
	id := uuid.New().String()
	// Store place key/value in in-memory db
	// Key is id and value is the Receipt
	// Ideally this would undergo a validation step before being
	// persisted so clients consuming this downstream can trust the value
	app.InMemoryDB.Save(id, rec)
	// Encode and return the response
	if err := json.NewEncoder(w).Encode(processedReceiptResponse{ID: id}); err != nil {
		slog.Error("Error encoding payload in handler ProcessReciept", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))
	}
}

func (app *Application) ReceiptPoints(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handling GET request for handler ReceiptPoints")
	// Look up the receipt ID in the persistence layer
	value, exists := app.findByID(r.PathValue("id"))

	if !exists {
		slog.Error("handler ReceiptPoints", "no such value exists with key", r.PathValue("id"))
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, http.StatusText(http.StatusNotFound))
		return
	}

	// Case the value to a Receipt type
	receipt, ok := castToReceipt(value)

	if !ok {
		slog.Error("handler ReceiptPoints", "could not cast value from db to reciept", value)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))
	}

	// Calculate the points for the receipt and return the response
	// or an error
	if err := json.NewEncoder(w).Encode(pointsResp{Points: receipt.CalculatePoints()}); err != nil {
		slog.Error("handler ReceiptPoints", "could not calculate points for reciept", err.Error(), "recipt", receipt)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))
	}
}

// Initialize the application
func (app *Application) Initialize() {
	app.InMemoryDB = inMemoryDB.NewClient()
	app.Router = http.NewServeMux()
	app.Server = &http.Server{
		Handler:      app.Router,
		Addr:         ":8080",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	app.Router.Handle("POST /receipts/process", jsonValidationMiddleware(http.HandlerFunc(app.ProcessReciept)))
	app.Router.HandleFunc("GET /receipts/{id}/points", app.ReceiptPoints)
}

// Run the application
func (a *Application) Run() {
	// Run the server in a goroutine so that it doesn't block
	go func() {
		slog.Info("Server is running on", slog.Any("port:", a.Server.Addr))
		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Could not listen on %s: %v\n", a.Server.Addr, err)
		}
	}()

	// Channel to listen for OS interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal
	<-c

	// Create a context with a timeout for the graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	a.Server.SetKeepAlivesEnabled(false)
	if err := a.Server.Shutdown(ctx); err != nil {
		slog.Error("Could not gracefully shutdown the server", "error", err.Error())
	}

	slog.Info("Server stopped")
}
