// Package marathon scheduler integration
package marathon

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/op/go-logging"

	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/scheduler"
)

var log = logging.MustGetLogger("mssscheduler")

// MarathonScheduler holds Marathon API URL.
type MarathonScheduler struct {
	baseMarathonURL string
}

// AppsMessage from the Marathon Apps API.
type AppsMessage struct {
	Apps []App `json:"apps"`
}

// App from the Marathon Apps API.
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
	return err
}

// StopStartTasks by calling the Marathon scaling API.
func (m *MarathonScheduler) StopStartTasks(tasks *demand.Tasks) error {
	// Create tasks if there aren't enough of them, and stop them if there are too many
	var tooMany []*demand.Task
	var tooFew []*demand.Task
	var diff int
	var err error = nil

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
		diff = task.Requested - task.Demand
		log.Debugf("Stop %d of task %s", diff, task.Name)
		err = m.stopStartTask(task)
		if err != nil {
			log.Errorf("Couldn't stop %s: %v ", task.Name, err)
		}
		log.Infof("now have %d", task.Requested)
	}

	// Now we can scale up
	for _, task := range tooFew {
		diff = task.Demand - task.Requested
		log.Debugf("Start %d of task %s", diff, task.Name)
		err = m.stopStartTask(task)
		if err != nil {
			log.Errorf("Couldn't start %s: %v ", task.Name, err)
		}
		log.Infof("now have %d", task.Requested)
	}

	log.Infof("%v", tasks)
	return err
}

// Marathon API to get the current running tasks.
func getJSONGet(url string) (body []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Failed to build API GET request err %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Errorf("Failed to GET err %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Http error %d: %s", resp.StatusCode, resp.Status)
	}

	body, err = ioutil.ReadAll(resp.Body)

	return body, err
}

// CountAllTasks tells us how many instances of each task are currently running.
func (m *MarathonScheduler) CountAllTasks(running *demand.Tasks) error {
	var (
		err         error
		appsMessage AppsMessage
	)

	url := m.baseMarathonURL + "apps/"

	body, err := getJSONGet(url)
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
	log.Infof("Start/stop PUT: %s", url)

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

	req, err := http.NewRequest("PUT", url, w)
	if err != nil {
		log.Errorf("Failed to build PUT request err %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Errorf("start/stop err %v", err)
		return err
	}

	if resp.StatusCode > 299 {
		log.Errorf("error response from marathon. %s", resp.Status)
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("start/stop read err %v", err)
		return err
	}

	// We do nothing with this body
	s := string(body)
	log.Infof("start/stop json: %s", s)

	// Now we've asked for this many, update the currentcount
	task.Requested = task.Demand

	return err
}

// getBaseMarathonURL returns the base API path.
func getBaseMarathonURL(marathonAPI string) string {
	return marathonAPI + "/v2/"
}
