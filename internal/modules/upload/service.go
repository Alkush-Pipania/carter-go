package upload

import (
	"github.com/Alkush-Pipania/carter-go/pkg/rabbitmq"
	"github.com/Alkush-Pipania/carter-go/pkg/s3"
)

type Repository interface {
}

type service struct {
	repo      Repository
	producer  *rabbitmq.Producer
	presigner *s3.Presigner
}

func NewService(repo repository, producer *rabbitmq.Producer, presign *s3.Presigner) *service {
	return &service{
		repo:      repo,
		producer:  producer,
		presigner: presign,
	}
}
