package collection

import "github.com/Alkush-Pipania/carter-go/pkg/db"

// ToCollectionResponse converts a db.Collection pointer to CollectionResponse
func ToCollectionResponse(c *db.Collection) *CollectionResponse {
	if c == nil {
		return nil
	}
	return &CollectionResponse{
		ID:        c.ID.String(),
		UserID:    c.UserID.String(),
		Name:      c.Name,
		CreatedAt: c.CreatedAt.Time,
	}
}

// ToCollectionResponses converts a slice of db.Collection to []CollectionResponse
func ToCollectionResponses(collections []db.Collection) []CollectionResponse {
	responses := make([]CollectionResponse, len(collections))
	for i := range collections {
		responses[i] = *ToCollectionResponse(&collections[i])
	}
	return responses
}
