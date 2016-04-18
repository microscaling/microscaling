// Package metric implements the queue metric for NSQ (http://nsq.io/).
package metric

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const constNSQStatsEndpoint string = "127.0.0.1:4151"
const constNSQStatsAPI string = "/stats?format=json"

// compile-time assert that we implement the right interface
var _ Metric = (*NSQMetric)(nil)

// NSQMetric stores the current value.
type NSQMetric struct {
	currentVal int
	queueName  string
}

// StatsMessage from NSQ stats API.
type StatsMessage struct {
	Data StatsData `json:"data"`
}

// StatsData from NSQ stats API.
type StatsData struct {
	Topics []Topic `json:"topics"`
}

// Topic from NSQ stats API.
type Topic struct {
	TopicName string    `json:"topic_name"`
	Channels  []Channel `json:"channels"`
}

// Channel from NSQ stats API.
type Channel struct {
	ChannelName string `json:"channel_name"`
	Depth       int    `json:"depth"`
}

var (
	httpClient = &http.Client{
		// TODO Make timeout configurable.
		Timeout: 10 * time.Second,
	}
	nsqStatsEndpoint string
	nsqInitialized   = false
)

// NSQInit sets up the NSQ Stats endpoint.
func NSQInit() {
	nsqStatsEndpoint = os.Getenv("NSQ_STATS_ENDPOINT")
	if nsqStatsEndpoint == "" {
		nsqStatsEndpoint = constNSQStatsEndpoint
	}

	nsqInitialized = true
}

// NewNSQMetric creates the metric.
func NewNSQMetric(queueName string) *NSQMetric {
	if !nsqInitialized {
		NSQInit()
	}

	return &NSQMetric{
		queueName: queueName,
	}
}

// Call NSQ Stats API to get the queue length.
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

// UpdateCurrent sets the current queue length.
func (nsqm *NSQMetric) UpdateCurrent() {
	var statsMessage StatsMessage

	url := "http://" + nsqStatsEndpoint + constNSQStatsAPI
	body, err := getJSONGet(url)
	if err != nil {
		log.Errorf("Error getting NSQ metric %v", err)
	}

	err = json.Unmarshal(body, &statsMessage)
	if err != nil {
		log.Errorf("Error %v unmarshalling from %s", err, string(body[:]))
	}

	// Loop through NSQ Channels and Metrics to find the correct value.
	// Currently queue name is used for both the Topic and Channel.
	// TODO May need to split these later.
	for _, topic := range statsMessage.Data.Topics {
		if topic.TopicName == nsqm.queueName {
			for _, channel := range topic.Channels {
				if channel.ChannelName == nsqm.queueName {
					nsqm.currentVal = channel.Depth
				}
			}
		}
	}

	log.Debugf("Queue name %s length %d", nsqm.queueName, nsqm.currentVal)
}

// Current returns the queue length.
func (nsqm *NSQMetric) Current() int {
	return nsqm.currentVal
}
