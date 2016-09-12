package gigglerepo

func NewOverleafRepo() *GiggleRepo {
  return &GiggleRepo{
    RType: TYPE_OVERLEAF,
  }
}

func (gs *GiggleRepo) pullOverleaf(localFolder, remoteFolder string) error {
  return nil
}
