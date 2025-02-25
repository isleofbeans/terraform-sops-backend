# HOW TO: Contribute

[![readme](./docs/assets/breadcrum-readme.drawio.svg)](./README.md)[![how-to-guides](./docs/assets/breadcrum-how-to-guides.drawio.svg)](./docs/how-to/index.md)

## About

Thank you for investing your time in contributing to our project!

Any contribution you make will be reflected on [terraform-sops-backend](https://github.com/isleofbeans/terraform-sops-backend).

In this guide you will get an overview of the contribution workflow from creating a PR, reviewing, and merging the PR.

Use the table of contents icon on the top left corner of this document to get to a specific section of this guide quickly.

## Getting started

This repository is verified on commit time by using [pre-commit](https://pre-commit.com/) hooks. To have them working you need to initialize them. Please follow the proper [installation instructions](https://pre-commit.com/#install).

### Required development binaries

For the list of required binaries see the [ASDF configuration](./.tool-versions)  
The simplest way to install them is to use [ASDF](https://asdf-vm.com/guide/getting-started.html)

### Make Changes

#### Make changes locally

1. Clone the repository.
    - Using GitHub Desktop:
        - [Getting started with GitHub Desktop](https://docs.github.com/en/desktop/installing-and-configuring-github-desktop/getting-started-with-github-desktop) will guide you through setting up Desktop.
        - Once Desktop is set up, you can use it to [clone the repo](https://docs.github.com/en/desktop/contributing-and-collaborating-using-github-desktop/adding-and-cloning-repositories/cloning-and-forking-repositories-from-github-desktop#cloning-a-repository)!

    - Using the command line:
        - [Clone the repo](https://docs.github.com/en/get-started/quickstart/fork-a-repo#cloning-your-forked-repository) so that you can make your changes without affecting the original project until you're ready to merge them.
            ```sh
            git clone git@github.com:isleofbeans/terraform-sops-backend.git
            ```

2. Enable pre-commit hooks

    Enable the pre-commit hooks for this repository

    ```sh
    pre-commit install
    ```

5. Setup a development environment

    Since the unit tests require some external functionality please setup your environment following the [HOW TO: Setup development environment](./docs/how-to/setup-development-environment.md)

6. Create a working branch and start with your changes!
    - [Using GitHub Desktop](https://docs.github.com/en/desktop/contributing-and-collaborating-using-github-desktop/making-changes-in-a-branch/managing-branches#creating-a-branch)
    - Using the command line:
        ```sh
        git checkout -b <useful-branch-name>
        ```

### Commit your update

Commit the changes once you are happy with them. See [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) to know how to craft good commit messages.

### Pull Request

When you're finished with the changes, create a pull request, also known as a PR.

- Give a describing name to the PR
- Fill the description body with
    - A description of the change done by the PR
    - (Optional) Required migration paths
    - The issue number related to this PR
- Once you submit your PR, a repository manager will review your proposal. We may ask questions or request for additional information.
- We may ask for changes to be made before a PR can be merged, either using [suggested changes](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/incorporating-feedback-in-your-pull-request) or pull request comments. You can apply suggested changes directly through the UI. You can make any other changes in your fork, then commit them to your branch.
- As you update your PR and apply changes, mark each conversation as [resolved](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/commenting-on-a-pull-request#resolving-conversations).
- If you run into any merge issues, checkout this [git tutorial](https://lab.github.com/githubtraining/managing-merge-conflicts) to help you resolve merge conflicts and other issues.

### Your PR is approved and merged

Congratulations :tada::tada: We thanks you :sparkles:.
