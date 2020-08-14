package main

import (
	"fmt"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/command/git"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/retry"
	"os"
)

type gitCommand interface {
	init() error
	addRemote(name, url string) error
	merge(arg string) error
	fetchWithRetry(opts ...string) error
	checkout(arg string) error
}

type gitCommandFactory interface {
	new(dir string) (gitCommand, error)
}

type realGitCommand struct {
	git git.Git
}

func (r realGitCommand) init() error {
	return run(r.git.Init())
}

func (r realGitCommand) addRemote(name, url string) error {
	return run(r.git.RemoteAdd(name, url))
}

func (r realGitCommand) merge(arg string) error {
	return run(r.git.Merge(arg))
}

func (r realGitCommand) fetchWithRetry(opts ...string) error {
	return runWithRetry(func() *command.Model {
		return r.git.Fetch(opts...)
	})
}

func (r realGitCommand) checkout(arg string) error {
	return run(r.git.Checkout(arg))
}

type realGitCommandFactory struct{}

func (r realGitCommandFactory) new(dir string) (gitCommand, error) {
	g, err := git.New(dir)
	if err != nil {
		return nil, err
	}
	return realGitCommand{
		git: g,
	}, nil
}

func run(c *command.Model) error {
	log.Infof(c.PrintableCommandArgs())
	return c.SetStdout(os.Stdout).SetStderr(os.Stderr).Run()
}

func runWithRetry(f func() *command.Model) error {
	return retry.Times(2).Wait(5).Try(func(attempt uint) error {
		if attempt > 0 {
			log.Warnf("Retrying...")
		}

		err := run(f())
		if err != nil {
			log.Warnf("Attempt %d failed:", attempt+1)
			fmt.Println(err.Error())
		}

		return err
	})
}

func checkout(gitCmd gitCommand, arg string, checkoutType checkoutType) error {
	opts := buildFetchOpts(checkoutType, arg)

	if err := gitCmd.fetchWithRetry(opts...); err != nil {
		return fmt.Errorf("fetch failed: %v", err)
	}

	if err := gitCmd.checkout(arg); err != nil {
		return fmt.Errorf("checkout failed %s: %v", arg, err)
	}

	return nil
}

func buildFetchOpts(checkoutType checkoutType, arg string) []string {
	var opts []string

	if checkoutType == tag {
		opts = append(opts, "--tags")
	}
	if checkoutType == branch {
		opts = append(opts, "origin", arg)
	}
	return opts
}
