// Package marathon scheduler integration
package marathon

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/op/go-logging"

	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/scheduler"
	"github.com/microscaling/microscaling/utils"
)

var log = logging.MustGetLogger("mssscheduler")

// MarathonScheduler holds Marathon API URL.
type MarathonScheduler struct {
	baseMarathonURL string
	taskBackoffs    map[string]*utils.Backoff
	sync.Mutex
}

// AppsMessage from the Marathon API.
type AppsMessage struct {
	Apps []App `json:"apps"`
}

// App from the Marathon API.
type App struct {
	ID        string `json:"id"`
	Instances int    `json:"instances"`
}

var (
	httpClient = &http.Client{
		// TODO Make timeout configurable.
		Timeout: 10 * time.Second,
	}
)

// NewScheduler returns a pointer to the scheduler.
func NewScheduler(marathonAPI string) *MarathonScheduler {
	return &MarathonScheduler{
		baseMarathonURL: getBaseMarathonURL(marathonAPI),
		taskBackoffs:    make(map[string]*utils.Backoff),
	}
}

type startStopPayload struct {
	Instances int `json:"instances"`
}

// compile-time assert that we implement the right interface
var _ scheduler.Scheduler = (*MarathonScheduler)(nil)

// InitScheduler initializes the scheduler.
func (m *MarathonScheduler) InitScheduler(task *demand.Task) (err error) {
	log.Infof("Marathon initializing task %s", task.Name)

	m.Lock()
	defer m.Unlock()

	m.taskBackoffs[task.Name] = &utils.Backoff{
		Min:    500 * time.Millisecond,
		Max:    10 * time.Second,
		Factor: 2,
	}
	return err
}

// StopStartTasks by calling the Marathon scaling API.
func (m *MarathonScheduler) StopStartTasks(tasks *demand.Tasks) error {
	// Create tasks if there aren't enough of them, and stop them if there are too many
	var tooMany []*demand.Task
	var tooFew []*demand.Task
	var err error = nil

	tasks.Lock()
	defer tasks.Unlock()
	m.Lock()
	defer m.Unlock()

	// TODO: Consider checking the number running before we start & stop
	for _, task := range tasks.Tasks {
		if task.Demand > task.Requested {
			// There aren't enough of these containers yet
			tooFew = append(tooFew, task)
		}
		if task.Demand < task.Requested {
			// there aren't enough of these containers yet
			tooMany = append(tooMany, task)
		}
	}

	// Scale down first to free up resources
	for _, task := range tooMany {
		err = m.stopStartTask(task)
		if err != nil {
			log.Errorf("Couldn't stop %s: %v ", task.Name, err)
		}
		log.Debugf("Now have %s: %d", task.Name, task.Requested)
	}

	// Now we can scale up
	for _, task := range tooFew {
		err = m.stopStartTask(task)
		if err != nil {
			log.Errorf("Couldn't start %s: %v ", task.Name, err)
		}
		log.Debugf("Now have %s: %d", task.Name, task.Requested)
	}

	return err
}

// CountAllTasks tells us how many instances of each task are currently running.
func (m *MarathonScheduler) CountAllTasks(running *demand.Tasks) error {
	var (
		err         error
		appsMessage AppsMessage
	)

	running.Lock()
	defer running.Unlock()
	m.Lock()
	defer m.Unlock()

	url := m.baseMarathonURL + "apps/"

	body, err := utils.GetJSON(url)
	if err != nil {
		log.Errorf("Error getting Marathon Apps %v", err)
		return err
	}

	err = json.Unmarshal(body, &appsMessage)
	if err != nil {
		log.Errorf("Error %v unmarshalling from %s", err, string(body[:]))
		return err
	}

	appCounts := make(map[string]int)

	// Remove leading slash from App IDs and set the instance counts.
	for _, app := range appsMessage.Apps {
		appCounts[strings.Replace(app.ID, "/", "", 1)] = app.Instances
	}

	// Set running counts. Defaults to 0 if the App does not exist.
	tasks := running.Tasks
	for _, t := range tasks {
		t.Running = appCounts[t.Name]
	}

	return err
}

// stopStartTask updates the number of running tasks using the Marathon API.
func (m *MarathonScheduler) stopStartTask(task *demand.Task) error {
	// Submit a post request to Marathon to match the requested number of the requested app
	// format looks like:
	// PUT http://marathon:8080/v2/apps/<app>
	//  Request:
	//  {
	//    "instances": 8
	//  }

	url := m.baseMarathonURL + "apps/" + task.Name
	log.Debugf("Start/stop PUT: %s", url)

	payload := startStopPayload{
		Instances: task.Demand,
	}
	w := &bytes.Buffer{}
	encoder := json.NewEncoder(w)
	err := encoder.Encode(&payload)
	if err != nil {
		log.Errorf("Failed to encode json. %v", err)
		return err
	}

	status, err := utils.PutJSON(url, w)

	b := m.taskBackoffs[task.Name]

	// Handle locked deployments by backing off until they complete.
	if status == 409 {
		dur := b.Duration(b.Attempt())

		log.Infof("Task: %s has a locked deployment - backing off for %s", task.Name, dur.String())
		time.Sleep(dur)

		return err
	} else {
		// Reset attempts when the deployment is no longer locked.
		if b.Attempt() > 0 {
			log.Infof("Task: %s succeeded set attempts to 0", task.Name)
			b.Reset()
		}

		if status > 299 {
			log.Errorf("Error response from Marathon. %v", err)
			return err
		}
	}

	// Now we've asked for this many, update the currentcount
	task.Requested = task.Demand

	return err
}

// getBaseMarathonURL returns the base API path.
func getBaseMarathonURL(marathonAPI string) string {
	return marathonAPI + "/v2/"
}
