# Simple Git Clone (internal step)

As step developers, we often have to clone sample repos in E2E testing workflows. For this the standard git clone step is not usable since it uses hard coded bitrise input variables.

The usage is pretty similar to the standard clone step but is a bit more
relaxed and has less options. See `step.yml` for details.

The checkout process depends on the checkout properties: the step checks out a repository state defined by a branch, git commit, or a git tag.
