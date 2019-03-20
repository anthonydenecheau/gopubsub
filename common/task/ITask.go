package task

import (
	"github.com/anthonydenecheau/gopubsub/common/pubsub"
	"github.com/anthonydenecheau/gopubsub/common/repository"
)

// Task is a parent class
type Task struct {
	dr         repository.SyncRepository
	pubService pubsub.PubSubService
}
