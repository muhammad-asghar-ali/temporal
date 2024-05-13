package main

import (
	"context"
	"log"
	"sync"

	"go.temporal.io/sdk/client"

	"money-transfer/interenal/bank"
	"money-transfer/shared"
	"money-transfer/temporal"
)

var (
	once sync.Once
	c    client.Client
	err  error
)

func main() {
	// Use sync.Once to ensure the client is created only once per process.
	once.Do(func() {
		t := &temporal.Temporal{}
		c, err = t.Client()
		if err != nil {
			log.Fatalln("Unable to create Temporal client:", err)
		}
	})

	defer c.Close()

	// Call RunTemporalWorker from the temporal package
	if err := temporal.RunTemporalWorker(&temporal.Temporal{}); err != nil {
		log.Fatalln("Unable to start Temporal worker:", err)
	}

	input := shared.PaymentDetails{
		SourceAccount: "85-150",
		TargetAccount: "43-812",
		Amount:        250,
		ReferenceID:   "12345",
	}

	options := client.StartWorkflowOptions{
		ID:        "pay-invoice-701",
		TaskQueue: shared.MoneyTransferTaskQueueName,
	}

	log.Printf("Starting transfer from account %s to account %s for %d", input.SourceAccount, input.TargetAccount, input.Amount)

	we, err := c.ExecuteWorkflow(context.Background(), options, bank.MoneyTransfer, input)
	if err != nil {
		log.Fatalln("Unable to start the Workflow:", err)
	}

	log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

	var result string

	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable to get Workflow result:", err)
	}

	log.Println(result)
}
