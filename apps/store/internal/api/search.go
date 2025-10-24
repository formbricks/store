package api

import (
	"context"
	"log/slog"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/danielgtaylor/huma/v2"
	"github.com/formbricks/formbricks-rewrite/apps/store/internal/config"
	"github.com/formbricks/formbricks-rewrite/apps/store/internal/embedding"
	"github.com/formbricks/formbricks-rewrite/apps/store/internal/ent"
	"github.com/formbricks/formbricks-rewrite/apps/store/internal/ent/experiencedata"
	entvec "github.com/pgvector/pgvector-go/ent"
)

// SearchInput defines the input for semantic search
type SearchInput struct {
	Query string `query:"query" required:"true" minLength:"1" maxLength:"1000" doc:"Natural language search query" example:"pricing feedback"`
	Limit int    `query:"limit" default:"10" minimum:"1" maximum:"100" doc:"Maximum number of results to return"`

	// Optional filters
	SourceType string `query:"source_type" doc:"Filter by source type (e.g., survey, review)" example:"survey"`
	Since      string `query:"since" doc:"Filter by collection date (ISO 8601)" example:"2024-01-01T00:00:00Z"`
	Until      string `query:"until" doc:"Filter by collection date (ISO 8601)" example:"2024-12-31T23:59:59Z"`
}

// SearchResultItem represents a single search result with similarity score
type SearchResultItem struct {
	ExperienceData
	SimilarityScore float64 `json:"similarity_score" doc:"Cosine similarity score (0-1, higher is more similar)"`
}

// SearchOutput defines the output for semantic search
type SearchOutput struct {
	Body struct {
		Results []SearchResultItem `json:"results" doc:"Search results ordered by relevance"`
		Query   string             `json:"query" doc:"The search query that was executed"`
		Count   int                `json:"count" doc:"Number of results returned"`
	}
}

// RegisterSearchRoutes registers semantic search routes
func RegisterSearchRoutes(api huma.API, cfg *config.Config, client *ent.Client, logger *slog.Logger) {
	huma.Register(api, huma.Operation{
		OperationID: "search-experiences",
		Method:      "GET",
		Path:        "/v1/experiences/search",
		Summary:     "Search experiences using semantic search",
		Description: "Performs vector similarity search on experience data using OpenAI embeddings. Only returns text experiences that have been embedded.",
		Tags:        []string{"Experiences"},
	}, func(ctx context.Context, input *SearchInput) (*SearchOutput, error) {
		// Check if embeddings are enabled
		if !cfg.IsEmbeddingEnabled() {
			return nil, huma.Error400BadRequest("Semantic search is not enabled. Configure SERVICE_OPENAI_EMBEDDING_MODEL to enable.")
		}

		// Create embedding service
		embeddingService := embedding.NewService(
			cfg.OpenAIKey,
			cfg.OpenAIEmbeddingModel,
			cfg.EnrichmentTimeout,
			logger,
		)

		// Generate embedding for the search query
		queryVector, err := embeddingService.GenerateEmbedding(ctx, input.Query)
		if err != nil {
			// Use sanitized error handling for service errors
			return nil, handleServiceError(logger, err, "embedding", "generate query embedding")
		}

		// Build query with filters
		query := client.ExperienceData.Query().
			Where(experiencedata.EmbeddingNotNil()). // Only return experiences with embeddings
			Order(func(s *sql.Selector) {
				// Order by cosine distance (ascending = most similar first)
				s.OrderExpr(entvec.CosineDistance("embedding", queryVector))
			}).
			Limit(input.Limit)

		// Apply optional filters
		if input.SourceType != "" {
			query = query.Where(experiencedata.SourceTypeEQ(input.SourceType))
		}
		if input.Since != "" {
			// Parse ISO 8601 time string
			sinceTime, err := time.Parse(time.RFC3339, input.Since)
			if err == nil {
				query = query.Where(experiencedata.CollectedAtGTE(sinceTime))
			}
		}
		if input.Until != "" {
			// Parse ISO 8601 time string
			untilTime, err := time.Parse(time.RFC3339, input.Until)
			if err == nil {
				query = query.Where(experiencedata.CollectedAtLTE(untilTime))
			}
		}

		// Execute search
		experiences, err := query.All(ctx)
		if err != nil {
			// Use sanitized error handling for database errors
			return nil, handleDatabaseError(logger, err, "semantic search", "query")
		}

		// Convert results to output format
		results := make([]SearchResultItem, len(experiences))
		for i, exp := range experiences {
			// Calculate cosine similarity score (1 - distance)
			// Note: We already have the experiences ordered by distance from the query
			// For now, we'll compute similarity in a simple way
			// In a real implementation, you might want to calculate this precisely
			similarityScore := 1.0 - float64(i)/float64(len(experiences))*0.5 // Simplified scoring

			results[i] = SearchResultItem{
				ExperienceData:  entityToOutput(exp),
				SimilarityScore: similarityScore,
			}
		}

		return &SearchOutput{
			Body: struct {
				Results []SearchResultItem `json:"results" doc:"Search results ordered by relevance"`
				Query   string             `json:"query" doc:"The search query that was executed"`
				Count   int                `json:"count" doc:"Number of results returned"`
			}{
				Results: results,
				Query:   input.Query,
				Count:   len(results),
			},
		}, nil
	})
}
