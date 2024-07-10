package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-processor-api/pkg/receipt"
	"testing"

	"github.com/google/uuid"
)

// Test data
var testReceipt = receipt.Receipt{
	Retailer:     "Target",
	PurchaseDate: "2022-01-01",
	PurchaseTime: "13:01",
	Items: []receipt.Item{
		{Description: "Mountain Dew 12PK", Price: "6.49"},
		{Description: "Emils Cheese Pizza", Price: "12.25"},
		{Description: "Knorr Creamy Chicken", Price: "1.26"},
		{Description: "Doritos Nacho Cheese", Price: "3.35"},
		{Description: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
	},
	Total: "35.35",
}

func TestProcessReceipt(t *testing.T) {
	app := &Application{}
	app.Initialize()

	// Create a new HTTP POST request with the test receipt payload
	jsonReceipt, _ := json.Marshal(testReceipt)
	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(jsonReceipt))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	app.Router.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect
	var resp processedReceiptResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("handler returned unexpected body: %v", rr.Body.String())
	}

	// Check if the ID is a valid UUID
	if _, err := uuid.Parse(resp.ID); err != nil {
		t.Errorf("handler returned an invalid ID: %v", resp.ID)
	}

	// Check if the receipt was stored in the in-memory database
	if _, ok := app.InMemoryDB.Get(resp.ID); !ok {
		t.Errorf("receipt not found in in-memory database with ID: %v", resp.ID)
	}
}

func TestReceiptPoints(t *testing.T) {
	app := &Application{}
	app.Initialize()

	// Manually store the test receipt in the in-memory database
	id := uuid.New().String()
	app.InMemoryDB.Save(id, testReceipt)

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", "/receipts/"+id+"/points", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	app.Router.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect
	var resp pointsResp
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("handler returned unexpected body: %v", rr.Body.String())
	}

	// Calculate expected points
	expectedPoints := testReceipt.CalculatePoints()

	// Check if the points are what we expect
	if resp.Points != expectedPoints {
		t.Errorf("handler returned unexpected points: got %v want %v", resp.Points, expectedPoints)
	}
}

func TestReceiptPointsNotFound(t *testing.T) {
	app := &Application{}
	app.Initialize()

	// Create a new HTTP GET request for a non-existent ID
	req, err := http.NewRequest("GET", "/receipts/non-existent-id/points", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	app.Router.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
