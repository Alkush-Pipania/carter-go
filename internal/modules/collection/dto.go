package collection

import (
	"time"

	"github.com/Alkush-Pipania/carter-go/pkg/db"
)

type CollectionResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func ToCollectionResponse(c db.Collection) CollectionResponse {
	return CollectionResponse{
		ID:        uuidToString(c.ID),
		UserID:    uuidToString(c.UserID),
		Name:      c.Name,
		CreatedAt: c.CreatedAt.Time,
	}
}

func ToCollectionResponses(collections []db.Collection) []CollectionResponse {
	responses := make([]CollectionResponse, len(collections))
	for i, c := range collections {
		responses[i] = ToCollectionResponse(c)
	}
	return responses
}

func uuidToString(u interface{ String() string }) string {
	return u.String()
}
