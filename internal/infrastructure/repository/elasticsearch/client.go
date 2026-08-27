package elasticsearch

import (
	"bytes"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

const IndexName = "global_search"

const indexMapping = `{
  "settings": {
    "analysis": {
      "char_filter": {
        "zero_width_spaces": {
          "type": "mapping",
          "mappings": ["\\u200C=>\\u0020"] 
        }
      },
      "filter": {
        "persian_stop": {
          "type": "stop",
          "stopwords": "_persian_"
        }
      },
      "analyzer": {
        "professional_persian": {
          "type": "custom",
          "char_filter": ["zero_width_spaces"],
          "tokenizer": "standard",
          "filter": [
            "lowercase",
            "decimal_digit",
            "arabic_normalization",
            "persian_normalization",
            "persian_stop"
          ]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "entity_id": { "type": "long" },
      "type": { "type": "keyword" },
      "title": {
        "type": "search_as_you_type",
        "analyzer": "professional_persian"
      },
      "subtitle": {
        "type": "text",
        "analyzer": "professional_persian"
      },
      "image_url": { "type": "keyword", "index": false }
    }
  }
}`

func NewClient(url string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating elastic client: %w", err)
	}

	res, err := client.Indices.Exists([]string{IndexName})
	if err != nil {
		return nil, fmt.Errorf("error checking index existence: %w", err)
	}

	if res.StatusCode == 404 {
		createRes, err := client.Indices.Create(
			IndexName,
			client.Indices.Create.WithBody(bytes.NewReader([]byte(indexMapping))),
		)
		if err != nil {
			return nil, fmt.Errorf("error creating index: %w", err)
		}
		defer createRes.Body.Close()
		log.Println("Elasticsearch index created with Persian mapping.")
	}

	return client, nil
}
