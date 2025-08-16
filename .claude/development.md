# Development Guidelines

## Testing Strategy

### Unit Tests
- Test all business logic in `internal/` packages
- Focus on critical paths that cannot fail
- Mock file system operations where possible
- Test plugin loading, theme resolution, validation

### Integration Tests  
- Test actual config file parsing/writing
- Use `tests/fixtures/` for sample configs
- Verify plugin reload commands work
- Test theme application end-to-end

### Test Coverage
- Aim for high coverage on business logic
- Don't write tests that don't improve code quality
- All shipped plugins must have tests
- User-added plugins are their responsibility

## Error Handling Patterns

### Plugin Loading
- Continue with other plugins if one fails
- Show clear error messages: "filename parse error/not found"
- Never break working system if save fails
- Log errors but don't crash application

### Config File Safety
- Always backup before modification
- Validate files can be parsed before writing
- Show what will change before applying
- Provide rollback for failed operations

## Code Style

### Go Conventions
- Follow standard Go project layout
- Use meaningful package names
- Export only what needs to be public
- Document exported functions and types

### JSON Schemas
- Validate all plugin JSON files
- Provide clear error messages for invalid JSON
- Use consistent field naming (snake_case)
- Include help text for user-facing fields

## Manual Testing Procedures

### Plugin Development
1. Create sample plugin in `./plugins/test-app/`
2. Verify plugin loads without errors
3. Test field validation works correctly
4. Confirm templates generate valid configs

### Theme Application
1. Start with clean backup of configs
2. Apply theme changes in staging area
3. Verify diffs show expected changes
4. Test rollback restores original state
5. Confirm reload commands execute properly

### UI Testing
1. Test keyboard navigation works
2. Verify color picker shows current values
3. Check form validation provides helpful feedback
4. Ensure error states are handled gracefully

## Development Workflow

### PR Requirements
- All tests must pass
- Code must be formatted (`go fmt`)
- Static analysis clean (`go vet`)
- Manual testing completed
- Documentation updated if needed

### Review Checklist
- Business logic separated from UI
- Error handling implemented
- Tests cover new functionality
- Configuration safety maintained
- Performance impact considered