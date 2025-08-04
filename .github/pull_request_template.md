# Pull Request Template

## Summary
<!-- Brief description of what this PR implements -->

## TDD Workflow Completion ✅
- [ ] **EXPLORE**: Requirements research completed and documented
- [ ] **PLAN**: Implementation plan created and approved by human reviewer
- [ ] **TEST**: Failing tests written first (Red phase) 
- [ ] **CODE**: Minimal implementation completed (Green phase)
- [ ] **REFACTOR**: Code quality improvements completed
- [ ] **COMMIT & PUSH**: Feature branch created and changes committed
- [ ] **CREATE PR**: This pull request created with comprehensive description

## Implementation Details

### TDD Cycles Completed
<!-- List each TDD cycle implemented -->
1. [ ] **Cycle 1**: [Description] - Red → Green → Refactor ✅
2. [ ] **Cycle 2**: [Description] - Red → Green → Refactor ✅
3. [ ] **Cycle 3**: [Description] - Red → Green → Refactor ✅
<!-- Add more cycles as needed -->

### Architecture Changes
- [ ] **Clean Architecture**: Changes follow clean architecture principles
- [ ] **DDD Compliance**: Domain-driven design patterns maintained
- [ ] **Layer Separation**: Proper separation of concerns maintained
- [ ] **Interface Contracts**: New interfaces documented and implemented

### API Changes
- [ ] **New Endpoints**: [List any new HTTP endpoints]
- [ ] **Request/Response**: DTOs follow established patterns
- [ ] **Error Handling**: Domain errors properly mapped to HTTP responses
- [ ] **Backward Compatibility**: No breaking changes to existing APIs

## Test Coverage

### Test Types Implemented
- [ ] **Unit Tests**: All layers tested with in-memory storage
- [ ] **Integration Tests**: Real database interactions tested
- [ ] **HTTP Tests**: Full request/response cycle testing
- [ ] **Edge Cases**: Error scenarios and boundary conditions covered

### Coverage Metrics
- [ ] **Test Coverage**: ≥95% coverage for all new code
- [ ] **All Tests Pass**: Unit, integration, and HTTP tests passing ✅
- [ ] **Test Data**: External test data files used where appropriate
- [ ] **Test Isolation**: Tests run independently without side effects

## Code Quality Standards

### Code Standards
- [ ] **Go Conventions**: Code follows Go idioms and style guide
- [ ] **Project Patterns**: Consistent with existing codebase patterns
- [ ] **No Debug Code**: No console.log, print, or debug statements
- [ ] **No TODOs**: No TODO or FIXME comments in production code
- [ ] **Error Handling**: Comprehensive error handling implemented

### Documentation
- [ ] **Code Comments**: Complex logic documented where needed
- [ ] **API Documentation**: New endpoints documented if applicable
- [ ] **README Updates**: Documentation updated if needed
- [ ] **Migration Notes**: Database changes documented if any

## Review Checklist

### Functionality
- [ ] **Feature Works**: Manual testing confirms feature works as specified
- [ ] **Requirements Met**: All acceptance criteria satisfied
- [ ] **Error Cases**: Error scenarios handled appropriately
- [ ] **Performance**: No significant performance degradation

### Security & Quality
- [ ] **Security Review**: No security vulnerabilities introduced
- [ ] **Data Validation**: Input validation implemented where needed
- [ ] **SQL Injection**: Database queries use parameterized statements
- [ ] **Authentication**: Proper authorization checks if applicable

### Database Changes
- [ ] **Schema Changes**: Database migrations included if needed
- [ ] **Migration Safety**: Migrations are reversible and safe
- [ ] **Data Integrity**: Foreign key constraints respected
- [ ] **Index Performance**: Database queries optimized

## Deployment Notes
<!-- Any special deployment considerations -->
- [ ] **Migration Required**: Database migration needed before deployment
- [ ] **Configuration**: New configuration parameters documented
- [ ] **Dependencies**: No new external dependencies added
- [ ] **Environment Variables**: New environment variables documented

## Related Issues
<!-- Link to related GitHub issues -->
Closes #[issue-number]
Related to #[issue-number]

## Screenshots/Demo
<!-- If applicable, add screenshots or demo GIFs -->

---

## Review Guidelines for Reviewers

### What to Focus On
1. **TDD Compliance**: Verify all TDD workflow steps were followed
2. **Test Quality**: Ensure comprehensive test coverage and meaningful tests
3. **Architecture**: Confirm clean architecture and DDD principles
4. **Code Quality**: Check for consistency with project standards
5. **Security**: Review for potential security issues

### Approval Criteria
- ✅ All tests passing with ≥95% coverage
- ✅ Code follows project conventions and style guide
- ✅ TDD workflow completion checklist verified
- ✅ No breaking changes to existing functionality
- ✅ Comprehensive test coverage for new features

🤖 **Generated with [Claude Code](https://claude.ai/code)**

**Co-Authored-By:** Claude <noreply@anthropic.com>