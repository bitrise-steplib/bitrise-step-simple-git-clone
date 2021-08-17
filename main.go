package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/log"
)

type configParser interface {
	parse(conf *config) error
}

type bitriseConfigParser struct{}

func (b bitriseConfigParser) parse(conf *config) error {
	return stepconf.Parse(conf)
}

type config struct {
	RepositoryURL string `env:"repository_url,required"`
	CloneIntoDir  string `env:"clone_into_dir,required"`
	Commit        string `env:"commit"`
	Tag           string `env:"tag"`
	Branch        string `env:"branch"`
}

type checkoutType string

const (
	commit checkoutType = "commit"
	tag                 = "tag"
	branch              = "branch"
)

type simpleGitCloner struct {
	parser     configParser
	gitFactory gitCommandFactory
}

func (c simpleGitCloner) clone() error {
	var cfg config
	if err := c.parser.parse(&cfg); err != nil {
		return fmt.Errorf("parse step configuration: %v", err)
	}
	stepconf.Print(cfg)

	var checkoutType checkoutType
	var checkoutArg string

	var setCheckoutArg = func(arg string) error {
		if checkoutArg != "" {
			return errors.New("exactly one of [branch, tag, commit] input must be set")
		}
		checkoutArg = arg
		return nil
	}

	if cfg.Commit != "" {
		checkoutType = commit
		checkoutArg = cfg.Commit
	}
	if cfg.Tag != "" {
		checkoutType = tag
		if err := setCheckoutArg(cfg.Tag); err != nil {
			return fmt.Errorf("process inputs: %v", err)
		}
	}
	if cfg.Branch != "" {
		checkoutType = branch
		if err := setCheckoutArg(cfg.Branch); err != nil {
			return fmt.Errorf("process inputs: %v", err)
		}
	}
	if checkoutArg == "" {
		return errors.New("tag, commit or branch input must be set")
	}

	gitCmd, err := c.gitFactory.new(cfg.CloneIntoDir)
	if err != nil {
		return fmt.Errorf("create gitCmd project: %v", err)
	}

	if err := gitCmd.init(); err != nil {
		return fmt.Errorf("init repository: %v", err)
	}
	if err := gitCmd.addRemote("origin", cfg.RepositoryURL); err != nil {
		return fmt.Errorf("add remote repository %s: %v", cfg.RepositoryURL, err)
	}

	if err := checkout(gitCmd, checkoutArg, checkoutType); err != nil {
		return fmt.Errorf("checkout %s: %v", checkoutArg, err)
	}
	// Update branch: 'git fetch' followed by a 'git merge' is the same as 'git pull'.
	if checkoutType == branch {
		if err := gitCmd.merge("origin/" + checkoutArg); err != nil {
			return fmt.Errorf("merge %s: %v", checkoutArg, err)
		}
	}

	return nil
}

var simpleCloner simpleGitCloner

func init() {
	simpleCloner = simpleGitCloner{
		parser:     bitriseConfigParser{},
		gitFactory: realGitCommandFactory{},
	}
}

func main() {
	if err := simpleCloner.clone(); err != nil {
		log.Errorf("ERROR: %v", err)
		os.Exit(1)
	}
	log.Donef("\nSuccess")
}
