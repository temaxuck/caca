package git

import (
	"time"

	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type GitService struct {
	Repository *git.Repository
	Author     *object.Signature
	Worktree   *git.Worktree
}

type GitAuthor struct {
	Name, Email string
}

// TODO: Allow to use either repository path or URL
func NewGitService(repoPath string, author *GitAuthor) (*GitService, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	cfg, err := repo.ConfigScoped(gitConfig.GlobalScope)
	if err != nil {
		return nil, err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	var _author object.Signature
	if author == nil {
		_author.Name = cfg.User.Name
		_author.Email = cfg.User.Email
	} else {
		_author.Name = author.Name
		_author.Email = author.Email
	}

	return &GitService{
		Author:     &_author,
		Repository: repo,
		Worktree:   wt,
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
	gs.Author.When = date
	commitHash, err := gs.Worktree.Commit(
		msg,
		&git.CommitOptions{
			AllowEmptyCommits: true,
			Author:            gs.Author,
		},
	)

	return commitHash.String(), err
}

func (gs *GitService) StagedFiles() ([]string, error) {
	status, err := gs.Worktree.Status()
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
	return gs.Worktree.Restore(&git.RestoreOptions{
		Staged: true,
		Files:  files,
	})
}

func (gs *GitService) StageFiles(files []string) error {
	for _, file := range files {
		if _, err := gs.Worktree.Add(file); err != nil {
			return err
		}
	}

	return nil
}
