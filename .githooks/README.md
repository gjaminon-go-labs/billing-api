# Git Hooks for Quality Assurance

This directory contains Git hooks that help maintain code quality by running checks before certain Git operations.

## Available Hooks

### pre-push
Runs before `git push` to ensure code quality:
- Runs `go vet` to catch potential bugs
- Checks code formatting with `gofmt`  
- Runs unit tests

## Installation

### Option 1: Configure Git to use this directory (Recommended)
```bash
# From the billing-api directory
git config core.hooksPath .githooks
```

### Option 2: Copy to .git/hooks/
```bash
# From the billing-api directory
cp .githooks/pre-push .git/hooks/
chmod +x .git/hooks/pre-push
```

## Bypassing Hooks (Emergency Only)

If you need to push without running checks (not recommended):
```bash
git push --no-verify
```

**⚠️ Warning**: Only bypass hooks in emergencies. The CI pipeline will still run these checks and fail if there are issues.

## Benefits

- **No CI Surprises**: Catch issues locally before they fail in CI
- **Faster Feedback**: Know immediately if code meets quality standards
- **Save CI Minutes**: Don't waste CI resources on code that won't pass
- **Team Consistency**: Everyone follows the same quality standards

## Troubleshooting

If the hook isn't running:
1. Check it's executable: `ls -l .githooks/pre-push`
2. Verify Git configuration: `git config core.hooksPath`
3. Ensure you're in the right directory when configuring