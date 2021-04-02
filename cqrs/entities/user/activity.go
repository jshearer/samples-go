package user

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

func CommandHandler(cmd string, events []Event) ([]Event, error) {
	var result Command

	if err := json.Unmarshal([]byte(cmd), &result); err != nil {
		return nil, err
	}

	switch result.Name {
	case "create_user":
		{
			var command CreateUser

			if err := json.Unmarshal([]byte(cmd), &command); err != nil {
				return nil, err
			} else {
				for _, v := range events {
					if v.Name == "user_created" {
						return nil, fmt.Errorf("invalid 'create_user': User already exists")
					}
				}
				return []Event{{
					Name: "user_created",
					Meta: UserCreated{
						ExternalId: command.ExternalId,
						Email:      command.Email,
						Name:       command.Name,
					},
				}}, nil
			}
		}
	case "create_token":
		{
			var command CreateToken

			if err := json.Unmarshal([]byte(cmd), &command); err != nil {
				return nil, err
			} else {
				return []Event{{
					Name: "token_created",
					Meta: TokenCreated{
						Uuid: uuid.NewString(),
					},
				}}, nil
			}
		}
	case "upload_spec":
		{
			var command UploadSpec

			if err := json.Unmarshal([]byte(cmd), &command); err != nil {
				return nil, err
			} else {
				return []Event{{
					Name: "token_created",
					Meta: command.Meta,
				}}, nil
			}
		}
	}

	return nil, fmt.Errorf("unknown command: %v", result.Name)
}

func EventReducer(evt Event, state UserProjection) (UserProjection, error) {
	switch evt.Name {
	case "user_created":
		{
			var event UserCreated
			err := mapstructure.Decode(evt.Meta, &event)

			if err != nil {
				return state, fmt.Errorf("invalid user_created event: %+v", evt.Meta)
			} else if state.Email == "" && state.ExternalId == "" && state.Name == "" {
				state.Email = event.Email
				state.ExternalId = event.ExternalId
				state.Name = event.Name

				return state, nil
			} else {
				return state, fmt.Errorf("invalid 'create_user': User already exists")
			}
		}
	case "token_created":
		{
			var event TokenCreated
			err := mapstructure.Decode(evt.Meta, &event)

			if err != nil {
				return state, fmt.Errorf("invalid token_created event: %+v", evt.Meta)
			} else {
				state.Tokens = append(state.Tokens, event.Uuid)
				return state, nil
			}
		}
	case "spec_uploaded":
		{
			var event SpecUploaded
			err := mapstructure.Decode(evt.Meta, &event)

			if err != nil {
				return state, fmt.Errorf("invalid spec_uploaded event: %+v", evt.Meta)
			} else {
				state.Specs = append(state.Specs, event.Meta)
				return state, nil
			}
		}
	}

	return state, fmt.Errorf("unknown command: %v", evt.Name)
}
