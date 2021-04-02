package user

import (
	"cqrs/shared"
	"time"

	"go.temporal.io/sdk/workflow"
)

const TaskQueueName = "USER_TASK_QUEUE"

func UserWorkflow(ctx workflow.Context, oldEvts []string, user UserProjection) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	signalChan := workflow.GetSignalChannel(ctx, shared.CommandSignalName)
	selector := workflow.NewSelector(ctx)

	encodedEvents := workflow.SideEffect(ctx, func(_ workflow.Context) interface{} {
		return oldEvts
	})

	var events []Event
	encodedEvents.Get(&events)

	if err := workflow.SetQueryHandler(ctx, "current_events", func() ([]Event, error) {
		return events, nil
	}); err != nil {
		return err
	}

	if err := workflow.SetQueryHandler(ctx, "current_state", func() (UserProjection, error) {
		return user, nil
	}); err != nil {
		return err
	}

	selector.AddReceive(signalChan, func(ch workflow.ReceiveChannel, more bool) {
		var signalValue string
		ch.Receive(ctx, &signalValue)

		var newEvents []Event

		if err := workflow.ExecuteActivity(ctx, CommandHandler, signalValue, events).Get(ctx, &newEvents); err != nil {
			panic(err)
		}
		events = append(events, newEvents...)

		for _, newEvent := range newEvents {
			if err := workflow.ExecuteActivity(ctx, EventReducer, newEvent, user).Get(ctx, &user); err != nil {
				panic(err)
			}
		}
	})

	for {
		selector.Select(ctx)
	}
}
