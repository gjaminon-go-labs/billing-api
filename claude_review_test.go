package main

import (
	"fmt"
	"strings"
	"time"
)

// UserProcessor handles user data processing
// This code has intentional issues for Claude to review
type UserProcessor struct {
	users []User
}

type User struct {
	Name  string
	Email string
	Age   int
}

// ProcessUsers processes a list of users with several potential improvements
func (up *UserProcessor) ProcessUsers(userList []User) {
	// Issue 1: No input validation
	for i := 0; i < len(userList); i++ {
		// Issue 2: Inefficient loop iteration
		user := userList[i]
		
		// Issue 3: No error handling
		up.processUser(user)
	}
}

func (up *UserProcessor) processUser(user User) {
	// Issue 4: Direct printing instead of logging
	fmt.Printf("Processing user: %s\n", user.Name)
	
	// Issue 5: Inefficient string operations
	var result string
	for i := 0; i < 10; i++ {
		result = result + user.Name + "_" + fmt.Sprintf("%d", i)
	}
	
	// Issue 6: Unnecessary string conversion
	email := strings.ToLower(user.Email)
	if email != user.Email {
		user.Email = email
	}
	
	// Issue 7: Magic numbers
	if user.Age > 65 {
		fmt.Println("Senior user")
	}
	
	// Issue 8: Potential race condition (if used concurrently)
	up.users = append(up.users, user)
}

// ValidateEmail has validation issues
func ValidateEmail(email string) bool {
	// Issue 9: Very basic email validation
	return strings.Contains(email, "@")
}

// GetUserStats returns statistics about users
func (up *UserProcessor) GetUserStats() map[string]interface{} {
	// Issue 10: Using interface{} instead of proper types
	stats := make(map[string]interface{})
	
	stats["total"] = len(up.users)
	stats["timestamp"] = time.Now()
	
	// Issue 11: No nil check
	averageAge := 0
	for _, user := range up.users {
		averageAge += user.Age
	}
	stats["average_age"] = averageAge / len(up.users) // Potential division by zero
	
	return stats
}

// BatchProcessor demonstrates more issues
func BatchProcessor(users []User) error {
	// Issue 12: Poor error handling pattern
	if users == nil {
		return fmt.Errorf("users is nil")
	}
	
	processor := &UserProcessor{}
	
	// Issue 13: No context or timeout for potentially long operation
	for _, user := range users {
		if !ValidateEmail(user.Email) {
			// Issue 14: Silent failure
			continue
		}
		processor.processUser(user)
	}
	
	return nil
}