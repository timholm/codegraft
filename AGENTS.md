# codegraft

## Overview

Codegraft is a CLI tool that automates the creation of LoRA adapters for open-source language models trained on private library source code. It addresses the problem that AI coding assistants struggle with internal SDKs and unfamiliar APIs — even with perfect documentation, LLMs lack learned patterns for invoking these APIs effectively. The tool ingests library source code, automatically synthesizes training examples using graph-based data evolution (inspired by PriCoder), and produces fine-tuned models via Ollama integration. The entire pipeline is fully automated and requires no ML expertise.

## Quick Start

```bash
# Build the project
make build

# Run tests
make test
```

## Key Files

- `internal/config/config.go`: Configuration management — loads and saves JSON config, provides paths for library data
- `internal/config/defaults.go`: Default configuration values including Ollama URL, base model, server port, evolution depth, and diversity thresholds
- `internal/store/store.go`: Core data models (Symbol, Example, TuneJob) and Store interface defining persistence operations
- `internal/store/sqlite.go`: File-based JSON storage implementation for persisting symbols, training examples, and fine-tuning jobs

## How to Extend

**Adding a new storage backend:**
Implement the Store interface defined in `internal/store/store.go`. The interface requires methods for symbol operations (SaveSymbols, GetSymbols, CountSymbols, DeleteSymbols), example operations (SaveExamples, GetExamples, CountExamples, DeleteExamples), job operations (SaveJob, GetJob, ListJobs, UpdateJobStatus), and library listing (ListLibraries). Return concrete implementations of Symbol, Example, and TuneJob types.

**Adding new configuration options:**
Add fields to the Config struct in `internal/config/config.go` with JSON tags. Add corresponding defaults to `internal/config/defaults.go`. The Default() function automatically applies all defaults, and Load() merges file values with defaults for missing fields.

**Working with symbols:**
Symbols represent extracted API elements with Kind values: func, method, type, const, var, interface. Each symbol tracks its package, signature, docstring, file location, and call relationships via the Calls field.

## Testing

Run tests with:
```bash
make test
```

Tests validate the configuration system, data model persistence, and storage operations. When adding new features, create tests in the same package directory with a `_test.go` suffix that test the interface contract rather than implementation details.
