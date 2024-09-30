package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"database/sql"
)

// Category Handlers
func getCategories(w http.ResponseWriter) {
	rows, err := db.Query("SELECT * FROM categories")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	
	var categories []Category
	
	for rows.Next() {
		var cat Category
		rows.Scan(&cat.ID, &cat.Name)
		categories = append(categories, cat)
	}

	sendResponse(w, "Categories loaded successfully", categories)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Name string `json:"name"`
	}
	var params Params
	json.NewDecoder(r.Body).Decode(&params)

	sql := "INSERT INTO categories (name) VALUES (?)"
	_, err := db.Exec(sql, params.Name)
	if err != nil {
		sendResponse(w, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	
	sendResponse(w, "Category created successfully", nil)
}

func categoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET" :
			getCategories(w)
		case "POST" :
			createCategory(w, r)
		default :
			sendResponse(w, "Method not allowed", nil, http.StatusMethodNotAllowed)
			return
	}
}
// End Category Handlers

// Item Handler
func getItems(w http.ResponseWriter) {
	rows, err := db.Query("SELECT * FROM items")
	if err != nil {
		sendResponse(w, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	
	var items []Item
	
	for rows.Next() {
		var item Item
		rows.Scan(&item.ID, &item.CategoryID, &item.Name, &item.Description, &item.Price, &item.CreatedAt)
		items = append(items, item)
	}

	sendResponse(w, "Items loaded successfully", items)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Name string `json:"name"`
		Description string `json:"description"`
		Price float64 `json:"price"`
		CategoryID int `json:"category_id"`
	}
	var params Params
	json.NewDecoder(r.Body).Decode(&params)

	sql := "INSERT INTO items (name, category_id, description, price) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(sql, params.Name, params.CategoryID, params.Description, params.Price)
	if err != nil {
		sendResponse(w, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	sendResponse(w, "Item created successfully", nil)
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	switch r.Method {
		case "GET" :
			log.Println()
			if parts[2] != "" {
				id, _ := strconv.Atoi(parts[2])
				getItemByID(w, id)
				return
			}
			getItems(w)
			return
		case "POST" :
			createItem(w, r)
			return
	}
	sendResponse(w, "Method not allowed", nil, http.StatusMethodNotAllowed)
}
// End Item Handler

// Item ID Handler
func getItemByID(w http.ResponseWriter, id int) {
	var item Item
	query := "SELECT id, category_id, name, description, price, created_at FROM items WHERE id=?"
	err := db.QueryRow(query, id).Scan(&item.ID, &item.CategoryID, &item.Name, &item.Description, &item.Price, &item.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			sendResponse(w, "Item not found!", nil, http.StatusNotFound)
			return
		}
		sendResponse(w, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	sendResponse(w, "Item loaded successfully!", item)
}

// func updateItem(w http.ResponseWriter, r *http.Request) {}
// func deleteItem(w http.ResponseWriter, r *http.Request) {}
// End Item ID Handler