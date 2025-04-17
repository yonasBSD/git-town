package gitlab

import "github.com/git-town/git-town/v19/internal/git/giturl"

// Detect indicates whether the current repository is hosted on a GitLab server.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "gitlab.com"
}
