package source

import "github.com/Alkush-Pipania/carter-go/pkg/db"

// ToSourceResponse converts a db.Source pointer to SourceResponse
func ToSourceResponse(s *db.Source) *SourceResponse {
	if s == nil {
		return nil
	}
	return &SourceResponse{
		ID:           s.ID.String(),
		UserID:       s.UserID.String(),
		CollectionID: s.CollectionID.String(),
		Type:         string(s.Type),
		Status:       string(s.Status),
		Title:        s.Title,
		OriginalUrl:  s.OriginalUrl.String,
		CreatedAt:    s.CreatedAt.Time,
	}
}

// ToSourceResponses converts a slice of db.Source to []SourceResponse
func ToSourceResponses(sources []db.Source) []SourceResponse {
	responses := make([]SourceResponse, len(sources))
	for i := range sources {
		responses[i] = *ToSourceResponse(&sources[i])
	}
	return responses
}
