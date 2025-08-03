package main

import (
	"fmt"
	"strings"
)

// TestFunction demonstrates some code patterns for Claude to review
// @claude please review this code for best practices, performance, and Go idioms
func TestFunction(input string) string {
	// Potential issues for Claude to catch:
	// 1. No error handling
	// 2. Inefficient string operations
	// 3. Missing input validation
	// 4. Poor variable naming
	
	var result string
	
	// Inefficient string concatenation in loop
	for i := 0; i < 10; i++ {
		result = result + input + fmt.Sprintf("_%d", i)
	}
	
	// Case conversion without checking if needed
	result = strings.ToUpper(result)
	
	return result
}

// AnotherTestFunction with more potential improvements
func AnotherTestFunction(data []string) {
	// No return value, just prints
	// Missing nil check
	for i := 0; i < len(data); i++ {
		fmt.Println(data[i])
	}
}

// UserData represents a user - missing validation
type UserData struct {
	Name  string
	Email string
	Age   int
}

// ProcessUser processes user data without proper validation
func ProcessUser(user UserData) {
	// Missing validation
	// No error handling
	// Direct printing instead of proper logging
	fmt.Printf("Processing user: %s with email %s\n", user.Name, user.Email)
}