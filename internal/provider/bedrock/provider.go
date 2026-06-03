package bedrock

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
	"github.com/uvalib/virgo4-pool-kb-ws/internal/provider"
)

// Provider retrieves image metadata from an AWS Bedrock Knowledge Base.
type Provider struct {
	Region          string
	KnowledgeBaseID string
	client          *bedrockagentruntime.Client
}

// New creates a Bedrock KB provider for the given region and knowledge base id.
func New(region, knowledgeBaseID string) (*Provider, error) {
	if knowledgeBaseID == "" {
		return nil, fmt.Errorf("knowledge base id is required")
	}
	region = defaultString(region, "us-east-1")
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}
	return &Provider{
		Region:          region,
		KnowledgeBaseID: knowledgeBaseID,
		client:          bedrockagentruntime.NewFromConfig(cfg),
	}, nil
}

// Search implements provider.KnowledgeBaseProvider.
func (p *Provider) Search(ctx context.Context, query string, limit int, threshold float64) ([]provider.Hit, error) {
	if limit <= 0 {
		limit = 20
	}
	log.Printf("[KB] retrieve kb=%s query=%q limit=%d", p.KnowledgeBaseID, query, limit)

	input := &bedrockagentruntime.RetrieveInput{
		KnowledgeBaseId: aws.String(p.KnowledgeBaseID),
		RetrievalQuery: &types.KnowledgeBaseQuery{
			Text: aws.String(query),
		},
		RetrievalConfiguration: &types.KnowledgeBaseRetrievalConfiguration{
			VectorSearchConfiguration: &types.KnowledgeBaseVectorSearchConfiguration{
				NumberOfResults:    aws.Int32(int32(limit)),
				//OverrideSearchType: types.SearchTypeHybrid,
			},
		},
	}

	resp, err := p.client.Retrieve(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("bedrock retrieve: %w", err)
	}

	hits := make([]provider.Hit, 0, len(resp.RetrievalResults))
	for _, ref := range resp.RetrievalResults {
		hit := provider.Hit{
			ID:              extractMetadataString(ref.Metadata, "id", "image_id"),
			IIIFID:          extractMetadataString(ref.Metadata, "iiif_id"),
			IIIFImageURL:    extractMetadataString(ref.Metadata, "url_iiif_image_a", "iiif_image_url"),
			IIIFManifestURL: extractMetadataString(ref.Metadata, "url_iiif_manifest_a", "iiif_manifest_url"),
			Title:           extractMetadataString(ref.Metadata, "title", "title_a"),
			Collection:      extractMetadataString(ref.Metadata, "collection", "digital_collection_a"),
		}
		if ref.Content != nil && ref.Content.Text != nil {
			hit.Content = *ref.Content.Text
		}
		if ref.Score != nil {
			hit.Score = *ref.Score
		}
		if hit.ID == "" && hit.IIIFID != "" {
			hit.ID = hit.IIIFID
		}
		if hit.Title == "" {
			hit.Title = "Image"
		}
		if hit.Score < threshold {
			continue
		}
		if hit.ID == "" && hit.Title == "" {
			continue
		}
		hits = append(hits, hit)
	}

	sort.Slice(hits, func(i, j int) bool { return hits[i].Score > hits[j].Score })
	return hits, nil
}

func defaultString(val, def string) string {
	if val == "" {
		return def
	}
	return val
}
