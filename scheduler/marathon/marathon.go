// Package marathon provides a scheduler using the Marathon REST API.
package marathon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/op/go-logging"

	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/scheduler"
	"github.com/microscaling/microscaling/utils"
)

var log = logging.MustGetLogger("mssscheduler")

// MarathonScheduler holds Marathon API URL and a Backoff struct for each task.
type MarathonScheduler struct {
	baseMarathonURL string
	demandUpdate    chan struct{}
	backoff         *utils.Backoff
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

type marathonError struct {
	e       error
	blocked bool
}

func (e marathonError) Error() string {
	if e.blocked {
		return "Deployment blocked"
	}

	return e.e.Error()
}

// NewScheduler returns a pointer to the scheduler.
func NewScheduler(marathonAPI string, demandUpdate chan struct{}) *MarathonScheduler {
	return &MarathonScheduler{
		baseMarathonURL: getBaseMarathonURL(marathonAPI),
		demandUpdate:    demandUpdate,
		backoff: &utils.Backoff{
			Min:    250 * time.Millisecond,
			Max:    5 * time.Second,
			Factor: 2,
		},
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
	return err
}

// StopStartTasks by calling the Marathon scaling API.
func (m *MarathonScheduler) StopStartTasks(tasks *demand.Tasks) error {
	// Create tasks if there aren't enough of them, and stop them if there are too many
	var tooMany []*demand.Task
	var tooFew []*demand.Task
	var err error

	// Check we're not already backed off. This could easily happen if we get a demand update arrive while we are in the midst
	// of a previous backoff.
	if m.backoff.Waiting() {
		log.Debug("Backoff timer still running")
		return nil
	}

	tasks.Lock()
	defer tasks.Unlock()

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

	// Concatentate the two lists - scale down first to free up resources
	tasksToScale := append(tooMany, tooFew...)
	for _, task := range tasksToScale {
		err = m.stopStartTask(task)
		if err != nil {
			if merr, ok := err.(marathonError); ok && merr.blocked {
				// Marathon can't make scale changes at the moment.
				// Trigger a new scaling operation by signalling a demandUpdate after a backoff delay
				err = m.backoff.Backoff(m.demandUpdate)
			} else {
				log.Errorf("Couldn't scale %s: %v ", task.Name, err)
			}
			return err
		}

		// Clear any backoffs on success
		m.backoff.Reset()
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
func (m *MarathonScheduler) stopStartTask(task *demand.Task) (err error) {

	// Scale app using the Marathon REST API.
	status, err := updateApp(m.baseMarathonURL, task.Name, task.Demand)
	if err != nil {
		return err
	}

	switch status {
	case 200:
		// Update was successful
		task.Requested = task.Demand
	case 409:
		// Deployment is locked and we need to back off
		log.Debugf("Deployment locked")
		err = &marathonError{err, true}
	default:
		err = fmt.Errorf("Error response code %d from Marathon API", status)
	}

	return err
}

// Submit a post request to Marathon to match the requested number of the requested app
// format looks like:
// PUT http://marathon:8080/v2/apps/<app>
//  Request:
//  {
//    "instances": 8
//  }
func updateApp(marathonURL string, taskName string, demand int) (status int, err error) {
	url := marathonURL + "apps/" + taskName
	log.Debugf("Start/stop PUT: %s", url)

	payload := startStopPayload{
		Instances: demand,
	}
	w := &bytes.Buffer{}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&payload)
	if err != nil {
		log.Errorf("Failed to encode json. %v", err)
		return 0, err
	}

	// Make scaling call to the Marathon API.
	return utils.PutJSON(url, w)
}

// getBaseMarathonURL returns the base API path.
func getBaseMarathonURL(marathonAPI string) string {
	return marathonAPI + "/v2/"
}

func (m *MarathonScheduler) Cleanup() error {
	m.backoff.Stop()
	return nil
}
