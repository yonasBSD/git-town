package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

type GitUserName string

func (self GitUserName) String() string {
	return string(self)
}

func NewGitUserNameOption(value string) Option[GitUserName] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[GitUserName]()
	}
	return Some(GitUserName(value))
}
