# codegraft

Automatically synthesize high-quality LoRA adapters for open-source models trained on your private library source code.

## What it does

Codegraft solves the problem that AI coding assistants cannot effectively use internal SDKs and private libraries — even with perfect documentation, LLMs lack learned patterns for unfamiliar APIs. The tool ingests a private library's source code and API surface, automatically synthesizes training examples using graph-based data evolution (inspired by PriCoder), and produces LoRA adapters via Ollama integration. The entire pipeline is automated and requires no ML expertise.

## Install

```bash
go install github.com/factory/codegraft@latest
```

## Usage

Initialize configuration and ingest a library:

```bash
codegraft config set --data-dir ~/.codegraft --ollama-url http://localhost:11434
codegraft ingest mylib /path/to/mylib/source
```

Generate training examples using graph-based evolution:

```bash
codegraft synthesize mylib --count 1000 --evolution-depth 3
```

Fine-tune a model using the synthesized examples:

```bash
codegraft tune mylib --base-model codellama:7b --output-model mylib-tuned
```

Check job status:

```bash
codegraft status
```

Serve the fine-tuned model:

```bash
codegraft serve mylib-tuned --port 8080
```

## API

### Configuration

```go
type Config struct {
  DataDir            string  // Root directory for codegraft data (~/.codegraft)
  OllamaURL          string  // Ollama server URL (http://localhost:11434)
  BaseModel          string  // Base model for fine-tuning (codellama:7b)
  ServerPort         int     // Port for completion server (8080)
  LogLevel           string  // Logging level: info, debug, warn, error
  MaxExamples        int     // Maximum training examples to generate (5000)
  EvolutionDepth     int     // Depth of graph evolution (3)
  DiversityThreshold float64 // Minimum diversity score (0.7)
}

// Load reads config from file or defaults
func Load(path string) (*Config, error)

// Save writes config to file
func (c *Config) Save(path string) error

// LibraryDir returns the data directory for a library
func (c *Config) LibraryDir(library string) string

// ExamplesFile returns the path to training examples
func (c *Config) ExamplesFile(library string) string
```

### Store

```go
type Store interface {
  // Symbols
  SaveSymbols(symbols []Symbol) error
  GetSymbols(library string) ([]Symbol, error)
  CountSymbols(library string) (int, error)
  DeleteSymbols(library string) error

  // Examples
  SaveExamples(examples []Example) error
  GetExamples(library string) ([]Example, error)
  CountExamples(library string) (int, error)
  DeleteExamples(library string) error

  // Jobs
  SaveJob(job *TuneJob) error
  GetJob(id string) (*TuneJob, error)
  ListJobs() ([]TuneJob, error)
  UpdateJobStatus(id string, status JobStatus, errMsg string) error

  // Libraries
  ListLibraries() ([]string, error)

  Close() error
}
```

### Data Models

```go
// Symbol represents an extracted API symbol
type Symbol struct {
  ID         string   // Unique identifier
  Library    string   // Library name
  Package    string   // Package name
  Name       string   // Symbol name
  Kind       string   // func, method, type, const, var, interface
  Signature  string   // Function/method signature
  Docstring  string   // Documentation
  File       string   // Source file path
  Line       int      // Line number
  Language   string   // Programming language
  Calls      []string // IDs of symbols this symbol calls
  Exported   bool     // Whether symbol is exported
}

// Example is a synthesized training example
type Example struct {
  ID          string    // Unique identifier
  Library     string    // Library name
  Instruction string    // Task instruction
  Input       string    // Optional input code snippet
  Response    string    // Generated response code
  Score       float64   // Quality score
  Diversity   float64   // Diversity score
  SourceIDs   []string  // Source symbol IDs
  CreatedAt   time.Time // Creation timestamp
}

// TuneJob tracks a model fine-tuning job
type TuneJob struct {
  ID            string    // Unique job ID
  Library       string    // Library name
  BaseModel     string    // Base model name
  OutputModel   string    // Output model name
  Status        JobStatus // pending, running, done, failed
  ExampleCount  int       // Number of examples used
  Error         string    // Error message if failed
  ModelfilePath string    // Path to generated Modelfile
  CreatedAt     time.Time // Job creation time
  UpdatedAt     time.Time // Last update time
}

type JobStatus string
const (
  JobStatusPending JobStatus = "pending"
  JobStatusRunning JobStatus = "running"
  JobStatusDone    JobStatus = "done"
  JobStatusFailed  JobStatus = "failed"
)
```

## Architecture

**Configuration** (`internal/config/`)
- `config.go`: Config struct with Load/Save methods for JSON files
- `defaults.go`: Default values for all configuration parameters

**Storage** (`internal/store/`)
- `store.go`: Store interface and data models (Symbol, Example, TuneJob)
- `sqlite.go`: File-based JSON storage implementation

The architecture follows a modular pattern where configuration is separated from storage, allowing for multiple storage backend implementations.

## References

- [PriCoder Paper](https://arxiv.org/abs/2603.15159) — Graph-based data evolution for synthesizing high-quality training examples from source code

## License

MIT
