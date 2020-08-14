package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSimpleCloneStepNoInputFails(t *testing.T) {
	// Arrange
	var err error
	parser := new(mockConfigParser)
	gitFactory := new(mockGitCommandFactory)
	cloner := simpleGitCloner{
		parser:     parser,
		gitFactory: gitFactory,
	}

	parser.On("parse", mock.Anything).Return(nil)

	// Act
	err = cloner.clone()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "tag, commit or branch input must be set")
}

func TestSimpleCloneStepMultipleInputsFails(t *testing.T) {
	// Arrange
	var err error
	parser := new(mockConfigParser)
	gitFactory := new(mockGitCommandFactory)
	cloner := simpleGitCloner{
		parser:     parser,
		gitFactory: gitFactory,
	}

	parser.On("parse", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		config, ok := args.Get(0).(*config)
		if ok {
		}
		config.Branch = "test-branch"
		config.Commit = "test-commit"
	})

	// Act
	err = cloner.clone()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "process inputs: exactly one of [branch, tag, commit] input must be set")
}

func TestSimpleCloneStepBranchInputCreatesBranchCommand(t *testing.T) {
	// Arrange
	var err error
	parser := new(mockConfigParser)
	gitFactory := new(mockGitCommandFactory)
	gitCommand := new(mockGitCommand)
	cloner := simpleGitCloner{
		parser:     parser,
		gitFactory: gitFactory,
	}

	parser.On("parse", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		config, ok := args.Get(0).(*config)
		if ok {
		}
		config.Branch = "test-branch"
	})
	gitFactory.On("new", mock.Anything).Return(gitCommand, nil)
	gitCommand.On("init").Return(nil)
	gitCommand.On("addRemote", mock.Anything, mock.Anything).Return(nil)
	gitCommand.On("fetchWithRetry", mock.Anything, mock.Anything).Return(nil)
	gitCommand.On("checkout", mock.Anything).Return(nil)
	gitCommand.On("merge", mock.Anything).Return(nil)

	// Act
	err = cloner.clone()

	// Assert
	assert.NoError(t, err)
	gitCommand.AssertNumberOfCalls(t, "fetchWithRetry", 1)
	gitCommand.AssertCalled(t, "fetchWithRetry", "origin", "test-branch")
	gitCommand.AssertNumberOfCalls(t, "checkout", 1)
	gitCommand.AssertCalled(t, "checkout", "test-branch")
}

func TestSimpleCloneStepTagInputCreatesTagCommand(t *testing.T) {
	// Arrange
	var err error
	parser := new(mockConfigParser)
	gitFactory := new(mockGitCommandFactory)
	gitCommand := new(mockGitCommand)
	cloner := simpleGitCloner{
		parser:     parser,
		gitFactory: gitFactory,
	}

	parser.On("parse", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		config, ok := args.Get(0).(*config)
		if ok {
		}
		config.Tag = "test-tag"
	})
	gitFactory.On("new", mock.Anything).Return(gitCommand, nil)
	gitCommand.On("init").Return(nil)
	gitCommand.On("addRemote", mock.Anything, mock.Anything).Return(nil)
	gitCommand.On("fetchWithRetry", mock.Anything).Return(nil)
	gitCommand.On("checkout", mock.Anything).Return(nil)
	gitCommand.On("merge", mock.Anything).Return(nil)

	// Act
	err = cloner.clone()

	// Assert
	assert.NoError(t, err)
	gitCommand.AssertNumberOfCalls(t, "fetchWithRetry", 1)
	gitCommand.AssertCalled(t, "fetchWithRetry", "--tags")
	gitCommand.AssertNumberOfCalls(t, "checkout", 1)
	gitCommand.AssertCalled(t, "checkout", "test-tag")
}

func TestSimpleCloneStepCommitInputCreatesCommitCommand(t *testing.T) {
	// Arrange
	var err error
	parser := new(mockConfigParser)
	gitFactory := new(mockGitCommandFactory)
	gitCommand := new(mockGitCommand)
	cloner := simpleGitCloner{
		parser:     parser,
		gitFactory: gitFactory,
	}

	parser.On("parse", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		config, ok := args.Get(0).(*config)
		if ok {
		}
		config.Commit = "test-commit"
	})
	gitFactory.On("new", mock.Anything).Return(gitCommand, nil)
	gitCommand.On("init").Return(nil)
	gitCommand.On("addRemote", mock.Anything, mock.Anything).Return(nil)
	gitCommand.On("fetchWithRetry").Return(nil)
	gitCommand.On("checkout", mock.Anything).Return(nil)
	gitCommand.On("merge", mock.Anything).Return(nil)

	// Act
	err = cloner.clone()

	// Assert
	assert.NoError(t, err)
	gitCommand.AssertNumberOfCalls(t, "fetchWithRetry", 1)
	gitCommand.AssertCalled(t, "fetchWithRetry")
	gitCommand.AssertNumberOfCalls(t, "checkout", 1)
	gitCommand.AssertCalled(t, "checkout", "test-commit")
}
