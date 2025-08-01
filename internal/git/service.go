package git

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type GitService struct {
	author     *object.Signature
	repository *git.Repository
	worktree   *git.Worktree
}

type GitAuthor struct {
	Name, Email string
}

// TODO: Allow to use either repository path or URL
func NewGitService(repoPath string, author *GitAuthor) (*GitService, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			err = fmt.Errorf("repository %q does not exist", repoPath)
		}
		return nil, err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	var authorSig object.Signature
	authorSig.Name = author.Name
	authorSig.Email = author.Email

	return &GitService{
		author:     &authorSig,
		repository: repo,
		worktree:   wt,
	}, nil
}

func (gs *GitService) CommitEmpty(msg string, date time.Time) (string, error) {
	staged, err := gs.StagedFiles()
	if err != nil {
		return "", err
	}
	err = gs.UnstageFiles(staged)
	if err != nil {
		return "", err
	}

	commitHash, err := gs.Commit(msg, date)
	if err != nil {
		return "", err
	}

	err = gs.StageFiles(staged)
	return commitHash, err
}

func (gs *GitService) Commit(msg string, date time.Time) (string, error) {
	gs.author.When = date
	commitHash, err := gs.worktree.Commit(
		msg,
		&git.CommitOptions{
			AllowEmptyCommits: true,
			Author:            gs.author,
		},
	)

	return commitHash.String(), err
}

func (gs *GitService) StagedFiles() ([]string, error) {
	status, err := gs.worktree.Status()
	if err != nil {
		return nil, err
	}

	var staged []string
	for filename, fileStatus := range status {
		switch fileStatus.Staging {
		case git.Unmodified, git.Untracked:
			continue
		}
		staged = append(staged, filename)
	}

	return staged, nil
}

func (gs *GitService) UnstageFiles(files []string) error {
	if len(files) == 0 {
		return nil
	}
	return gs.worktree.Restore(&git.RestoreOptions{
		Staged: true,
		Files:  files,
	})
}

func (gs *GitService) StageFiles(files []string) error {
	for _, file := range files {
		if _, err := gs.worktree.Add(file); err != nil {
			return err
		}
	}

	return nil
}

func GetGlobalUser() (*GitAuthor, error) {
	cfg, err := gitConfig.LoadConfig(gitConfig.GlobalScope)
	if err != nil {
		return nil, err
	}

	return &GitAuthor{
		Name:  cfg.User.Name,
		Email: cfg.User.Email,
	}, nil
}
