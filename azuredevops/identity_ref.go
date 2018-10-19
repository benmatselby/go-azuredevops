package azuredevops

type IdentityRef struct {
	Descriptor     string `json:"descriptor"`
	DirectoryAlias string `json:"directoryAlias"`
	DisplayName    string `json:"displayName"`
	ID             string `json:"id"`
	ImageURL       string `json:"imageUrl"`
	Inactive       bool   `json:"inactive"`
	IsAadIdentity  bool   `json:"isAadIdentity"`
	IsContainer    bool   `json:"isContainer"`
	ProfileUrl     string `json:"profileUrl"`
	UniqueName     string `json:"uniqueName"`
	URL            string `json:"url"`
}
