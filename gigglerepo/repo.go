package gigglerepo

type RepoType int

const (
	TYPE_DROPBOX  RepoType = iota
	TYPE_OVERLEAF          = iota
	TYPE_GITHUB            = iota
)

type GiggleRepo struct {
	RType RepoType `json:"rtype"`
}
