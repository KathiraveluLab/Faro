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
		if err := srv.Start(":8089"); err != nil {
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
	db.Exec("INSERT INTO records VALUES ('REC005', 'Postgres', 'Bob Wilson', 'Metformin 500mg')")
	db.Exec("INSERT INTO records VALUES ('REC006', 'Dicom', 'Bob Wilson', 'Glucophage 500 mg')")
	db.Exec("INSERT INTO records VALUES ('REC007', 'Excel', 'Alice Brown', 'Lisinopril 10mg')")
	db.Exec("INSERT INTO records VALUES ('REC008', 'Redis', 'Alice Brown', 'Prinivil 10 mg')")
	db.Exec("INSERT INTO records VALUES ('REC009', 'Oracle', 'Charlie Davis', 'Atorvastatin 20mg')")
	db.Exec("INSERT INTO records VALUES ('REC010', 'FHIR', 'Charlie Davis', 'Lipitor 20 mg')")
	db.Exec("INSERT INTO records VALUES ('REC011', 'S3', 'Eve Miller', 'Amoxicillin 500mg')")
	db.Exec("INSERT INTO records VALUES ('REC012', 'GCS', 'Eve Miller', 'Amoxil 500 mg')")
	db.Exec("INSERT INTO records VALUES ('REC013', 'Local', 'Frank White', 'Ibuprofen 400mg')")
	db.Exec("INSERT INTO records VALUES ('REC014', 'Cloud', 'Frank White', 'Advil 400 mg')")
	db.Exec("INSERT INTO records VALUES ('REC015', 'Legacy', 'Grace Hopper', 'Warfarin 5mg TAB')")
	db.Exec("INSERT INTO records VALUES ('REC016', 'Web', 'Grace Hopper', 'Warfarin 5 mg tablet')")
	db.Exec("INSERT INTO records VALUES ('REC017', 'API', 'Alan Turing', 'Atorvastatin 40mg')")
	db.Exec("INSERT INTO records VALUES ('REC018', 'App', 'Alan Turing', 'Atorvastaten 40 mg')")
	db.Exec("INSERT INTO records VALUES ('REC019', 'CSV2', 'Ada Lovelace', 'Metoprolol 25mg')")
	db.Exec("INSERT INTO records VALUES ('REC020', 'DB2', 'Ada Lovelace', 'Metoprolol Tartrate 25mg')")

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
		// Save to store
		store.PutDuplicate(res)
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
			store.PutDuplicate(res)
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

	fmt.Println("\nEngine is running. Dashboard active at http://localhost:8089")
	fmt.Println("Press Ctrl+C to terminate.")
	select {} // Keep running for the dashboard
}
