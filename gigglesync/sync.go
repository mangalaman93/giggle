package gigglesync

import (
	"github.com/mangalaman93/giggle/gigglerepo"
)

type GiggleSyncRepo struct {
	Name       string                   `json:"name"`
	SourceRepo *gigglerepo.GiggleRepo   `json:"source_repo"`
	DestRepos  []*gigglerepo.GiggleRepo `json:"dest_repos"`
}
