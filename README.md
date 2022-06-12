# helloGo

`helloGo` is a repository with quick samples on how to do some common operations and features from `goLang`.

## Pipelines

We use GitHub Actions to set up and execute our pipelines.

The files controlling each pipeline can be found within the [.gitHub](./.github) repository.

### Commits and Pull Requests

For every commit and pull request the [test.yml](./.github/workflows/test.yml) pipeline is triggered to execute all the
tests within the project.

### Dependabot

For semi-automatic updates on our dependencies we use GitHub's Dependabot. The settings can be found on
the [dependabot.yml](./.github/dependabot.yml) file.
