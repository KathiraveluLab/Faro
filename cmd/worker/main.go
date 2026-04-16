package main

import (
	"encoding/json"
	"faro/internal/pkg/detect"
	"faro/internal/pkg/types"
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	fmt.Printf("Faro Distributed Worker starting. Connecting to NATS at %s...\n", natsURL)

	nc, err := nats.Connect(natsURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Create pull subscriber for comparison jobs
	sub, err := js.PullSubscribe(detect.JobSubject, "faro-worker", nats.Durable("faro-worker-durable"))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	measure := detect.Levenshtein{}
	normalizer := detect.NewClinicalNormalizer()

	fmt.Println("Worker ready. Listening for comparison jobs...")

	// Graceful shutdown handling
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sig:
			fmt.Println("Shutting down worker...")
			return
		default:
			msgs, err := sub.Fetch(1, nats.MaxWait(1*time.Second))
			if err != nil {
				continue
			}

			for _, msg := range msgs {
				var job detect.Job
				if err := json.Unmarshal(msg.Data, &job); err != nil {
					msg.Term()
					continue
				}

				// Perform comparison
				nameA := normalizer.Normalize(job.RecordA.Attributes["name"])
				nameB := normalizer.Normalize(job.RecordB.Attributes["name"])
				score := measure.Compare(nameA, nameB)

				if score > 0.8 {
					result := types.SimilarityResult{
						RecordA:     job.RecordA.ID,
						RecordB:     job.RecordB.ID,
						Score:       score,
						Algorithm:   measure.Name(),
						IsDuplicate: true,
					}
					data, _ := json.Marshal(result)
					nc.Publish(detect.ResultSubject, data)
					fmt.Printf("Duplicate Found: %s <-> %s (%.2f)\n", job.RecordA.ID, job.RecordB.ID, score)
				}
				
				// Acknowledge the message
				msg.Ack()
			}
		}
	}
}
