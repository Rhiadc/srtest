package git

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	. "github.com/go-git/go-git/v5/_examples"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type GitConfig struct {
	RepoURL           string
	RepoInternalPath  string
	RepoReferenceName string
	Username          string
	Password          string
	OwnerName         string
	OnwerEmail        string
}

type GitClient struct {
	config GitConfig
}

func NewGitClient(config GitConfig) *GitClient {
	return &GitClient{config: config}
}

func (gc GitClient) CloneRepo() (*git.Repository, error) {
	//@todo pass logs
	// Clone the repository if it doesn't exist
	repo, err := git.PlainOpen(gc.config.RepoInternalPath)
	basicAuth := gc.getBasicAuth()
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			repo, err = git.PlainClone(gc.config.RepoInternalPath, false, &git.CloneOptions{
				URL:           gc.config.RepoURL,
				Auth:          basicAuth,
				ReferenceName: plumbing.ReferenceName(gc.config.RepoReferenceName),
				SingleBranch:  true,
				Progress:      os.Stdout,
			})
			if err != nil && !errors.Is(err, git.ErrRepositoryAlreadyExists) {
				fmt.Println("Failed to clone repository:", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Failed to open repository:")
		}
	}
	return repo, err
}

func (GitClient) PullFromMain(repo *git.Repository) (*git.Worktree, error) {
	Info("git pull origin main")
	w, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil, err
	}
	return w, nil
}

func (gc GitClient) AddCommitAndPush(repo *git.Repository, w *git.Worktree, path string) error {
	if err := gc.add(w, path); err != nil {
		return err
	}

	if err := gc.commit(w); err != nil {
		return err
	}

	if err := gc.push(repo); err != nil {
		return err
	}

	return nil
}

func (GitClient) add(w *git.Worktree, path string) error {
	Info("git add example-git-file")
	_, err := w.Add(path)
	return err
}

func (gc GitClient) commit(w *git.Worktree) error {
	Info("git commit")
	_, err := w.Commit("Adding a new file", &git.CommitOptions{
		Author: &object.Signature{
			Name:  gc.config.OwnerName,
			Email: gc.config.OnwerEmail,
			When:  time.Now(),
		},
	})
	return err
}

func (gc GitClient) push(repo *git.Repository) error {
	Info("git push origin main")
	basicAuth := gc.getBasicAuth()
	err := repo.Push(&git.PushOptions{
		Auth: basicAuth,
	})
	return err
}

func (gc GitClient) getBasicAuth() *http.BasicAuth {
	return &http.BasicAuth{Username: gc.config.Username, Password: gc.config.Password}
}
