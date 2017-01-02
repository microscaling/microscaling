package metric

import (
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type mockedQueueAttributes struct {
	sqsiface.SQSAPI
	Resp sqs.GetQueueAttributesOutput
}

type sqsTest struct {
	Resp     sqs.GetQueueAttributesOutput
	Expected int
}

// Mock SQS API call and just return the response that is parsed
func (m mockedQueueAttributes) GetQueueAttributes(in *sqs.GetQueueAttributesInput) (*sqs.GetQueueAttributesOutput, error) {
	return &m.Resp, nil
}

func TestUpdateCurrent(t *testing.T) {
	queueURL := "https://sqs.us-east-1.amazonaws.com/1234567890/microscaling-test"

	cases := make([]sqsTest, 2)
	cases[0] = sqsTest{
		Resp:     getQueueAttributes(0),
		Expected: 0,
	}
	cases[1] = sqsTest{
		Resp:     getQueueAttributes(42),
		Expected: 42,
	}

	for i, c := range cases {
		m := SQSMetric{
			client:   mockedQueueAttributes{Resp: c.Resp},
			queueURL: queueURL,
		}

		m.UpdateCurrent()

		if m.Current() != c.Expected {
			t.Errorf("Test %d: expected count %d but was %d", i, c.Expected, m.Current)
		}
	}
}

func getQueueAttributes(count int) sqs.GetQueueAttributesOutput {
	a := make(map[string]*string)
	a["ApproximateNumberOfMessages"] = aws.String(strconv.Itoa(count))

	return sqs.GetQueueAttributesOutput{
		Attributes: a,
	}
}
