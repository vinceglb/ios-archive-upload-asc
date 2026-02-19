package wizard

type Inputs struct {
	Workspace        string
	Scheme           string
	BundleID         string
	TeamID           string
	AppID            string
	AppName          string // display only
	ASCKeyID         string
	ASCIssuerID      string
	ASCPrivateKeyB64 string
	GitHubRepo       string // "owner/repo"
	SecretsWereSet   bool
	WorkflowWasWritten bool
	WorkflowPath     string
}
