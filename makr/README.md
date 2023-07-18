# Origine

## Makr

### Introduction
A CLI tool for converting source files to container images. 

### Installation
1. Download the right binary for your platform from `assets` directory.
2. Run `install.sh` if your platform is linux/darwin
```bash
$ install.sh 0.0.1 darwin
```

### Usage & Commands
1. Run as a HTTP server
```
$ makr serve --port 25623

Sample request:
POST http://localhost:25623
{
    "git-host": "https://github.com",
    "git-username": "origine-run",
    "git-repo": "health",
    "git-password": "xxx",
    "git-remote": "origin",
    "git-branch": "main",
    "docker-host": "https://index.docker.io/v1",
    "docker-username": "isaiahiroko",
    "docker-registry": "health",
    "docker-password": "xxx",
    "docker-version": "0.0.0 "
}
```

2. Run on CLI
```bash
$ makr run --git-host=https://github.com \
    --git-username=origine-run \
    --git-repo=health \
    --git-password=xxx \
    --git-remote=origin \
    --git-branch=main \
    --docker-host=https://index.docker.io/v1 \
    --docker-username=isaiahiroko \
    --docker-registry=health \
    --docker-password=xxx \
    --image-version=0.0.0 
```

### [License](./LICENSE.md)

### Todo
+ version in yaml env
x delete one checksum
- add mongo db
- add mysql db
- re-design
```
- mark running on k8s
- accept, save and queue request via http
- every x minute, check for pending jobs
- if available:
- - lunch an ec2 instance (the size should be based on previous mem & cpu comsumption for the job type, default to 1vcpu, 2gb)
- - ssh into it
- - setup 
- - clone n install mark
- - run `mark run ...` command (for every pending job)
- - shutdown ec2 if there's not pending job and its 5mins to the end of the hour


   23  sudo apt-get update
   24  sudo apt-get install ca-certificates curl gnupg
   25  sudo install -m 0755 -d /etc/apt/keyrings
   26  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
   27  sudo chmod a+r /etc/apt/keyrings/docker.gpg
   28  echo   "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" |   sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
   29  sudo apt-get update
   30  sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
   32  cd makr/
   38  rm -rf tmp
```
