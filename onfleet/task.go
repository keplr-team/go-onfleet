package onfleet

import (
	"context"
	"errors"

	"github.com/google/go-querystring/query"
)

type TasksService service

type Task struct {
	Id                string            `json:"id"`
	TimeCreated       int64             `json:"timeCreated"`
	TimeLastModified  int64             `json:"timeLastModified"`
	Organization      string            `json:"organization"`
	ShortID           string            `json:"shortId"`
	TrackingURL       string            `json:"trackingURL"`
	Worker            string            `json:"worker"`
	Merchant          string            `json:"merchant"`
	Executor          string            `json:"executor"`
	Creator           string            `json:"creator"`
	Dependencies      []interface{}     `json:"dependencies"`
	State             int               `json:"state"`
	CompleteAfter     int64             `json:"completeAfter"`
	CompleteBefore    interface{}       `json:"completeBefore"`
	PickupTask        bool              `json:"pickupTask"`
	Notes             string            `json:"notes"`
	CompletionDetails CompletionDetails `json:"completionDetails"`
	Feedback          []interface{}     `json:"feedback"`
	Metadata          []interface{}     `json:"metadata"`
	Overrides         Overrides         `json:"overrides"`
	Container         Container         `json:"container"`
	Recipients        []Recipients      `json:"recipients"`
	Destination       Destination       `json:"destination"`
	DidAutoAssign     bool              `json:"didAutoAssign"`
}

type CompletionDetails struct {
	Events []interface{} `json:"events"`
	Time   interface{}   `json:"time"`
}

type Overrides struct {
	RecipientSkipSMSNotifications interface{} `json:"recipientSkipSMSNotifications"`
	RecipientNotes                interface{} `json:"recipientNotes"`
	RecipientName                 interface{} `json:"recipientName"`
}

type Container struct {
	Type   string `json:"type"`
	Worker string `json:"worker"`
}

type Recipients struct {
	Id                   string        `json:"id"`
	Organization         string        `json:"organization"`
	TimeCreated          int64         `json:"timeCreated"`
	TimeLastModified     int64         `json:"timeLastModified"`
	Name                 string        `json:"name"`
	Phone                string        `json:"phone"`
	Notes                string        `json:"notes"`
	SkipSMSNotifications bool          `json:"skipSMSNotifications"`
	Metadata             []interface{} `json:"metadata"`
}

type Destination struct {
	Id               string        `json:"id"`
	TimeCreated      int64         `json:"timeCreated"`
	TimeLastModified int64         `json:"timeLastModified"`
	Address          Address       `json:"address"`
	Notes            string        `json:"notes"`
	Metadata         []interface{} `json:"metadata"`
}

type TasksListState string

var (
	TasksListStateUnassigned TasksListState = "0"
	TasksListStateAssigned   TasksListState = "1"
	TasksListStateActive     TasksListState = "2"
	TasksListStateCompleted  TasksListState = "3"
)

type TasksListOptions struct {
	From   int64            `url:"from"`
	States []TasksListState `url:"state,omitempty"`
	Worker string           `url:"worker,omitempty"`
}

type TaskPayload struct {
	Destination   Destination  `json:"destination"`
	Recipients    []Recipients `json:"recipients"`
	CompleteAfter int64        `json:"completeAfter"`
	Notes         string       `json:"notes"`
	Container     Container    `json:"container"`
}

type TasksCreatePayload struct {
	Tasks []TaskPayload `json:"tasks"`
}

type TaskError struct {
	StatusCode int    `json:"statusCode"`
	Error      int    `json:"error"`
	Message    string `json:"message"`
	Cause      string `json:"cause"`
}

type TasksCreateReturn struct {
	Tasks  []Task `json:"tasks"`
	Errors []struct {
		Error TaskError   `json:"error"`
		Task  TaskPayload `json:"task"`
	} `json:"errors"`
}

// List all tasks
// https://docs.onfleet.com/reference#list-tasks
func (s *TasksService) List(ctx context.Context, opts *TasksListOptions) ([]Task, error) {
	var tasks []Task
	err := s.getMany(ctx, "tasks", opts, &tasks)
	return tasks, err
}

func (s *TasksService) getMany(ctx context.Context, path string, opts interface{}, v interface{}) error {
	if opts != nil {
		v, err := query.Values(opts)
		if err != nil {
			return err
		}
		path += "?" + v.Encode()
	}

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}

	err = s.client.Do(ctx, req, v)
	if err != nil {
		return err
	}

	return nil
}

// Create tasks
// https://docs.onfleet.com/reference#create-tasks-in-batch
func (s *TasksService) Create(ctx context.Context, payload *TasksCreatePayload) ([]Task, error) {
	var res TasksCreateReturn
	req, err := s.client.NewRequest("POST", "tasks/batch", payload)
	if err != nil {
		return nil, err
	}

	err = s.client.Do(ctx, req, &res)
	if err != nil {
		return nil, err
	}
	if len(res.Errors) != 0 {
		return nil, errors.New(res.Errors[0].Error.Message)
	}

	return res.Tasks, nil
}

// Update task
// https://docs.onfleet.com/reference#update-task
func (s *TasksService) Update(ctx context.Context, taskId string, payload *TaskPayload) (*Task, error) {
	var res Task
	req, err := s.client.NewRequest("PUT", "tasks/"+taskId, payload)
	if err != nil {
		return nil, err
	}

	err = s.client.Do(ctx, req, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
