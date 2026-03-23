// Package store provides persistent storage for codegraft data.
package store

import (
	"time"
)

// Symbol represents an extracted API symbol from a library.
type Symbol struct {
	ID        string   `json:"id"`
	Library   string   `json:"library"`
	Package   string   `json:"package"`
	Name      string   `json:"name"`
	Kind      string   `json:"kind"` // func, method, type, const, var, interface
	Signature string   `json:"signature"`
	Docstring string   `json:"docstring"`
	File      string   `json:"file"`
	Line      int      `json:"line"`
	Language  string   `json:"language"`
	Calls     []string `json:"calls,omitempty"` // IDs of symbols this symbol calls
	Exported  bool     `json:"exported"`
}

// Example is a synthesized training example.
type Example struct {
	ID          string    `json:"id"`
	Library     string    `json:"library"`
	Instruction string    `json:"instruction"`
	Input       string    `json:"input,omitempty"`
	Response    string    `json:"response"`
	Score       float64   `json:"score"`
	Diversity   float64   `json:"diversity"`
	SourceIDs   []string  `json:"source_ids,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// JobStatus represents the state of a tuning job.
type JobStatus string

const (
	JobStatusPending  JobStatus = "pending"
	JobStatusRunning  JobStatus = "running"
	JobStatusDone     JobStatus = "done"
	JobStatusFailed   JobStatus = "failed"
)

// TuneJob tracks a model fine-tuning job.
type TuneJob struct {
	ID          string    `json:"id"`
	Library     string    `json:"library"`
	BaseModel   string    `json:"base_model"`
	OutputModel string    `json:"output_model"`
	Status      JobStatus `json:"status"`
	ExampleCount int      `json:"example_count"`
	Error       string    `json:"error,omitempty"`
	ModelfilePath string  `json:"modelfile_path,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Store defines the data persistence interface for codegraft.
type Store interface {
	// Symbol operations
	SaveSymbols(symbols []Symbol) error
	GetSymbols(library string) ([]Symbol, error)
	CountSymbols(library string) (int, error)
	DeleteSymbols(library string) error

	// Example operations
	SaveExamples(examples []Example) error
	GetExamples(library string) ([]Example, error)
	CountExamples(library string) (int, error)
	DeleteExamples(library string) error

	// Job operations
	SaveJob(job *TuneJob) error
	GetJob(id string) (*TuneJob, error)
	ListJobs() ([]TuneJob, error)
	UpdateJobStatus(id string, status JobStatus, errMsg string) error

	// Library listing
	ListLibraries() ([]string, error)

	Close() error
}
