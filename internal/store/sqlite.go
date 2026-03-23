package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// fileStore is a JSON-file-backed implementation of Store.
// Each library gets its own directory with separate JSON files.
type fileStore struct {
	dataDir string
}

// NewFileStore creates a new file-backed store rooted at dataDir.
func NewFileStore(dataDir string) (Store, error) {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("creating data dir: %w", err)
	}
	return &fileStore{dataDir: dataDir}, nil
}

func (s *fileStore) libraryDir(library string) string {
	return filepath.Join(s.dataDir, "libraries", library)
}

func (s *fileStore) symbolsFile(library string) string {
	return filepath.Join(s.libraryDir(library), "symbols.json")
}

func (s *fileStore) examplesFile(library string) string {
	return filepath.Join(s.libraryDir(library), "examples.json")
}

func (s *fileStore) jobsFile() string {
	return filepath.Join(s.dataDir, "jobs.json")
}

func writeJSON(path string, v interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating dir for %s: %w", path, err)
	}
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshaling to %s: %w", path, err)
	}
	return os.WriteFile(path, data, 0o644)
}

func readJSON(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}
	return json.Unmarshal(data, v)
}

// SaveSymbols stores symbols for a library, replacing any existing set.
func (s *fileStore) SaveSymbols(symbols []Symbol) error {
	if len(symbols) == 0 {
		return nil
	}
	return writeJSON(s.symbolsFile(symbols[0].Library), symbols)
}

// GetSymbols returns all symbols for a library.
func (s *fileStore) GetSymbols(library string) ([]Symbol, error) {
	var symbols []Symbol
	if err := readJSON(s.symbolsFile(library), &symbols); err != nil {
		return nil, err
	}
	return symbols, nil
}

// CountSymbols returns the number of stored symbols for a library.
func (s *fileStore) CountSymbols(library string) (int, error) {
	syms, err := s.GetSymbols(library)
	return len(syms), err
}

// DeleteSymbols removes all symbols for a library.
func (s *fileStore) DeleteSymbols(library string) error {
	path := s.symbolsFile(library)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// SaveExamples appends examples for a library.
func (s *fileStore) SaveExamples(examples []Example) error {
	if len(examples) == 0 {
		return nil
	}
	library := examples[0].Library
	existing, err := s.GetExamples(library)
	if err != nil {
		return err
	}
	existing = append(existing, examples...)
	return writeJSON(s.examplesFile(library), existing)
}

// GetExamples returns all training examples for a library.
func (s *fileStore) GetExamples(library string) ([]Example, error) {
	var examples []Example
	if err := readJSON(s.examplesFile(library), &examples); err != nil {
		return nil, err
	}
	return examples, nil
}

// CountExamples returns the number of training examples for a library.
func (s *fileStore) CountExamples(library string) (int, error) {
	ex, err := s.GetExamples(library)
	return len(ex), err
}

// DeleteExamples removes all examples for a library.
func (s *fileStore) DeleteExamples(library string) error {
	path := s.examplesFile(library)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

type jobsData struct {
	Jobs []TuneJob `json:"jobs"`
}

func (s *fileStore) loadJobs() ([]TuneJob, error) {
	var data jobsData
	if err := readJSON(s.jobsFile(), &data); err != nil {
		return nil, err
	}
	return data.Jobs, nil
}

func (s *fileStore) saveJobs(jobs []TuneJob) error {
	return writeJSON(s.jobsFile(), jobsData{Jobs: jobs})
}

// SaveJob creates or updates a tuning job.
func (s *fileStore) SaveJob(job *TuneJob) error {
	jobs, err := s.loadJobs()
	if err != nil {
		return err
	}
	for i, j := range jobs {
		if j.ID == job.ID {
			jobs[i] = *job
			return s.saveJobs(jobs)
		}
	}
	jobs = append(jobs, *job)
	return s.saveJobs(jobs)
}

// GetJob retrieves a job by ID.
func (s *fileStore) GetJob(id string) (*TuneJob, error) {
	jobs, err := s.loadJobs()
	if err != nil {
		return nil, err
	}
	for _, j := range jobs {
		if j.ID == id {
			jj := j
			return &jj, nil
		}
	}
	return nil, fmt.Errorf("job %s not found", id)
}

// ListJobs returns all jobs, sorted newest first.
func (s *fileStore) ListJobs() ([]TuneJob, error) {
	jobs, err := s.loadJobs()
	if err != nil {
		return nil, err
	}
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].CreatedAt.After(jobs[j].CreatedAt)
	})
	return jobs, nil
}

// UpdateJobStatus updates the status and optional error for a job.
func (s *fileStore) UpdateJobStatus(id string, status JobStatus, errMsg string) error {
	jobs, err := s.loadJobs()
	if err != nil {
		return err
	}
	for i, j := range jobs {
		if j.ID == id {
			jobs[i].Status = status
			jobs[i].Error = errMsg
			jobs[i].UpdatedAt = time.Now()
			return s.saveJobs(jobs)
		}
	}
	return fmt.Errorf("job %s not found", id)
}

// ListLibraries returns the names of all libraries that have been ingested.
func (s *fileStore) ListLibraries() ([]string, error) {
	libDir := filepath.Join(s.dataDir, "libraries")
	entries, err := os.ReadDir(libDir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading library dir: %w", err)
	}
	var libs []string
	for _, e := range entries {
		if e.IsDir() {
			libs = append(libs, e.Name())
		}
	}
	return libs, nil
}

// Close is a no-op for the file store.
func (s *fileStore) Close() error { return nil }
