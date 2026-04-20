package memory

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/drhunn/SEAL_HAT_LLM/internal/telemetry"
)

type RetrievalResult struct {
	RecordID      string
	RecordKind    string
	Title         string
	Summary       string
	Body          string
	Status        string
	SemanticScore float64
	ClusterScore  float64
	FinalScore    float64
}

func EncodeVector(values []float32) string {
	parts := make([]string, 0, len(values))
	for _, v := range values {
		parts = append(parts, strconv.FormatFloat(float64(v), 'f', -1, 32))
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func ZeroVector(dim int) string {
	parts := make([]string, dim)
	for i := range parts {
		parts[i] = "0"
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func (s *PostgresStore) RunCoarseToFineSearch(ctx context.Context, namespace, specialistID, vectorLiteral string, topRegions, topClusters, topRecords int) ([]RetrievalResult, error) {
	const q = `
		SELECT
			record_id::text,
			record_kind::text,
			title,
			summary,
			body,
			status::text,
			semantic_score,
			cluster_score,
			final_score
		FROM agent_core.fn_run_coarse_to_fine_search(
			$1,
			$2,
			($3)::vector,
			$4,
			$5,
			$6
		)`

	rows, err := s.db.Query(ctx, q, namespace, specialistID, vectorLiteral, topRegions, topClusters, topRecords)
	if err != nil {
		return nil, fmt.Errorf("run coarse-to-fine search: %w", err)
	}
	defer rows.Close()

	results := make([]RetrievalResult, 0, topRecords)
	for rows.Next() {
		var item RetrievalResult
		if err := rows.Scan(
			&item.RecordID,
			&item.RecordKind,
			&item.Title,
			&item.Summary,
			&item.Body,
			&item.Status,
			&item.SemanticScore,
			&item.ClusterScore,
			&item.FinalScore,
		); err != nil {
			return nil, fmt.Errorf("scan retrieval row: %w", err)
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate retrieval rows: %w", err)
	}

	return results, nil
}

func SignalsForRetrieval(specialistID, taskClass string, results []RetrievalResult, retrievalErr error, collector *telemetry.Collector) []telemetry.Signal {
	if collector == nil {
		return nil
	}
	signals := make([]telemetry.Signal, 0)
	if retrievalErr != nil {
		signals = append(signals, collector.NewSignal(specialistID, "retrieval", "retrieval", taskClass, retrievalErr.Error(), telemetry.SeverityHigh))
		return signals
	}
	if len(results) == 0 {
		signals = append(signals, collector.NewSignal(specialistID, "retrieval", "retrieval", taskClass, "retrieval returned no records", telemetry.SeverityModerate))
		return signals
	}
	if results[0].FinalScore < 0.30 {
		signals = append(signals, collector.NewSignal(specialistID, "retrieval", "retrieval", taskClass, "retrieval top score is weak", telemetry.SeverityLow, results[0].RecordID))
	}
	return signals
}
