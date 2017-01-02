package metric

import (
	"errors"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

const constQueueLengthAttribute string = "ApproximateNumberOfMessages"

// compiletime assert that we implement the right interface
var _ Metric = (*SQSMetric)(nil)

// SQSMetric is used to measure the length of an SQS Queue
type SQSMetric struct {
	client     sqsiface.SQSAPI
	currentVal int
	queueURL   string
}

// NewSQSMetric makes sure we have access to the SQS client
func NewSQSMetric(queueURL string) (metric *SQSMetric, err error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return nil, errors.New("AWS_REGION env var must be set")
	}

	var config *aws.Config = &aws.Config{Region: aws.String(region)}

	sess, err := session.NewSession(config)
	if err != nil {
		log.Errorf("Failed to create AWS session: %v", err)
		return
	}

	metric = &SQSMetric{
		client:   sqs.New(sess),
		queueURL: queueURL,
	}

	return
}

// UpdateCurrent calls the SQS API to get the queue length and stores the value in the metric.
func (sm *SQSMetric) UpdateCurrent() {
	a := make([]*string, 1)
	a[0] = aws.String(constQueueLengthAttribute)

	params := sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(sm.queueURL),
		AttributeNames: a,
	}

	m, err := sm.client.GetQueueAttributes(&params)
	if err != nil {
		log.Errorf("Failed to get SQS queue info: %v", err)
	}

	v := aws.StringValue(m.Attributes[constQueueLengthAttribute])
	sm.currentVal, err = strconv.Atoi(v)
	if err != nil {
		log.Errorf("Failed to convert queue length to int: %v", err)
	} else {
		log.Debugf("Queue URL %s length %d", sm.queueURL, sm.currentVal)
	}
}

// Current reads out the value of the current queue length
func (sm *SQSMetric) Current() int {
	return sm.currentVal
}
