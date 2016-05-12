// Package marathon provides a scheduler using the Marathon REST API.
package marathon

import (
	"bytes"
	"encoding/json"
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

// NewScheduler returns a pointer to the scheduler.
func NewScheduler(marathonAPI string) *MarathonScheduler {
	return &MarathonScheduler{
		baseMarathonURL: getBaseMarathonURL(marathonAPI),
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

	m.backoff = &utils.Backoff{
		Min:    250 * time.Millisecond,
		Max:    5 * time.Second,
		Factor: 2,
	}
	return err
}

// StopStartTasks by calling the Marathon scaling API.
func (m *MarathonScheduler) StopStartTasks(tasks *demand.Tasks) error {
	// Create tasks if there aren't enough of them, and stop them if there are too many
	var tooMany []*demand.Task
	var tooFew []*demand.Task
	var err error

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
	var (
		err    error
		status int
	)

	// Get the backoff.
	b := m.backoff

	// Function called by AfterFunc once the backoff duration has expired.
	f := func() {
		log.Debugf("Cleared wait for task %s", task.Name)
		b.Clear()
	}

	// Keep attempting to scale app until successful.
	for task.Requested != task.Demand {
		// Don't scale the app while waiting for the backoff duration.
		if b.Waiting() == true {
			log.Debugf("Waiting 250ms for task %s", task.Name)
			time.Sleep(250 * time.Millisecond)

		} else {
			// Scale app using the Marathon REST API.
			status, err = updateApp(m.baseMarathonURL, task.Name, task.Demand)

			if status >= 200 && status <= 299 {
				// Update was successful so set the requested count.
				task.Requested = task.Demand

				// Reset attempts since the deployment is no longer locked.
				if b.Attempt() > 0 {
					log.Infof("Task: %s succeeded set attempts to 0", task.Name)
					b.Reset()
				}
			} else {
				// Clear waiting flag and break loop if the maximum attempts is reached.
				if b.MaxAttempts() {
					log.Debugf("Max attempts reached for task %s", task.Name)
					b.Clear()
					break

				} else {
					// Update failed so back off and attempt again.
					dur := b.Duration(b.Attempt())
					log.Infof("Task: %s failed to scale backing off for %s", task.Name, dur.String())

					timer := time.AfterFunc(dur, f)
					defer timer.Stop()
				}
			}
		}
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
