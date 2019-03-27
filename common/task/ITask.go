package task

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	pubsub "github.com/anthonydenecheau/gopubsub/common/config"
	"github.com/anthonydenecheau/gopubsub/common/repository"
)

// Task is a parent class
type Task struct {
	dr         repository.SyncRepository
	pubService pubsub.PubSubService
	log        *logrus.Logger
	closed     chan struct{}
	wg         sync.WaitGroup
	ticker     *time.Ticker
	fn         func()
}

func (t *Task) Run() {
	for {
		select {
		case <-t.closed:
			return
		case <-t.ticker.C:
			t.fn()
		}
	}

}

func (t *Task) Stop() {
	close(t.closed)
	t.wg.Wait()
}
