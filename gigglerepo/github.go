package gigglerepo

func NewGithubRepo() *GiggleRepo {
  return &GiggleRepo{
    RType: TYPE_GITHUB,
  }
}

func (gs *GiggleRepo) syncGithub(localFolder, remoteFolder string) error {
  return nil
}
