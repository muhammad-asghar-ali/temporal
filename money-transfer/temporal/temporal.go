package temporal

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"money-transfer/interenal/bank"
	"money-transfer/shared"
)

type Temporal struct{}

func (t *Temporal) Client() (client.Client, error) {
	// Implement the method to create and return the Temporal client
	c, err := client.Dial(client.Options{})
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (t *Temporal) Worker(queueName string) worker.Worker {
	// Implement the method to create and return the Temporal worker for the specified queue
	c, _ := t.Client()
	wf := worker.New(c, queueName, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	wf.RegisterWorkflow(bank.MoneyTransfer)
	wf.RegisterActivity(bank.Withdraw)
	wf.RegisterActivity(bank.Deposit)
	wf.RegisterActivity(bank.Refund)

	return wf
}

func RunTemporalWorker(t *Temporal) error {
	c, err := t.Client()
	if err != nil {
		return err
	}
	defer c.Close()

	// Get the Temporal worker from the Temporal instance
	w := t.Worker(shared.MoneyTransferTaskQueueName)

	// Start listening to the Task Queue.
	if err := w.Run(worker.InterruptCh()); err != nil {
		return err
	}

	return nil
}
