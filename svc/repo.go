package svc

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/mangalaman93/giggle/conf"
)

func performSync(ctx context.Context, ch *conf.ConfigHolder) error {
	log.Println("[INFO] periodic sync begin")
	defer log.Println("[INFO] periodic sync end")

	return nil
}

// openRepo opens the git repo and returns an instance to it.
// If the repository doesn't exist, it would clone the repo first.
func openRepo(ctx context.Context, cr conf.Repo, auth *http.BasicAuth, target string) (
	*git.Repository, error) {

	if _, errExist := os.Stat(target); os.IsNotExist(errExist) {
		repo, err := git.PlainCloneContext(ctx, target, false, &git.CloneOptions{
			URL:        cr.URLToRepo,
			RemoteName: cr.Name,
			Auth:       auth,
		})
		if err != nil {
			return nil, fmt.Errorf("error in cloning the repo [%v] :: %w", cr.Name, err)
		}

		return repo, nil
	}

	repo, err := git.PlainOpen(target)
	if err != nil {
		return nil, fmt.Errorf("error in opening the repo [%v] :: %w", cr.Name, err)
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
	} else if err != nil && err != git.ErrRemoteNotFound {
		return nil, fmt.Errorf("error in finding remote [%v] :: %w", name, err)
	}

	c := &config.RemoteConfig{
		Name: name,
		URLs: []string{url},
	}
	if remote, err = repo.CreateRemote(c); err != nil {
		return nil, fmt.Errorf("error in creating remote [%v] :: %w", name, err)
	}

	return remote, nil
}

// fetch fetches from provided remote.
func fetch(ctx context.Context, from *git.Remote) error {
	err := from.FetchContext(ctx, &git.FetchOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error in fetching from [%v] :: %w", from.Config().Name, err)
	}

	return nil
}

// push pushes from the `from` remote to the `to` remote.
// `auth` is authentication for `to` auth.
func push(ctx context.Context, from, to *git.Remote, auth *http.BasicAuth) error {
	refSpec := config.RefSpec(
		fmt.Sprintf("refs/remotes/%s/master:refs/remotes/%s/master",
			from.Config().Name, to.Config().Name),
	)
	err := to.PushContext(ctx, &git.PushOptions{
		RemoteName: to.Config().Name,
		RefSpecs:   []config.RefSpec{refSpec},
		Auth:       auth,
	})
	if err != nil {
		return fmt.Errorf("error in pushing [%v] :: %w", refSpec, err)
	}

	return nil
}
