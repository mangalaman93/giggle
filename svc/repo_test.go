package svc

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/mangalaman93/giggle/conf"
)

func createTestDir(t *testing.T, testDir string) {
	if err := os.MkdirAll(testDir, 0700); err != nil {
		t.Fatalf("unable to create test dir: %v", err)
	}
}

func deleteTestDir(t *testing.T, testDir string) {
	if err := os.RemoveAll(testDir); err != nil {
		t.Fatalf("unable to delete test dir: %v", err)
	}
}

func createFile(filePath, content string) error {
	fileFD, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v :: %w", filePath, err)
	}
	if _, err := fileFD.WriteString(content); err != nil {
		return fmt.Errorf("error writing to file: %v :: %w", filePath, err)
	}
	if err := fileFD.Close(); err != nil {
		return fmt.Errorf("error closing file: %v :: %w", filePath, err)
	}
	return nil
}

func commit(repo *git.Repository, message string) error {
	workTree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("error getting worktree :: %w", err)
	}
	if err := workTree.AddGlob("."); err != nil {
		return fmt.Errorf("error adding file to worktree :: %w", err)
	}

	sig := &object.Signature{
		Name:  "user",
		Email: "user@giggle",
		When:  time.Now(),
	}
	co := &git.CommitOptions{Author: sig}
	if _, err := workTree.Commit(message, co); err != nil {
		return fmt.Errorf("error committing :: %w", err)
	}

	return nil
}

func setupGitRepo(repoPath string) (*git.Repository, error) {
	repo, err := git.PlainInit(repoPath, false)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize a git repo :: %w", err)
	}

	if err := repo.CreateBranch(&config.Branch{Name: "master"}); err != nil {
		return nil, fmt.Errorf("unable to create master branch :: %w", err)
	}

	readmeFile := filepath.Join(repoPath, "README.md")
	if err := createFile(readmeFile, "This is an example project."); err != nil {
		return nil, err
	}

	if err := commit(repo, "Add README"); err != nil {
		return nil, err
	}

	return repo, nil
}

func TestOpenRepo(t *testing.T) {
	testDir := path.Join(os.TempDir(), fmt.Sprintf("testdir_%d", rand.Intn(1000)))
	createTestDir(t, testDir)
	defer deleteTestDir(t, testDir)
	if _, err := setupGitRepo(testDir); err != nil {
		t.Fatalf("error setting up git repo :: %v", err)
	}
	repoDir := path.Join(os.TempDir(), fmt.Sprintf("testdir_%d", rand.Intn(1000)))
	defer deleteTestDir(t, repoDir)

	// open repo
	cr := conf.Repo{
		Name:      "local",
		Kind:      "local",
		URLToRepo: fmt.Sprintf("file://%v", testDir),
	}
	clonedRepo, err := openRepo(context.Background(), cr, nil, repoDir)
	if err != nil {
		t.Fatalf("unable to open a git repo :: %v", err)
	}
	if _, err := clonedRepo.Branch("master"); err != nil {
		t.Fatalf("master branch not found :: %v", err)
	}

	// we should be able to open the repo again.
	clonedRepo2, err := openRepo(context.Background(), cr, nil, repoDir)
	if err != nil {
		t.Fatalf("unable to open a git repo again :: %v", err)
	}
	if _, err := clonedRepo2.Branch("master"); err != nil {
		t.Fatalf("master branch not found :: %v", err)
	}

	// ensure all the files are in place.
	files, err := ioutil.ReadDir(repoDir)
	if err != nil {
		t.Fatalf("unable to read cloned repo dir :: %v", err)
	}
	found := false
	for _, file := range files {
		if file.Name() == "README.md" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("README.md not found in repo")
	}
}

