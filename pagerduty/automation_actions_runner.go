package pagerduty

import "fmt"

// AutomationActionsRunner handles the communication with schedule
// related methods of the PagerDuty API.
type AutomationActionsRunnerService service

type AutomationActionsRunner struct {
	ID           string                       `json:"id"`
	Name         string                       `json:"name"`
	Summary      string                       `json:"summary,omitempty"`
	Type         string                       `json:"type"`
	Description  string                       `json:"description,omitempty"`
	CreationTime string                       `json:"creation_time"`
	RunnerType   string                       `json:"runner_type,omitempty"`
	Status       string                       `json:"status"`
	Teams        []*TeamReference             `json:"teams,omitempty"`
	Privileges   *AutomationActionsPriviliges `json:"privileges,omitempty"`
}

type AutomationActionsPriviliges struct {
	Permissions []*string `json:"permissions,omitempty"`
}

type AutomationActionsRunnerPayload struct {
	Runner *AutomationActionsRunner `json:"runner,omitempty"`
}

// Create creates a new runner
func (s *AutomationActionsRunnerService) Create(runner *AutomationActionsRunner) (*AutomationActionsRunner, *Response, error) {
	u := "/automation_actions/runners"
	v := new(AutomationActionsRunnerPayload)
	o := RequestOptions{
		Type:  "header",
		Label: "X-EARLY-ACCESS",
		Value: "automation-actions-early-access",
	}

	resp, err := s.client.newRequestDoOptions("POST", u, nil, &AutomationActionsRunnerPayload{Runner: runner}, &v, o)
	if err != nil {
		return nil, nil, err
	}

	return v.Runner, resp, nil
}

// Get retrieves information about a runner.
func (s *AutomationActionsRunnerService) Get(id string) (*AutomationActionsRunner, *Response, error) {
	u := fmt.Sprintf("/automation_actions/runners/%s", id)
	v := new(AutomationActionsRunnerPayload)
	o := RequestOptions{
		Type:  "header",
		Label: "X-EARLY-ACCESS",
		Value: "automation-actions-early-access",
	}

	resp, err := s.client.newRequestDoOptions("GET", u, nil, nil, &v, o)
	if err != nil {
		return nil, nil, err
	}

	return v.Runner, resp, nil
}

// Delete deletes an existing runner.
func (s *AutomationActionsRunnerService) Delete(id string) (*Response, error) {
	u := fmt.Sprintf("/automation_actions/runners/%s", id)
	o := RequestOptions{
		Type:  "header",
		Label: "X-EARLY-ACCESS",
		Value: "automation-actions-early-access",
	}

	return s.client.newRequestDoOptions("DELETE", u, nil, nil, nil, o)
}
