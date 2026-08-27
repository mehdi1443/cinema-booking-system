package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/elastic"
)

type searchRepo struct {
	client *elasticsearch.Client
}

func NewSearchRepository(client *elasticsearch.Client) elastic.SearchRepository {
	return &searchRepo{client: client}
}

func (r *searchRepo) Index(ctx context.Context, doc entity.SearchDocument) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	res, err := r.client.Index(
		IndexName,
		bytes.NewReader(data),
		r.client.Index.WithDocumentID(doc.ID),
		r.client.Index.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}
	return nil
}

func (r *searchRepo) Delete(ctx context.Context, id string) error {
	res, err := r.client.Delete(IndexName, id, r.client.Delete.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (r *searchRepo) Suggest(ctx context.Context, keyword string, limit int) ([]entity.SearchDocument, error) {
	query := map[string]interface{}{
		"size": limit,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": keyword,
				"type":  "bool_prefix",
				"fields": []string{
					"title",
					"title._2gram",
					"title._3gram",
				},
			},
		},
	}

	return r.executeSearch(ctx, query)
}

func (r *searchRepo) Search(ctx context.Context, keyword string, types []enum.DocType, limit int) ([]entity.SearchDocument, int64, error) {
	mustQuery := map[string]interface{}{
		"multi_match": map[string]interface{}{
			"query":     keyword,
			"fields":    []string{"title^3", "subtitle"},
			"fuzziness": "AUTO",
		},
	}

	boolQuery := map[string]interface{}{
		"must": mustQuery,
	}

	if len(types) > 0 {
		boolQuery["filter"] = map[string]interface{}{
			"terms": map[string]interface{}{
				"type": types,
			},
		}
	}

	query := map[string]interface{}{
		"size": limit,
		"query": map[string]interface{}{
			"bool": boolQuery,
		},
	}

	docs, total, err := r.executeSearchWithTotal(ctx, query)
	return docs, total, err
}

func (r *searchRepo) executeSearch(ctx context.Context, query map[string]interface{}) ([]entity.SearchDocument, error) {
	docs, _, err := r.executeSearchWithTotal(ctx, query)
	return docs, err
}

func (r *searchRepo) executeSearchWithTotal(ctx context.Context, query map[string]interface{}) ([]entity.SearchDocument, int64, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(IndexName),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("search failed: %s", res.String())
	}

	var result struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source entity.SearchDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	docs := make([]entity.SearchDocument, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		docs = append(docs, hit.Source)
	}

	return docs, result.Hits.Total.Value, nil
}