func TestCreateRemote(t *testing.T) {
	testDir := path.Join(os.TempDir(), fmt.Sprintf("testdir_%d", rand.Intn(1000)))
	createTestDir(t, testDir)
	defer deleteTestDir(t, testDir)
	if _, err := setupGitRepo(testDir); err != nil {
		t.Fatalf("error setting up git repo :: %v", err)
	}
	repoDir := path.Join(os.TempDir(), fmt.Sprintf("testdir_%d", rand.Intn(1000)))
	defer deleteTestDir(t, repoDir)

	cr := conf.Repo{
		Name:      "local",
		Kind:      "local",
		URLToRepo: fmt.Sprintf("file://%v", testDir),
	}
	repo, err := openRepo(context.Background(), cr, nil, repoDir)
	if err != nil {
		t.Fatalf("error opening repo :: %v", err)
	}

	// create a remote
	remoteName := "overleaf"
	remoteURL := "https://git.overleaf.com/temp"
	remote, err := createRemote(repo, remoteName, remoteURL)
	if err != nil {
		t.Fatalf("error creating remote :: %v", err)
	}
	if remote.Config().Name != remoteName {
		t.Fatalf("unexpected name for remote: %v", remote.Config().Name)
	}
	if len(remote.Config().URLs) != 1 || remote.Config().URLs[0] != remoteURL {
		t.Fatalf("unexpected URLs for remote: %v", remote.Config().URLs)
	}

	// create the remote again with updated URL
	remoteNewURL := "https://git.overleaf.com/newtemp"
	remote, err = createRemote(repo, remoteName, remoteNewURL)
	if err != nil {
		t.Fatalf("error creating remote :: %v", err)
	}
	if len(remote.Config().URLs) != 1 || remote.Config().URLs[0] != remoteNewURL {
		t.Fatalf("unexpected URLs for remote: %v", remote.Config().URLs)
	}

	// create one more remote
	remote2Name := "github"
	remote2URL := "https://github.com/temp/temp"
	remote2, err := createRemote(repo, remote2Name, remote2URL)
	if err != nil {
		t.Fatalf("error creating remote :: %v", err)
	}
	if remote2.Config().Name != remote2Name {
		t.Fatalf("unexpected name for remote: %v", remote2.Config().Name)
	}
}

func TestFetch(t *testing.T) {
	testDir := path.Join(os.TempDir(), fmt.Sprintf("testdir_%d", rand.Intn(1000)))
	createTestDir(t, testDir)
	defer deleteTestDir(t, testDir)

	testRepo, err := setupGitRepo(testDir)
	if err != nil {
		t.Fatalf("error setting up git repo :: %v", err)
	}
	repoDir := path.Join(os.TempDir(), fmt.Sprintf("testdir_%d", rand.Intn(1000)))
	defer deleteTestDir(t, repoDir)

	cr := conf.Repo{
		Name:      "local",
		Kind:      "local",
		URLToRepo: fmt.Sprintf("file://%v", testDir),
	}
	repo, err := openRepo(context.Background(), cr, nil, repoDir)
	if err != nil {
		t.Fatalf("error opening repo :: %v", err)
	}

	remote, err := repo.Remote("local")
	if err != nil {
		t.Fatalf("error finding local remote :: %v", err)
	}
	if err := fetch(context.Background(), remote, nil); err != nil {
		t.Fatalf("error fetching from remote :: %v", err)
	}

	// add new file main.go to test repo.
	codeFile := filepath.Join(testDir, "main.go")
	if err := createFile(codeFile, "package main\n"); err != nil {
		t.Fatalf("error creating main.go :: %v", err)
	}
	if err := commit(testRepo, "Add main.go"); err != nil {
		t.Fatalf("error committing main.go :: %v", err)
	}

	// fetch new file
	if err := fetch(context.Background(), remote, nil); err != nil {
		t.Fatalf("error fetching from remote :: %v", err)
	}
	files, err := ioutil.ReadDir(repoDir)
	if err != nil {
		t.Fatalf("unable to read cloned repo dir :: %v", err)
	}
	fileNames := make([]string, len(files))
	for i, file := range files {
		fileNames[i] = file.Name()
	}

	// count the commits
	iter, err := repo.Log(&git.LogOptions{All: true})
	if err != nil {
		t.Fatalf("error getting all the commits :: %v", err)
	}
	counter := 0
	if err := iter.ForEach(func(*object.Commit) error {
		counter++
		return nil
	}); err != nil {
		t.Fatalf("error iterating commits :: %v", err)
	}
	iter.Close()
	if counter != 2 {
		t.Fatalf("fetch didn't work as expected")
	}
}

