package extractor

import (
	"io"

	parser "github.com/axetroy/whatchanged/1_parser"
	"github.com/axetroy/whatchanged/internal/client"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
)

type ExtractSplice struct {
	Name   string           // The name of splice
	Commit []*object.Commit // The commits including in this splice
	Tag    *client.Tag      // If this splice is a tag splice
}

func getTagOfCommit(tags []*client.Tag, commit *object.Commit) *client.Tag {
	for _, tag := range tags {
		if tag.Commit.Hash.String() == commit.Hash.String() {
			return tag
		}
	}

	return nil
}

func Extract(g *client.GitClient, scope *parser.Scope) ([]*ExtractSplice, error) {
	splices := make([]*ExtractSplice, 0)

	options := git.LogOptions{From: plumbing.NewHash(scope.Start.Commit.Hash.String())}

	cIter, err := g.Repository.Log(&options)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	commits := make([]*object.Commit, 0)

	for {
		if commit, err := cIter.Next(); err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.WithStack(err)
		} else if commit == nil {
			break
		} else if commit.Hash.String() == scope.End.Commit.Hash.String() {
			commits = append(commits, commit)
			break
		} else {
			commits = append(commits, commit)
		}
	}

	if err != nil {
		return nil, errors.WithStack(err)
	}

	// if there is no tags in this scope
	if scope.Tags == nil || len(scope.Tags) == 0 {
		splices = append(splices, &ExtractSplice{
			Name:   "Unreleased",
			Commit: commits,
			Tag:    nil,
		})

		return splices, nil
	}

	index := 0

	for {
		if index == len(commits) {
			break
		}

		commit := commits[index]

		item := &ExtractSplice{
			Name:   "Unreleased",
			Commit: make([]*object.Commit, 0),
		}

		item.Commit = append(item.Commit, commit)

		if tag := getTagOfCommit(scope.Tags, commit); tag != nil {
			item.Tag = tag
			item.Name = tag.Name
		}

		index++

	loop:
		for {
			if index == len(commits) {
				break loop
			}

			nextCommit := commits[index]

			if t := getTagOfCommit(scope.Tags, nextCommit); t != nil {
				break loop
			}

			item.Commit = append(item.Commit, nextCommit)

			index++
		}

		splices = append(splices, item)

	}

	return splices, nil
}
