package detect

import (
	"encoding/json"
	"faro/internal/pkg/types"
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

const (
	StreamName    = "FARO"
	JobSubject    = "faro.jobs"
	ResultSubject = "faro.results"
)

// NatsOrchestrator manages distributed jobs using NATS JetStream.
type NatsOrchestrator struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func NewNatsOrchestrator(url string) (*NatsOrchestrator, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	// Ensure stream exists
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     StreamName,
		Subjects: []string{JobSubject, ResultSubject},
	})
	if err != nil {
		// Stream might already exist
		fmt.Printf("Stream info: %v\n", err)
	}

	return &NatsOrchestrator{nc: nc, js: js}, nil
}

// PublishJobs pushes record pairs into the JetStream for workers to process.
func (o *NatsOrchestrator) PublishJobs(records []types.Record) error {
	for i := 0; i < len(records); i++ {
		for j := i + 1; j < len(records); j++ {
			job := Job{RecordA: records[i], RecordB: records[j]}
			data, _ := json.Marshal(job)
			_, err := o.js.Publish(JobSubject, data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ConsumeResults listens for similarity results published by workers.
func (o *NatsOrchestrator) ConsumeResults(timeout time.Duration) ([]types.SimilarityResult, error) {
	var results []types.SimilarityResult
	sub, err := o.nc.SubscribeSync(ResultSubject)
	if err != nil {
		return nil, err
	}
	defer sub.Unsubscribe()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		msg, err := sub.NextMsg(200 * time.Millisecond)
		if err != nil {
			continue // No message yet, wait until deadline
		}
		var res types.SimilarityResult
		if err := json.Unmarshal(msg.Data, &res); err == nil {
			results = append(results, res)
		}
	}

	return results, nil
}

func (o *NatsOrchestrator) Close() {
	o.nc.Close()
}
