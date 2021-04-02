package main

import (
	"cqrs/entities/user"
	"cqrs/shared"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_Create(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	cmd := `{"command_name":"create_user", "ExternalId": "eid", "Name": "my name", "Email": "test@example.com"}`

	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(shared.CommandSignalName, cmd)
	}, time.Minute)

	env.RegisterActivity(user.EventReducer)
	env.RegisterActivity(user.CommandHandler)
	env.ExecuteWorkflow(user.UserWorkflow, nil, nil)

	events_encoded, err := env.QueryWorkflow("current_events")
	require.Nil(t, err)

	var events []user.Event

	err = events_encoded.Get(&events)
	require.Nil(t, err)
	require.Equal(t, events[0], user.Event{
		Name: "user_created",
		Meta: map[string]interface{}{
			"Email":      "test@example.com",
			"Name":       "my name",
			"ExternalId": "eid",
		},
	})

	user_encoded, err := env.QueryWorkflow("current_state")
	require.Nil(t, err)

	var state user.UserProjection

	err = user_encoded.Get(&state)
	require.Nil(t, err)
	require.Equal(t, state, user.UserProjection{
		ExternalId: "eid",
		Name:       "my name",
		Email:      "test@example.com",
	})
}

func Test_Token(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	cmd_create := `{"command_name":"create_user", "ExternalId": "eid", "Name": "my name", "Email": "test@example.com"}`
	cmd_token := `{"command_name":"create_token", "Uuid": "a"}`

	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(shared.CommandSignalName, cmd_create)
		env.SignalWorkflow(shared.CommandSignalName, cmd_token)
	}, time.Minute)

	env.RegisterActivity(user.EventReducer)
	env.RegisterActivity(user.CommandHandler)
	env.ExecuteWorkflow(user.UserWorkflow, nil, nil)

	events_encoded, err := env.QueryWorkflow("current_events")
	require.Nil(t, err)

	var events []user.Event

	err = events_encoded.Get(&events)
	require.Nil(t, err)
	require.Equal(t, events[0], user.Event{
		Name: "user_created",
		Meta: map[string]interface{}{
			"Email":      "test@example.com",
			"Name":       "my name",
			"ExternalId": "eid",
		},
	})
	// require.Contains(t, events, cmd_token)

	user_encoded, err := env.QueryWorkflow("current_state")
	require.Nil(t, err)

	var state user.UserProjection

	err = user_encoded.Get(&state)
	require.Nil(t, err)

	require.NotNil(t, state.Tokens[0])

	state.Tokens[0] = "fixed_for_testing"

	require.Equal(t, state, user.UserProjection{
		ExternalId: "eid",
		Name:       "my name",
		Email:      "test@example.com",
		Tokens:     []string{"fixed_for_testing"},
	})
}
