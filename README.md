<p align="center">
  <img src="https://img.shields.io/badge/go-1.21%2B-blue">
  <img src="https://img.shields.io/badge/license-MIT-green">
</p>

# DEP — Git Multi-Repository Project Manager

DEP manages a collection of independent Git repositories as a single project. `dep sync` restores every local repository to the exact commit recorded in the lock file, enabling reproducible project-level checkouts.

[中文版](README.zh.md)

## Install

Requires: Go 1.21+, Git

### Option 1: Download from Releases (recommended)

Download the pre-built binary for your platform from the [Releases](https://github.com/gaamingcan/DEP/releases) page:

```bash
# Example: Linux amd64
curl -LO https://github.com/gaamingcan/DEP/releases/download/v0.2.0/dep-0.2.0-linux-amd64.tar.gz
tar xzf dep-0.2.0-linux-amd64.tar.gz
sudo cp dep /usr/local/bin/
```

### Option 2: Build from source

```bash
git clone https://github.com/gaamingcan/DEP && cd DEP
go build -o dep .
sudo cp dep /usr/local/bin/
```

## Quick Start

```bash
# Initialize the project
dep init

# Add repositories
dep add git@github.com:org/project-a.git
dep add git@github.com:org/project-b.git

# Lock HEAD commits and commit the lock file
dep lock
git add dep.lock && git commit -m "lock versions"

# Push to remote
git push
```

Other developers sync the project:

```bash
git pull
dep sync
```

All repositories are now checked out to the exact commits recorded in `dep.lock`.

## Commands

| Command | Description |
|---------|-------------|
| `dep init` | Initialize a project in the current directory |
| `dep add <url>` | Add and clone a repository |
| `dep remove <repo>` | Remove a repository (keep local files) |
| `dep lock` | Lock all HEAD commits to dep.lock |
| `dep sync` | Restore all repositories from dep.lock |
| `dep status` | Show abnormal repository states |
| `dep exec -- <cmd>` | Run a command in every repository |
| `dep exec -p -- <cmd>` | Run in parallel (`-p` or `--parallel`) |

## Files

| File | Description |
|------|-------------|
| `dep.lock` | Project management file (repository list + locked commits), commit to VCS |

```toml
version = 1

[[repo]]
url = "git@github.com:org/project-a.git"

[[repo]]
url = "git@github.com:org/project-b.git"
commit = "f8d32418a64b1f4e61d2d5e7e17a3d52fce3d9d2"
```

## Principles

- Repositories are peers, each keeps its own Git history
- Git commit is the sole version identifier
- Minimal configuration — only 7 commands
- dep.lock is committed to version control, ensuring team-wide consistency
