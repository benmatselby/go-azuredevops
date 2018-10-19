package azuredevops

import (
	"fmt"
)

// BuildsService handles communication with the builds methods on the API
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/build/builds
type BuildsService struct {
	client *Client
}

// BuildsListResponse is the wrapper around the main response for the List of Builds
type BuildsListResponse struct {
	Builds []Build `json:"value"`
}

type buildOrchestrationPlanSchema struct {
	Type   int    `json:"orchestrationType"`
	PlanID string `json:"planId"`
}

// Build represents a build
type Build struct {
	Definition    BuildDefinition `json:"definition"`
	Controller    BuildController `json:"controller"`
	LastChangedBy *IdentityRef    `json:"lastChangedBy,omitempty"`
	DeletedBy     *IdentityRef    `json:"deletedBy,omitempty"`
	BuildNumber   string          `json:"buildNumber,omitempty"`
	FinishTime    string          `json:"finishTime,omitempty"`
	Branch        string          `json:"sourceBranch"`
	Repository    Repository      `json:"repository"`
	Demands       []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"demands"`
	Logs *struct {
		ID   int    `json:"id"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"logs,omitempty"`
	Project *struct {
		Abbreviation string `json:"abbreviation"`
		Description  string `json:"description"`
		ID           string `json:"id"`
		Name         string `json:"name"`
		Revision     string `json:"revision"`
		State        string `json:"state"`
		URL          string `json:"url"`
		Visibility   string `json:"visibility"`
	} `json:"project,omitempty"`
	Properties          map[string]string
	Priority            string                         `json:"priority,omitempty"`
	OrchestrationPlan   *buildOrchestrationPlanSchema  `json:"orchestrationPlan,omitempty"`
	Plans               []buildOrchestrationPlanSchema `json:"plans,omitempty"`
	BuildNumberRevision int                            `json:"buildNumberRevision,omitempty"`
	Deleted             *bool                          `json:"deleted,omitempty"`
	DeletedDate         string                         `json:"deletedDate,omitempty"`
	DeletedReason       string                         `json:"deletedReason,omitempty"`
	ID                  string                         `json:"id,omitempty"`
	KeepForever         string                         `json:"keepForever,omitempty"`
	ChangedDate         string                         `json:"lastChangedDate,omitempty"`
	Params              string                         `json:"parameters,omitempty"`
	Quality             string                         `json:"quality,omitempty"`
	Queue               struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
		Pool *struct {
			ID       int    `json:"id"`
			IsHosted bool   `json:"is_hosted"`
			Name     string `json:"name"`
		} `json:"pool,omitempty"`
	} `json:"queue"`
	QueueOptions      map[string]string `json:"queue_options"`
	QueuePosition     *int              `json:"queuePosition,omitempty"`
	QueueTime         string            `json:"queueTime,omitempty"`
	RetainedByRelease *bool             `json:"retainedByRelease,omitempty"`
	Version           string            `json:"sourceVersion,omitempty"`
	StartTime         string            `json:"startTime,omitempty"`
	Status            string            `json:"status,omitempty"`
	Result            string            `json:"result,omitempty"`
	ValidationResults []struct {
		Message string `json:"message"`
		Result  string `json:"result"`
	}
	Tags         []string `json:"tags,omitempty"`
	TriggerBuild *Build   `json:"triggeredByBuild,omitempty"`
	URI          string   `json:"uri,omitempty"`
	URL          string   `json:"url,omitempty"`
}

// BuildsListOptions describes what the request to the API should look like
type BuildsListOptions struct {
	Definitions      string `url:"definitions,omitempty"`
	Branch           string `url:"branchName,omitempty"`
	Count            int    `url:"$top,omitempty"`
	Repository       string `url:"repositoryId,omitempty"`
	BuildIDs         string `url:"buildIds,omitempty"`
	Order            string `url:"queryOrder,omitempty"`
	Deleted          string `url:"deletedFilter,omitempty"`
	MaxPerDefinition string `url:"maxBuildsPerDefinition,omitempty"`
	Token            string `url:"continuationToken,omitempty"`
	Props            string `url:"properties,omitempty"`
	Tags             string `url:"tagFilters,omitempty"`
	Result           string `url:"resultFilter,omitempty"`
	Status           string `url:"statusFilter,omitempty"`
	Reason           string `url:"reasonFilter,omitempty"`
	UserID           string `url:"requestedFor,omitempty"`
	MaxTime          string `url:"maxTime,omitempty"`
	MinTime          string `url:"minTime,omitempty"`
	BuildNumber      string `url:"buildNumber,omitempty"`
	Queues           string `url:"queues,omitempty"`
	RepoType         string `url:"repositoryType,omitempty"`
}

// List returns list of the builds
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/build/builds/list
func (s *BuildsService) List(opts *BuildsListOptions) ([]Build, error) {
	URL := fmt.Sprintf("/_apis/build/builds?api-version=4.1")
	URL, err := addOptions(URL, opts)

	request, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	var response BuildsListResponse
	_, err = s.client.Execute(request, &response)

	return response.Builds, err
}

// QueueBuildOptions describes what the request to the API should look like
type QueueBuildOptions struct {
	IgnoreWarnings bool   `url:"ignoreWarnings,omitempty"`
	CheckInTicket  string `url:"checkInTicket,omitempty"`
}

// Queue inserts new build creation to queue
// utilising https://docs.microsoft.com/en-us/rest/api/vsts/build/builds/queue?view=vsts-rest-4.1
func (s *BuildsService) Queue(build *Build, opts *QueueBuildOptions) error {
	URL := "_apis/build/builds?api-version=4.1"
	URL, err := addOptions(URL, opts)

	if err != nil {
		return err
	}

	request, err := s.client.NewRequest("POST", URL, build)

	if err != nil {
		return err
	}

	_, err = s.client.Execute(request, &build)

	return err
}