func TestPush(t *testing.T) {
	repo1Dir := path.Join(os.TempDir(), fmt.Sprintf("testdir_%d", rand.Intn(1000)))
	createTestDir(t, repo1Dir)
	defer deleteTestDir(t, repo1Dir)
	repo1, err := setupGitRepo(repo1Dir)
	if err != nil {
		t.Fatalf("error setting up git repo :: %v", err)
	}

	repo2Dir := path.Join(os.TempDir(), fmt.Sprintf("testdir_%d", rand.Intn(1000)))
	defer deleteTestDir(t, repo2Dir)
	cr := conf.Repo{
		Name:      "repo1",
		Kind:      "repo1",
		URLToRepo: fmt.Sprintf("file://%v", repo1Dir),
	}
	repo2, err := openRepo(context.Background(), cr, nil, repo2Dir)
	if err != nil {
		t.Fatalf("error opening repo :: %v", err)
	}

	repo3Dir := path.Join(os.TempDir(), fmt.Sprintf("testdir_%d", rand.Intn(1000)))
	defer deleteTestDir(t, repo3Dir)
	repo3, err := openRepo(context.Background(), cr, nil, repo3Dir)
	if err != nil {
		t.Fatalf("error opening repo :: %v", err)
	}
	wt3, err := repo3.Worktree()
	if err != nil {
		t.Fatalf("error in getting the work tree :: %v", err)
	}
	// checkout another branch, so that push doesn't fail.
	if err := wt3.Checkout(&git.CheckoutOptions{
		Branch: "dev",
		Create: true,
	}); err != nil {
		t.Fatalf("error in creating a branch :: %v", err)
	}

	remote1, err := repo2.Remote("repo1")
	if err != nil {
		t.Fatalf("error finding repo1 remote :: %v", err)
	}
	repo3URL := fmt.Sprintf("file://%v", repo3Dir)
	remote3, err := createRemote(repo2, "repo3", repo3URL)
	if err != nil {
		t.Fatalf("error creating repo3 remote :: %v", err)
	}

	// add new file main.go to repo1.
	codeFile := filepath.Join(repo1Dir, "main.go")
	if err := createFile(codeFile, "package main\n"); err != nil {
		t.Fatalf("error creating main.go :: %v", err)
	}
	if err := commit(repo1, "Add main.go"); err != nil {
		t.Fatalf("error committing main.go :: %v", err)
	}

	// pull changes from repo1 to repo2
	if err := fetch(context.Background(), remote1, nil); err != nil {
		t.Fatalf("error fetching from repo1 :: %v", err)
	}

	// push changes from repo2 to repo3 (that are on repo1 remote)
	if err := push(context.Background(), remote1, remote3, nil); err != nil {
		t.Fatalf("error pushing from repo1 to repo3 :: %v", err)
	}

	// count the commits
	iter, err := repo3.Log(&git.LogOptions{All: true})
	if err != nil {
		t.Fatalf("error getting all the commits :: %v", err)
	}
	counter := 0
	if err := iter.ForEach(func(*object.Commit) error {
		counter++
		return nil
	}); err != nil {
		t.Fatalf("error iterating commits :: %v", err)
	}
	iter.Close()
	if counter != 2 {
		t.Fatalf("fetch didn't work as expected")
	}
}
