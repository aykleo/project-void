package git

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Commit struct {
	Hash      string
	Author    string
	Message   string
	Timestamp time.Time
}

func GetCommitsSince(repoPath string, since time.Time) ([]Commit, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	refs, err := repo.References()
	if err != nil {
		return nil, fmt.Errorf("failed to get references: %w", err)
	}

	uniqueCommits := make(map[string]Commit)

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name() == plumbing.HEAD {
			return nil
		}

		commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
		if err != nil {
			return nil
		}
		defer commitIter.Close()

		err = commitIter.ForEach(func(c *object.Commit) error {
			if c.Author.When.After(since) {
				hash := c.Hash.String()
				if _, exists := uniqueCommits[hash]; !exists {
					uniqueCommits[hash] = Commit{
						Hash:      hash,
						Author:    c.Author.Name,
						Message:   c.Message,
						Timestamp: c.Author.When,
					}
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error iterating references: %w", err)
	}

	commits := make([]Commit, 0, len(uniqueCommits))
	for _, commit := range uniqueCommits {
		commits = append(commits, commit)
	}

	return commits, nil
}
