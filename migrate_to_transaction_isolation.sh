#!/bin/bash

# Migrate all integration tests to use transaction isolation

# List of files to migrate
files=(
    "tests/integration/api/client_handler_test.go"
    "tests/integration/api/client_list_handler_test.go"
    "tests/integration/api/client_delete_handler_test.go"
    "tests/integration/api/client_update_handler_test.go"
    "tests/integration/repository/client_repository_test.go"
    "tests/integration/repository/client_repository_crud_test.go"
    "tests/integration/storage/postgres_storage_test.go"
)

for file in "${files[@]}"; do
    echo "Migrating $file..."
    
    # Replace NewIntegrationTestServer() with WithTransaction pattern
    sed -i 's/server := testhelpers\.NewIntegrationTestServer()/stack, cleanup := testhelpers.WithTransaction(t)\n\tdefer cleanup()\n\tserver := stack.HTTPServer/g' "$file"
    
    # Replace NewIntegrationTestStack() with WithTransaction pattern
    sed -i 's/stack := testhelpers\.NewIntegrationTestStack()/stack, cleanup := testhelpers.WithTransaction(t)\n\tdefer cleanup()/g' "$file"
    
    # Replace NewCleanIntegrationTestStack() with WithTransaction pattern
    sed -i 's/stack := testhelpers\.NewCleanIntegrationTestStack()/stack, cleanup := testhelpers.WithTransaction(t)\n\tdefer cleanup()/g' "$file"
done

echo "Migration complete!"