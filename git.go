package main

import (
	"fmt"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/command/git"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/retry"
	"os"
)

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

func checkout(gitCmd git.Git, arg string, checkoutType checkoutType) error {
	if err := runWithRetry(func() *command.Model {
		var opts []string

		if checkoutType == tag {
			opts = append(opts, "--tags")
		}
		if checkoutType == branch {
			opts = append(opts, "origin", arg)
		}
		return gitCmd.Fetch(opts...)
	}); err != nil {
		return fmt.Errorf("fetch failed: %v", err)
	}

	if err := run(gitCmd.Checkout(arg)); err != nil {
		return fmt.Errorf("checkout failed %s: %v", arg, err)
	}

	return nil
}
