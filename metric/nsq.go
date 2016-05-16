// Package metric implements the queue metric for NSQ (http://nsq.io/).
package metric

import (
	"encoding/json"
	"os"

	"github.com/microscaling/microscaling/utils"
)

const constNSQStatsEndpoint string = "127.0.0.1:4151"
const constNSQStatsAPI string = "/stats?format=json"

// compile-time assert that we implement the right interface
var _ Metric = (*NSQMetric)(nil)

// NSQMetric stores the current value.
type NSQMetric struct {
	currentVal  int
	topicName   string
	channelName string
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
	return
}

// NewNSQMetric creates the metric.
func NewNSQMetric(topicName string, channelName string) *NSQMetric {
	if !nsqInitialized {
		NSQInit()
	}

	return &NSQMetric{
		topicName:   topicName,
		channelName: channelName,
	}
}

// UpdateCurrent sets the current queue length.
func (nsqm *NSQMetric) UpdateCurrent() {
	var statsMessage StatsMessage

	url := "http://" + nsqStatsEndpoint + constNSQStatsAPI
	body, err := utils.GetJSON(url)
	if err != nil {
		log.Errorf("Error getting NSQ metric %v", err)
	}

	err = json.Unmarshal(body, &statsMessage)
	if err != nil {
		log.Errorf("Error %v unmarshalling from %s", err, string(body[:]))
	}

	// Loop through NSQ Channels and Metrics to find the correct value.
	for _, topic := range statsMessage.Data.Topics {
		if topic.TopicName == nsqm.topicName {
			for _, channel := range topic.Channels {
				if channel.ChannelName == nsqm.channelName {
					nsqm.currentVal = channel.Depth
				}
			}
		}
	}

	log.Debugf("Topic: %s Channel: %s Length: %d", nsqm.topicName, nsqm.channelName, nsqm.currentVal)
}

// Current returns the queue length.
func (nsqm *NSQMetric) Current() int {
	return nsqm.currentVal
}
