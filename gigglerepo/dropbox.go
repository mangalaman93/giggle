package gigglerepo

func NewDropboxRepo() *GiggleRepo {
	return &GiggleRepo{
		RType: TYPE_DROPBOX,
	}
}

func (gs *GiggleRepo) syncDropbox(localFolder, remoteFolder string) error {
	return nil
}
