# helloGo

`helloGo` is a repository with quick samples on how to do some common operations and features from `goLang`.

## Status

[![SonarCloud](https://sonarcloud.io/images/project_badges/sonarcloud-orange.svg)](https://sonarcloud.io/summary/new_code?id=rodolphocastro_hellogo)

---

[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=rodolphocastro_hellogo&metric=bugs)](https://sonarcloud.io/summary/new_code?id=rodolphocastro_hellogo)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=rodolphocastro_hellogo&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=rodolphocastro_hellogo)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=rodolphocastro_hellogo&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=rodolphocastro_hellogo)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=rodolphocastro_hellogo&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=rodolphocastro_hellogo)

## Pipelines

We use GitHub Actions to set up and execute our pipelines.

The files controlling each pipeline can be found within the [.gitHub](./.github) repository.

### Commits and Pull Requests

For every commit and pull request the [test.yml](./.github/workflows/test.yml) pipeline is triggered to execute all the
tests within the project.

#### Code Scanning

The code within this repository is scanned by Sonarqube (hosted at [SonarCloud](https://sonarcloud.io/)) while commits are being tested. 

This means a `secret` `SONAR_TOKEN` is set within this repository's secrets and that settings may be changed by tuning the [sonar-project.properties file](sonar-project.properties).

### Dependabot

For semi-automatic updates on our dependencies we use GitHub's Dependabot. The settings can be found on
the [dependabot.yml](./.github/dependabot.yml) file.

## Reference

The following websites were queried for the making of this repository:

+ [Go by Example](https://gobyexample.com/)
