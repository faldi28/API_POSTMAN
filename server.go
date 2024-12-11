package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type Item struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

var items []Item
var currentID int

// Get all items
func getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// Get a single item by ID
func getItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for _, item := range items {
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Item not found", http.StatusNotFound)
}

// Create a new item
func createItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newItem Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	currentID++
	newItem.ID = currentID
	items = append(items, newItem)
	json.NewEncoder(w).Encode(newItem)
}

// Update an existing item by ID
func updateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for i, item := range items {
		if item.ID == id {
			var updatedItem Item
			if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			updatedItem.ID = id
			items[i] = updatedItem
			json.NewEncoder(w).Encode(updatedItem)
			return
		}
	}
	http.Error(w, "Item not found", http.StatusNotFound)
}

// Delete an item by ID
func deleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for i, item := range items {
		if item.ID == id {
			items = append(items[:i], items[i+1:]...)
			json.NewEncoder(w).Encode(map[string]string{"message": "Item deleted"})
			return
		}
	}
	http.Error(w, "Item not found", http.StatusNotFound)
}

func main() {
	r := mux.NewRouter()

	// Mock data
	items = append(items, Item{ID: 1, Name: "Laptop", Quantity: 10, Price: 15000.00})
	items = append(items, Item{ID: 2, Name: "Mouse", Quantity: 50, Price: 500.00})
	currentID = 2

	// Routes
	r.HandleFunc("/items", getItems).Methods("GET")
	r.HandleFunc("/items/{id}", getItem).Methods("GET")
	r.HandleFunc("/items", createItem).Methods("POST")
	r.HandleFunc("/items/{id}", updateItem).Methods("PUT")
	r.HandleFunc("/items/{id}", deleteItem).Methods("DELETE")

	// Start server
	fmt.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
