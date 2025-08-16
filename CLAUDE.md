# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Palettesmith is a Go-based TUI application for theme management with a plugin-based architecture. The TUI is the first interface; a separate UI will follow, so display logic must stay separate from business logic.

## Core Architecture Rules

- **Separation of Concerns**: Business logic in `internal/`, UI implementations in `ui/`
- **Plugin System**: JSON-based plugin discovery from `./plugins/<id>/` directories
- **Theme Safety**: Never break existing user configurations - protection first, enhancement second
- **Modular Design**: Each component must be independently testable and UI-agnostic

## Development Commands

- `go run cmd/palettesmith/main.go` - Run the application
- `go build -o palettesmith cmd/palettesmith/main.go` - Build executable  
- `go test ./...` - Run all tests
- `go vet ./...` - Static analysis
- `go fmt ./...` - Format code

## Key Dependencies

- **Bubble Tea**: TUI framework (`github.com/charmbracelet/bubbletea`)
- **Bubbles**: UI components (`github.com/charmbracelet/bubbles`)
- **Lipgloss**: Terminal styling (`github.com/charmbracelet/lipgloss`)

## Import Additional Documentation

@.claude/architecture.md
@.claude/development.md