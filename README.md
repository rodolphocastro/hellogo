# helloGo

`helloGo` is a repository with quick samples on how to do some common operations and features from `goLang`.

## Status

[![SonarCloud](https://sonarcloud.io/images/project_badges/sonarcloud-orange.svg)](https://sonarcloud.io/summary/new_code?id=rodolphocastro_hellogo)

---

[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=rodolphocastro_hellogo&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=rodolphocastro_hellogo)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=rodolphocastro_hellogo&metric=coverage)](https://sonarcloud.io/summary/new_code?id=rodolphocastro_hellogo)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=rodolphocastro_hellogo&metric=bugs)](https://sonarcloud.io/summary/new_code?id=rodolphocastro_hellogo)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=rodolphocastro_hellogo&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=rodolphocastro_hellogo)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=rodolphocastro_hellogo&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=rodolphocastro_hellogo)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=rodolphocastro_hellogo&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=rodolphocastro_hellogo)

## Pipelines

We use GitHub Actions to set up and execute our pipelines.

The files controlling each pipeline can be found within the [.gitHub](./.github) repository.

### Commits and Pull Requests

For every pull request the [pull-request.yml](./.github/workflows/pull-request.yml) pipeline is triggered to execute all the
tests within the project, integration and unit both.

Every commit triggers the [vet.yml](./.github/workflows/vet.yml) pipeline to run unit test and `go vet` upon the files.

Finally, commits to `master` trigger the [deploy.yml](./.github/workflows/deploy.yml) pipeline that may deploy stuff once tests are executed.

#### Code Scanning

The code within this repository is scanned by Sonarqube (hosted at [SonarCloud](https://sonarcloud.io/)) while commits are being tested. 

This means a `secret` `SONAR_TOKEN` is set within this repository's secrets and that settings may be changed by tuning the [sonar-project.properties file](sonar-project.properties).

### Dependabot

For semi-automatic updates on our dependencies we use GitHub's Dependabot. The settings can be found on
the [dependabot.yml](./.github/dependabot.yml) file.

## Coding

In order to contribute to this project you'll need the following dependencies installed in your machine:

1. `goLang`
2. `minikube` or other k8s distro and its ctl

If you'd rather have an easy time setting up your environment consider using the `.devcontainer` defined in this project.

## Reference

The following websites were queried for the making of this repository:

+ [Go by Example](https://gobyexample.com/)
+ [Awesome Go](https://github.com/avelino/awesome-go)