package repository

import (
	"errors"
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type ReferenceType string

var (
	BranchReferenceType ReferenceType = "BranchReferenceType"
	TagReferenceType    ReferenceType = "TagReferenceType"
)

type repository struct {
	url, remote, branch, key string
	repo                     *git.Repository
	logger                   io.Writer
}

func (r *repository) Open(directory string) error {
	repo, err := git.PlainOpenWithOptions(directory, &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: false,
	})

	if err != nil {
		return err
	}

	if r.repo != nil {
		return errors.New("Another repository is already loaded")
	}

	r.repo = repo

	return nil
}

func (r *repository) Clone(directory string) error {
	if r.remote == "" {
		r.remote = "origin"
	}

	referenceName := plumbing.NewBranchReferenceName(r.branch)

	repo, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL: r.url,
		Auth: &http.BasicAuth{
			Username: "place-holder",
			Password: r.key,
		},
		RemoteName:    r.remote,
		ReferenceName: referenceName,
		SingleBranch:  true,
		NoCheckout:    false,
		Progress:      r.logger,
	})

	if err != nil {
		return err
	}

	if r.repo != nil {
		return errors.New("another repository is already loaded")
	}

	r.repo = repo

	return nil
}

func (r *repository) Checkout(commit string) error {
	w, err := r.repo.Worktree()
	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commit),
	})

	return err
}

func (r *repository) Push() error {
	err := r.repo.Push(&git.PushOptions{
		RemoteName: r.remote,
	})

	return err
}

func (r *repository) Pull(branch string) error {
	var referenceName plumbing.ReferenceName
	if r.branch == "" {
		referenceName = plumbing.NewRemoteHEADReferenceName(r.remote)
	} else {
		referenceName = plumbing.NewRemoteReferenceName(r.remote, r.branch)
	}

	w, err := r.repo.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{
		RemoteName:    r.remote,
		ReferenceName: referenceName,
	})

	return err
}

func New(url, key, remote, branch string, logger io.Writer) *repository {
	return &repository{url: url, key: key, remote: remote, branch: branch, logger: logger}
}
