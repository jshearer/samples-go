package user

import (
	"encoding/json"
	"fmt"
)

// Commands
type Command struct {
	Name string `json:"command_name"`
	Meta interface{}
}

type CreateUser struct {
	ExternalId string
	Name       string
	Email      string
}

type CreateToken struct {
}

type UploadSpec struct {
	Meta SpecMeta
}

func (e Command) parse(evt string) (interface{}, error) {
	var result Command

	if err := json.Unmarshal([]byte(evt), &result); err != nil {
		return nil, err
	}

	switch result.Name {
	case "create_user":
		{
			var parsed CreateUser
			if err := json.Unmarshal([]byte(evt), &evt); err != nil {
				return nil, err
			} else {
				return parsed, nil
			}
		}
	case "create_token":
		{
			var parsed CreateToken
			if err := json.Unmarshal([]byte(evt), &evt); err != nil {
				return nil, err
			} else {
				return parsed, nil
			}
		}
	case "upload_spec":
		{
			var parsed UploadSpec
			if err := json.Unmarshal([]byte(evt), &evt); err != nil {
				return nil, err
			} else {
				return parsed, nil
			}
		}
	}

	return nil, fmt.Errorf("unknown event type: %v", result.Name)
}

// Events
type Event struct {
	Name string `json:"event_name"`
	Meta interface{}
}

type UserCreated struct {
	ExternalId string
	Name       string
	Email      string
}

type TokenCreated struct {
	Uuid string
}

type SpecUploaded struct {
	Meta SpecMeta
}

func (e Event) parse(evt string) (interface{}, error) {
	var result Event

	if err := json.Unmarshal([]byte(evt), &result); err != nil {
		return nil, err
	}

	switch result.Name {
	case "user_created":
		{
			var parsed UserCreated
			if err := json.Unmarshal([]byte(evt), &evt); err != nil {
				return nil, err
			} else {
				return parsed, nil
			}
		}
	case "token_created":
		{
			var parsed TokenCreated
			if err := json.Unmarshal([]byte(evt), &evt); err != nil {
				return nil, err
			} else {
				return parsed, nil
			}
		}
	case "spec_uploaded":
		{
			var parsed SpecUploaded
			if err := json.Unmarshal([]byte(evt), &evt); err != nil {
				return nil, err
			} else {
				return parsed, nil
			}
		}
	}

	return nil, fmt.Errorf("unknown event type: %v", result.Name)
}

type SpecMeta struct {
	Uuid             string
	S3Path           string
	GithubUser       string
	GithubOrg        string
	GithubRepo       string
	DatetimeUploaded string
}

type UserProjection struct {
	Uuid       string
	ExternalId string
	Name       string
	Email      string
	Tokens     []string
	Specs      []SpecMeta
}
