package metric

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/storage"
)

// compile-time assert that we implement the right interface
var _ Metric = (*AzureQueueMetric)(nil)

var azureAccountName string
var azureQueueClient storage.QueueServiceClient
var azureClient storage.Client
var azureInitialized bool = false

type AzureQueueMetric struct {
	currentVal     int
	azureQueueName string
}

func AcsInit() (err error) {
	azureAccountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	azureClient, err = storage.NewBasicClient(azureAccountName, os.Getenv("AZURE_STORAGE_ACCOUNT_KEY"))
	if err == nil {
		azureQueueClient = azureClient.GetQueueService()
	}

	return
}

func NewAzureQueueMetric(queueName string) *AzureQueueMetric {
	if !azureInitialized {
		AcsInit()
	}

	return &AzureQueueMetric{
		azureQueueName: queueName,
	}
}

func (aqm *AzureQueueMetric) UpdateCurrent() {
	metadata, err := azureQueueClient.GetMetadata(aqm.azureQueueName)
	if err != nil {
		log.Errorf("Error getting Azure queue info: %v", err)
	}
	aqm.currentVal = metadata.ApproximateMessageCount
}

func (aqm *AzureQueueMetric) Current() int {
	return aqm.currentVal
}
