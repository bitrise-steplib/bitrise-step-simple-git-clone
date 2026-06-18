package main

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/retry"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/git"
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
	gitFactory git.Factory
}

func (r realGitCommand) init() error {
	return run(r.gitFactory.Init().Create(os.Stdout, os.Stderr, nil))
}

func (r realGitCommand) addRemote(name, url string) error {
	return run(r.gitFactory.RemoteAdd(name, url).Create(os.Stdout, os.Stderr, nil))
}

func (r realGitCommand) merge(arg string) error {
	return run(r.gitFactory.Merge(arg).Create(os.Stdout, os.Stderr, nil))
}

func (r realGitCommand) fetchWithRetry(opts ...string) error {
	return runWithRetry(func() command.Command {
		return r.gitFactory.Fetch(opts...).Create(os.Stdout, os.Stderr, nil)
	})
}

func (r realGitCommand) checkout(arg string) error {
	return run(r.gitFactory.Checkout(arg).Create(os.Stdout, os.Stderr, nil))
}

type realGitCommandFactory struct{}

func (r realGitCommandFactory) new(dir string) (gitCommand, error) {
	envRepo := env.NewRepository()
	cmdFactory := command.NewFactory(envRepo)
	g, err := git.NewFactory(dir, cmdFactory, nil)
	if err != nil {
		return nil, err
	}
	return realGitCommand{
		gitFactory: g,
	}, nil
}

func run(c command.Command) error {
	log.Infof(c.PrintableCommandArgs())
	return c.Run()
}

func runWithRetry(f func() command.Command) error {
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
		opts = append(opts, "origin", "refs/heads/"+arg)
	}
	return opts
}
