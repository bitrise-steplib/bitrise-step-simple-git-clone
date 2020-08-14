package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

var configParserMock func(conf *config) error
var storeOpts func(opts []string)
var storeArg func(arg string)

type mockConfigParser struct{}
type mockGitCommandFactory struct{}
type mockGitCommand struct{}

func (m mockGitCommand) init() error {
	return nil
}

func (m mockGitCommand) addRemote(_, _ string) error {
	return nil
}

func (m mockGitCommand) merge(_ string) error {
	return nil
}

func (m mockGitCommand) fetchWithRetry(opts ...string) error {
	storeOpts(opts)
	return nil
}

func (m mockGitCommand) checkout(arg string) error {
	storeArg(arg)
	return nil
}

func (m mockGitCommandFactory) new(_ string) (gitCommand, error) {
	return mockGitCommand{}, nil
}

func (m mockConfigParser) parse(conf *config) error {
	return configParserMock(conf)
}

func TestCloneStep(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Clone Step Suite")
}

var _ = Describe("clone step", func() {
	parser = mockConfigParser{}
	gitFactory = mockGitCommandFactory{}
	Context("parsing inputs", func() {
		When("there is no branch, tag or commit set", func() {
			var e error
			configParserMock = func(conf *config) error {
				return nil
			}
			e = mainE()
			It("should set error", func() {
				Expect(e).To(HaveOccurred())
				Expect(e.Error()).To(Equal("tag, commit or branch input must be set"))
			})
		})
		When("there are more than one branch, tag or commit set", func() {
			var e error
			configParserMock = func(conf *config) error {
				conf.Branch = "test-branch"
				conf.Commit = "test-commit"
				return nil
			}
			e = mainE()
			It("should set error", func() {
				Expect(e).To(HaveOccurred())
				Expect(e.Error()).To(Equal("process inputs: exactly one of [branch, tag, commit] input must be set"))
			})
		})
	})
	Context("checking out", func() {
		When("branch is the input", func() {
			var e error
			var arg string
			var opts []string
			configParserMock = func(conf *config) error {
				conf.Branch = "test-branch"
				return nil
			}
			storeArg = func(a string) {
				arg = a
			}
			storeOpts = func(o []string) {
				opts = o
			}
			e = mainE()
			It("should have branch checkout command", func() {
				Expect(e).NotTo(HaveOccurred())
				Expect(arg).To(Equal("test-branch"))
				Expect(opts).To(HaveLen(2))
				Expect(opts[0]).To(Equal("origin"))
				Expect(opts[1]).To(Equal("test-branch"))
			})
		})
		When("commit is the input", func() {
			var e error
			var arg string
			var opts []string
			configParserMock = func(conf *config) error {
				conf.Commit = "test-commit"
				return nil
			}
			storeArg = func(a string) {
				arg = a
			}
			storeOpts = func(o []string) {
				opts = o
			}
			e = mainE()
			It("should have commit checkout command", func() {
				Expect(e).NotTo(HaveOccurred())
				Expect(arg).To(Equal("test-commit"))
				Expect(opts).To(BeEmpty())
			})
		})
		When("tag is the input", func() {
			var e error
			var arg string
			var opts []string
			configParserMock = func(conf *config) error {
				conf.Tag = "test-tag"
				return nil
			}
			storeArg = func(a string) {
				arg = a
			}
			storeOpts = func(o []string) {
				opts = o
			}
			e = mainE()
			It("should have commit checkout command", func() {
				Expect(e).NotTo(HaveOccurred())
				Expect(arg).To(Equal("test-tag"))
				Expect(opts).To(HaveLen(1))
				Expect(opts[0]).To(Equal("--tags"))
			})
		})
	})
})
