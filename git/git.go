package git

// GitRepo is the basic store object.
type GitRepo struct {
	repo   string
	folder string
}

// New returns a new GPGStore that can then needs to be initialized with Init()
func New(repo, folder string) (*Repo, error) {
	gr := new(GitRepo)
	gr.repo = repo
	gr.folder = folder
	return gr, nil
}

// Update will clone a repo if it doesn't exist or pull a repo, if it does.
func (gr *GitRepo) Update() error {
}
