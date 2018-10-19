package azuredevops

// IdentityRef represents a Azure Devops user
type IdentityRef struct {
	Descriptor     string `json:"descriptor"`
	DirectoryAlias string `json:"directoryAlias"`
	DisplayName    string `json:"displayName"`
	ID             string `json:"id"`
	ImageURL       string `json:"imageUrl"`
	Inactive       bool   `json:"inactive"`
	IsAadIdentity  bool   `json:"isAadIdentity"`
	IsContainer    bool   `json:"isContainer"`
	ProfileURL     string `json:"profileUrl"`
	UniqueName     string `json:"uniqueName"`
	URL            string `json:"url"`
}
