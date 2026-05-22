# Simple Git Clone

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-simple-git-clone?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-simple-git-clone/releases)

The step checks out the defined repository state.

<details>
<summary>Description</summary>

The checkout process depends on the checkout properties: the step checks out a repository state defined by a branch, git commit or a git tag.

### Configuring the Step

1. The **Git repository URL** is required field, the address of the repository that needs to be cloned.
2. The **Clone destination** is required field, the local path where we need to clone.
3. In the **Clone Config** section there should be exactly one set of the following:
   - **Branch**
   - **Tag**
   - **Commit**

### Related Steps

- [Activate SSH key (RSA private key)](https://www.bitrise.io/integrations/steps/activate-ssh-key)
- [Bitrise.io Cache:Pull](https://www.bitrise.io/integrations/steps/cache-pull)
- [Bitrise.io Cache:Push](https://www.bitrise.io/integrations/steps/cache-push)
- [Git Clone Repository](https://www.bitrise.io/integrations/steps/git-clone)

</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `repository_url` | The step will try to fetch the repository from the given URL and clone it to the defined clone directory. | required |  |
| `clone_into_dir` | Path of the directory where the cloned repository will be placed. | required |  |
| `commit` | Hash of a commit that needs to be checked out. |  |  |
| `tag` | Git tag that needs to be checked out. |  |  |
| `branch` | Git branch that needs to be checked out. |  |  |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-simple-git-clone/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-simple-git-clone/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Note: this step's end-to-end tests (defined in e2e/bitrise.yml) are working with secrets which are intentionally not stored in this repo. External contributors won't be able to run those tests. Don't worry, if you open a PR with your contribution, we will help with running tests and make sure that they pass.

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
