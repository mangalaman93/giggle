package svc

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/mangalaman93/giggle/conf"
)

func performSync(ctx context.Context, cf *conf.Config) error {
	log.Println("[INFO] periodic sync begin")
	defer log.Println("[INFO] periodic sync end")

	for _, sc := range cf.Sync {
		log.Printf("[INFO] syncing %v\n", sc.Name)
		if err := syncRepo(ctx, sc, cf.Auth); err != nil {
			log.Printf("[WARN] error syncing %v :: %v\n", sc.Name, err)
			continue
		}

		log.Printf("[INFO] synced %v\n", sc.Name)
	}

	return nil
}

func syncRepo(ctx context.Context, sc conf.SyncConfig,
	authMap map[string]*conf.AuthMethod) error {

	repoFolder := conf.GetSyncTarget(sc.Name)
	fromAuth := authMap[sc.From.AuthToUse]
	fromRepo, err := openRepo(ctx, sc.From, fromAuth, repoFolder)
	if err != nil {
		return err
	}

	fromRemote, err := fromRepo.Remote(sc.From.Name)
	if err != nil {
		return err
	}

	if err := fetch(ctx, fromRemote, fromAuth); err != nil {
		return err
	}

	var errRet error
	for _, to := range sc.ToList {
		toAuth := authMap[to.AuthToUse]
		toRemote, err := createRemote(fromRepo, to.Name, to.URLToRepo)
		if err != nil {
			log.Printf("[WARN] error creating remote in: %v :: %v\n", sc.From.Name, err)
			errRet = err
			continue
		}

		if err := push(ctx, fromRemote, toRemote, toAuth); err != nil {
			log.Printf("[WARN] error syncing to repo: %v :: %v\n", to.Name, err)
			errRet = err
			continue
		}
	}

	return errRet
}

// openRepo opens the git repo and returns an instance to it.
// If the repository doesn't exist, it would clone the repo first.
func openRepo(ctx context.Context, cr conf.Repo, auth *conf.AuthMethod, folder string) (
	*git.Repository, error) {

	if _, errExist := os.Stat(folder); os.IsNotExist(errExist) {
		repo, err := git.PlainCloneContext(ctx, folder, false, &git.CloneOptions{
			URL:        cr.URLToRepo,
			RemoteName: cr.Name,
			Auth:       auth.GetAuth(),
		})
		if err != nil {
			return nil, fmt.Errorf("error cloning the repo [%v] :: %w", cr.Name, err)
		}

		return repo, nil
	}

	repo, err := git.PlainOpen(folder)
	if err != nil {
		return nil, fmt.Errorf("error opening the repo [%v] :: %w", cr.Name, err)
	}

	return repo, nil
}

// createRemote creates a remote in a repository. It would first try to find whether
// the given remote exists. If it does, it will ensure the URL is correctly updated.
// In case it doesn't, it will create a brand new remote in the repository.
func createRemote(repo *git.Repository, name, url string) (*git.Remote, error) {
	remote, err := repo.Remote(name)
	if err == nil {
		remoteURLs := remote.Config().URLs
		if len(remoteURLs) != 1 || remoteURLs[0] != url {
			remote.Config().URLs = []string{url}
		}

		return remote, nil
	} else if err != git.ErrRemoteNotFound {
		return nil, fmt.Errorf("error finding remote [%v] :: %w", name, err)
	}

	c := &config.RemoteConfig{
		Name: name,
		URLs: []string{url},
	}
	if remote, err = repo.CreateRemote(c); err != nil {
		return nil, fmt.Errorf("error creating remote [%v] :: %w", name, err)
	}

	return remote, nil
}

// fetch fetches from provided remote.
func fetch(ctx context.Context, from *git.Remote, am *conf.AuthMethod) error {
	err := from.FetchContext(ctx, &git.FetchOptions{Auth: am.GetAuth()})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error fetching from [%v] :: %w", from.Config().Name, err)
	}

	return nil
}

// push pushes from the `from` remote to the `to` remote.
// `auth` is authentication for `to` auth.
func push(ctx context.Context, from, to *git.Remote, am *conf.AuthMethod) error {
	o := &git.PushOptions{
		RemoteName: to.Config().Name,
		Auth:       am.GetAuth(),
		RefSpecs: []config.RefSpec{
			config.RefSpec(fmt.Sprintf("refs/remotes/%v/master:refs/heads/master", from.Config().Name)),
		},
	}
	if err := to.PushContext(ctx, o); err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error pushing [%v] :: %w", o.RefSpecs, err)
	}

	return nil
}
