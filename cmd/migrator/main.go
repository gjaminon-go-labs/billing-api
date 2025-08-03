// Database Migration CLI Tool
//
// This is a standalone CLI tool for managing database migrations.
// Provides: Migration commands (up, down, status, force), manual database management
// Features: Support for all migration operations, environment configuration
// Usage: go run cmd/migrator/main.go [command] [args]
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gjaminon-go-labs/billing-api/internal/config"
	"github.com/gjaminon-go-labs/billing-api/internal/migration"
)

const (
	cmdUp     = "up"
	cmdDown   = "down"
	cmdSteps  = "steps"
	cmdStatus = "status"
	cmdForce  = "force"
	cmdHelp   = "help"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}

func run() error {
	// Parse command line arguments
	if len(os.Args) < 2 {
		printUsage()
		return nil
	}

	command := os.Args[1]
	
	// Handle help command
	if command == cmdHelp {
		printUsage()
		return nil
	}

	// Load configuration
	environment := config.GetEnvironment()
	log.Printf("📋 Environment: %s", environment)

	appConfig, err := config.LoadConfig(environment)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create migration service
	migrationConfig := &migration.Config{
		DatabaseURL:    appConfig.Database.Host, // We'll use the database URL from config
		MigrationsPath: "database/migrations",
		SchemaName:     appConfig.Database.Schema,
	}

	// Build proper database URL
	migrationConfig.DatabaseURL = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		appConfig.Database.User,
		appConfig.Database.Password,
		appConfig.Database.Host,
		appConfig.Database.Port,
		appConfig.Database.DBName,
		appConfig.Database.SSLMode)

	migrationService, err := migration.NewService(migrationConfig)
	if err != nil {
		return fmt.Errorf("failed to create migration service: %w", err)
	}
	defer migrationService.Close()

	// Execute command
	switch command {
	case cmdUp:
		return handleUp(migrationService)
	case cmdDown:
		return handleDown(migrationService)
	case cmdSteps:
		return handleSteps(migrationService, os.Args[2:])
	case cmdStatus:
		return handleStatus(migrationService)
	case cmdForce:
		return handleForce(migrationService, os.Args[2:])
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func handleUp(service *migration.Service) error {
	log.Println("🚀 Running all pending migrations...")
	return service.Up()
}

func handleDown(service *migration.Service) error {
	log.Println("🔄 Rolling back one migration...")
	return service.Down()
}

func handleSteps(service *migration.Service, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("steps command requires number of steps")
	}

	steps, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid number of steps: %s", args[0])
	}

	return service.Steps(steps)
}

func handleStatus(service *migration.Service) error {
	status, err := service.Status()
	if err != nil {
		return err
	}

	fmt.Printf("📊 Migration Status:\n")
	fmt.Printf("   Version: %d\n", status.Version)
	fmt.Printf("   Dirty: %t\n", status.Dirty)
	fmt.Printf("   Has Schema: %t\n", status.HasSchema)
	fmt.Printf("   Message: %s\n", status.Message)

	if status.Dirty {
		fmt.Printf("⚠️  Database is in dirty state. Use 'force' command to fix.\n")
	}

	return nil
}

func handleForce(service *migration.Service, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("force command requires version number")
	}

	version, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid version number: %s", args[0])
	}

	fmt.Printf("⚠️  WARNING: This will force the migration version to %d without running migrations.\n", version)
	fmt.Printf("⚠️  This should only be used to fix dirty state or skip broken migrations.\n")
	fmt.Printf("⚠️  Are you sure you want to continue? (y/N): ")

	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" && confirm != "Y" {
		fmt.Println("Operation cancelled.")
		return nil
	}

	return service.Force(version)
}

func printUsage() {
	fmt.Printf("Database Migration CLI Tool\n\n")
	fmt.Printf("Usage: go run cmd/migrator/main.go <command> [args]\n\n")
	fmt.Printf("Commands:\n")
	fmt.Printf("  up             Run all pending migrations\n")
	fmt.Printf("  down           Roll back one migration\n")
	fmt.Printf("  steps <n>      Run n migrations (positive=up, negative=down)\n")
	fmt.Printf("  status         Show current migration status\n")
	fmt.Printf("  force <v>      Force migration version (use with caution)\n")
	fmt.Printf("  help           Show this help message\n\n")
	fmt.Printf("Environment Variables:\n")
	fmt.Printf("  ENVIRONMENT    Set environment (development, production)\n")
	fmt.Printf("  DB_HOST        Override database host\n")
	fmt.Printf("  DB_PORT        Override database port\n")
	fmt.Printf("  DB_USER        Override database user\n")
	fmt.Printf("  DB_PASSWORD    Override database password\n")
	fmt.Printf("  DB_NAME        Override database name\n\n")
	fmt.Printf("Examples:\n")
	fmt.Printf("  go run cmd/migrator/main.go up\n")
	fmt.Printf("  go run cmd/migrator/main.go down\n")
	fmt.Printf("  go run cmd/migrator/main.go steps 2\n")
	fmt.Printf("  go run cmd/migrator/main.go status\n")
	fmt.Printf("  ENVIRONMENT=production go run cmd/migrator/main.go up\n")
}

// init configures logging
func init() {
	log.SetFlags(log.LstdFlags)
	log.SetPrefix("[MIGRATOR] ")
	
	// Parse global flags for configuration
	var helpFlag = flag.Bool("help", false, "Show help")
	var versionFlag = flag.Bool("version", false, "Show version")
	
	flag.Parse()
	
	if *helpFlag {
		printUsage()
		os.Exit(0)
	}
	
	if *versionFlag {
		fmt.Println("Database Migration CLI Tool v1.0.0")
		os.Exit(0)
	}
}