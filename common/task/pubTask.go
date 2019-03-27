package task

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Sirupsen/logrus"

	pubsub "github.com/anthonydenecheau/gopubsub/common/config"
	"github.com/anthonydenecheau/gopubsub/common/model"
	"github.com/anthonydenecheau/gopubsub/common/service"

	syncRepository "github.com/anthonydenecheau/gopubsub/common/repository"
)

func (t Task) pubTopic() string { return t.pubService.GetTopicName() }

// PubTask is children Task
type pubTask struct {
	Task
	dogService      service.DogService
	breederService  service.BreederService
	ownerService    service.OwnerService
	parentService   service.ParentService
	pedigreeService service.PedigreeService
	titleServiceF   service.TitleService
	titleServiceE   service.TitleService
	personService   service.PersonService
}

type messages struct {
	dogs      []*model.Dog
	breeders  []*model.Breeder
	owners    []*model.Owner
	parents   []*model.Parent
	pedigrees []*model.Pedigree
	titles    []*model.Title
	action    string
}

func makeTimestamp() int64 {
	//return time.Now().UnixNano() / int64(time.Millisecond)
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// NewTask initialize all tasks
func NewTask(
	db *sql.DB, ds service.DogService,
	bs service.BreederService,
	ws service.OwnerService,
	ps service.ParentService,
	ls service.PedigreeService,
	tsf service.TitleService,
	tse service.TitleService,
	ns service.PersonService,
	pubService pubsub.PubSubService,
	log *logrus.Logger) {

	log.Info("PubTask is running ...")
	dr := syncRepository.NewOraSyncRepository(db)

	// la tâche s'exécute toutes les 5 secondes
	task := &Task{
		dr:         dr,
		pubService: pubService,
		log:        log,
		closed:     make(chan struct{}),
		ticker:     time.NewTicker(time.Second * 5),
	}

	d := pubTask{
		Task:            *task,
		dogService:      ds,
		breederService:  bs,
		ownerService:    ws,
		parentService:   ps,
		pedigreeService: ls,
		titleServiceF:   tsf,
		titleServiceE:   tse,
		personService:   ns}

	fmt.Println("PubTopic: ", d.pubTopic()) //has Task  behavior

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	// En cas d'arrêt, on attend que la tâche se termine ...
	task.wg.Add(1)
	go func() {
		defer task.wg.Done()
		task.fn = func() {
			d.syncDogChanges()
			d.syncBreederChanges()
			d.syncOwnerChanges()
			d.syncParentChanges()
			d.syncPedigreeChanges()
			d.syncTitleChangesF()
			d.syncTitleChangesE()
			d.syncPersonChanges()
		}
		task.Run()
	}()

	select {
	case sig := <-c:
		log.Infof("Got %s signal. Aborting...\n", sig)
		task.pubService.GetTopic().Stop()
		task.Stop()
	}
}
