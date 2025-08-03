package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

// GetUserData handles user data requests
func GetUserData(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	
	// SQL Injection vulnerability
	query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userID)
	
	db, _ := sql.Open("postgres", "connection_string")
	rows, err := db.Query(query)
	
	// No error handling
	defer rows.Close()
	
	var userData string
	rows.Next()
	rows.Scan(&userData)
	
	// No authorization check - any user can access any data
	// Sensitive data exposure - returns raw DB data
	w.Write([]byte(userData))
}