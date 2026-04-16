package main

import (
	"database/sql"
	"faro/internal/pkg/detect"
	"faro/internal/pkg/server"
	"faro/internal/pkg/source"
	"faro/internal/pkg/storage"
	"fmt"

	_ "modernc.org/sqlite"
)

func main() {
	fmt.Println("Faro Near-Duplicate Detection Engine (Prototype)")

	// 1. Initialize storage and components
	store := storage.NewMemoryStore()
	measure := detect.Levenshtein{}
	normalizer := detect.NewClinicalNormalizer()

	// 2. Initialize and start the Resolution Dashboard (Phase 7)
	srv := server.NewServer(store)
	go func() {
		if err := srv.Start(":8080"); err != nil {
			fmt.Printf("Dashboard error: %v\n", err)
		}
	}()

	// 3. Initialize Source Adapters (Phase 5)
	fmt.Println("\nInitializing Source Adapters...")
	db, _ := sql.Open("sqlite", ":memory:")
	db.Exec("CREATE TABLE records (id TEXT, source TEXT, patient_name TEXT, record_text TEXT)")
	db.Exec("INSERT INTO records VALUES ('REC001', 'MySQL', 'John Doe', 'ASA 81mg')")
	db.Exec("INSERT INTO records VALUES ('REC002', 'MongoDB', 'John Doe', 'Aspirin 81 mg')")
	db.Exec("INSERT INTO records VALUES ('REC003', 'CSV', 'Jane Smith', 'HCTZ 25 milligram')")
	db.Exec("INSERT INTO records VALUES ('REC004', 'SQL', 'Jane Smith', 'Hydrochlorothiazide 25mg')")

	sqlSource := source.NewSqlSource(db, "records")
	records, err := sqlSource.FetchRecords()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Fetched %d records via SQLSource.\n", len(records))

	for _, r := range records {
		store.PutRecord(r)
	}

	// 4. Run Concurrent Discovery (Phase 4) and Broadcast (Phase 7)
	fmt.Println("\nRunning concurrent discovery...")
	orchestrator := detect.NewOrchestrator(measure, normalizer, 4)
	
	finalResults := orchestrator.Run(records)
	for _, res := range finalResults {
		fmt.Printf("Duplicate Found! [%s <-> %s] Similarity: %.2f%%\n", res.RecordA, res.RecordB, res.Score*100)
		// Broadcast to Dashboard via WebSocket
		srv.BroadcastDiscovery(res)
	}

	// 5. Hierarchical Metadata (Phase 2)
	fmt.Println("\nChecking Hierarchical Metadata...")
	src := &source.MockSource{}
	imageSetA, _ := src.FetchMetadata("PAT001")
	imageSetB, _ := src.FetchMetadata("PAT001-CONFLICT")
	
	metaDetector := detect.MetadataDetector{}
	metaResults := metaDetector.CompareStudies(imageSetA, imageSetB)
	for _, res := range metaResults {
		if res.IsDuplicate {
			fmt.Printf("Hierarchy Match! [%s <-> %s] Similarity: %.2f%%\n", res.RecordA, res.RecordB, res.Score*100)
			srv.BroadcastDiscovery(res)
		}
	}

	// 6. Distributed Mode (Phase 6)
	natsURL := "nats://localhost:4222"
	natsOrch, err := detect.NewNatsOrchestrator(natsURL)
	if err == nil {
		fmt.Printf("\nNATS found at %s. Distributing jobs...\n", natsURL)
		defer natsOrch.Close()
		natsOrch.PublishJobs(records)
	}

	fmt.Println("\nEngine is running. Dashboard active at http://localhost:8080")
	fmt.Println("Press Ctrl+C to terminate.")
	select {} // Keep running for the dashboard
}
