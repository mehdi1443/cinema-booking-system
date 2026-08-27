package dto

import (
	"strings"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type SuggestRequest struct {
	Keyword string `form:"q" binding:"required,min=2"`
	Limit   int    `form:"limit,default=5"`
}

type SearchQueryRequest struct {
	Keyword string `form:"q" binding:"required,min=2"`
	Types   string `form:"types"`
	Limit   int    `form:"limit,default=10"`
}

func (req *SearchQueryRequest) GetDocTypes() []enum.DocType {
	if req.Types == "" {
		return nil
	}

	rawTypes := strings.Split(req.Types, ",")
	var docTypes []enum.DocType

	for _, t := range rawTypes {
		cleanType := strings.TrimSpace(strings.ToLower(t))
		if cleanType == string(enum.DocTypeMovie) ||
			cleanType == string(enum.DocTypeCinema) ||
			cleanType == string(enum.DocTypePerson) {
			docTypes = append(docTypes, enum.DocType(cleanType))
		}
	}
	return docTypes
}

type SearchResponseItem struct {
	ID       string `json:"id"`
	EntityID int64  `json:"entity_id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	ImageURL string `json:"image_url"`
}

func ToSearchResponseSlice(docs []entity.SearchDocument) []SearchResponseItem {
	if len(docs) == 0 {
		return []SearchResponseItem{}
	}

	responses := make([]SearchResponseItem, 0, len(docs))
	for _, doc := range docs {
		responses = append(responses, SearchResponseItem{
			ID:       doc.ID,
			EntityID: doc.EntityID,
			Type:     string(doc.Type),
			Title:    doc.Title,
			Subtitle: doc.Subtitle,
			ImageURL: doc.ImageURL,
		})
	}
	return responses
}
