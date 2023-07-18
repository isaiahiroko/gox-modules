# Origine

## Portr

### Introduction
Portr is a pull-only gateman for kurbenetes cluster. It securely makes call to an external source and make changes to the cluster based on response.

### How It Works
1. Fetch Resources
### Installation
Download the right binary for your platform:
- [Windows](./)
- [Mac](./)
- [Linux](./)

### Usage
The service expose only a single command, which is ran as follows:
```bash
$ go run main.go portr --git-host=https://github.com \
    --git-username=origine-run \
    --git-repo=health \
    --git-password=xxx \
    --git-remote=origin \
    --git-branch=main \
    --git-commit=7bc995e2ccece9c5b40848864ed69044718a262f \
    --docker-host=https://index.docker.io/v1 \
    --docker-username=isaiahiroko \
    --docker-registry=health \
    --docker-password=xxx \
    --docker-version=0.0.0 
```

### [License](./LICENSE.md)
