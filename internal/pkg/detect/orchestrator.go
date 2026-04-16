package detect

import (
	"faro/internal/pkg/types"
	"sync"
)

// Job represents a pair of records to be compared by a worker.
type Job struct {
	RecordA types.Record
	RecordB types.Record
}

// Orchestrator manages a pool of workers for concurrent duplicate detection.
type Orchestrator struct {
	Measure    SimilarityMeasure
	Normalizer Normalizer
	NumWorkers int
}

func NewOrchestrator(measure SimilarityMeasure, norm Normalizer, workers int) *Orchestrator {
	return &Orchestrator{
		Measure:    measure,
		Normalizer: norm,
		NumWorkers: workers,
	}
}

// Run executes the detection process in parallel.
func (o *Orchestrator) Run(records []types.Record) []types.SimilarityResult {
	jobs := make(chan Job, len(records)*2)
	results := make(chan types.SimilarityResult, len(records)*2)
	var wg sync.WaitGroup

	// Start workers
	for w := 1; w <= o.NumWorkers; w++ {
		wg.Add(1)
		go o.worker(jobs, results, &wg)
	}

	// Generate jobs (Cartesian product for simplicity in this prototype)
	go func() {
		for i := 0; i < len(records); i++ {
			for j := i + 1; j < len(records); j++ {
				jobs <- Job{RecordA: records[i], RecordB: records[j]}
			}
		}
		close(jobs)
	}()

	// Wait for workers in a separate goroutine
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var finalResults []types.SimilarityResult
	for res := range results {
		if res.IsDuplicate {
			finalResults = append(finalResults, res)
		}
	}

	return finalResults
}

func (o *Orchestrator) worker(jobs <-chan Job, results chan<- types.SimilarityResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		nameA := o.Normalizer.Normalize(job.RecordA.Attributes["name"])
		nameB := o.Normalizer.Normalize(job.RecordB.Attributes["name"])
		score := o.Measure.Compare(nameA, nameB)

		if score > 0.8 {
			results <- types.SimilarityResult{
				RecordA:     job.RecordA.ID,
				RecordB:     job.RecordB.ID,
				Score:       score,
				Algorithm:   o.Measure.Name(),
				IsDuplicate: true,
			}
		}
	}
}
