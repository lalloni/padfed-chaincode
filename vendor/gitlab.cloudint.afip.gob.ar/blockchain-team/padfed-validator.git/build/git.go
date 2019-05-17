package build

import (
	"errors"
	"io"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func GitWorktreeNotDirty(path string) error {
	dirty, err := GitWorktreeDirty(path)
	if err != nil {
		return err
	}
	if dirty {
		return errors.New("working directory is dirty")
	}
	return nil
}

func GitWorktreeDirty(path string) (bool, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return false, err

	}
	wt, err := repo.Worktree()
	if err != nil {
		return false, err
	}
	st, err := wt.Status()
	if err != nil {
		return false, err
	}
	return !st.IsClean(), nil
}

func GitTag(path, tag string) (*plumbing.Reference, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	head, err := repo.Head()
	if err != nil {
		return nil, err
	}
	return repo.CreateTag(tag, head.Hash(), nil)
}

func GitTags(path string) ([]string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	tags, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	defer tags.Close()
	res := []string{}
	for {
		tag, err := tags.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		res = append(res, tag.Name().Short())
	}
	return res, nil
}

func GitTagExist(path, tag string) (bool, error) {
	tags, err := GitTags(path)
	if err != nil {
		return false, err
	}
	for _, t := range tags {
		if t == tag {
			return true, nil
		}
	}
	return false, nil
}

func GitTagNotExist(path, tag string) error {
	exist, err := GitTagExist(path, tag)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("tag exist")
	}
	return nil
}
