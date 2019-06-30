package azuredevops

// ProjectsService handles communication with the project methods on the API
// utilising https://docs.microsoft.com/en-us/rest/api/azure/devops/core/projects
type ProjectsService struct {
	client *DevOpsClient
}

// ProjectsListResponse describes what the list API call returns for Projects
type ProjectsListResponse struct {
	Projects []Project `json:"value"`
	Count    int       `json:"count,omitempty"`
}

// Project describes a project
type Project struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	URL            string `json:"url"`
	State          string `json:"state"`
	Revision       int    `json:"revision"`
	Visibility     string `json:"visibility"`
	LastUpdateTime string `json:"lastUpdateTime"`
}

// List returns a list of the projects
func (s *ProjectsService) List() ([]Project, error) {
	URL := "/_apis/projects?api-version=5.0"

	request, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	var response ProjectsListResponse
	_, err = s.client.Execute(request, &response)

	return response.Projects, err
}
