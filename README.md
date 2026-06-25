<p align="center">
  <img src="https://img.shields.io/badge/go-1.21%2B-blue">
  <img src="https://img.shields.io/badge/license-MIT-green">
</p>

# DEP — Git Multi-Repository Project Manager

DEP 管理一组独立 Git 仓库，将其视为一个项目。开发者使用 `dep sync` 即可将本地仓库恢复到锁文件指定的精确版本，实现项目级源码可复现。

## 安装

依赖：Go 1.21+、Git

```bash
# 此命令安装到 `$HOME/go/bin/dep`，需将 `$HOME/go/bin` 加入 `PATH`
go install github.com/org/dep@latest

or

git clone <repo-url> && cd dep
go build -o dep 
sudo cp dep /usr/local/bin/ or /usr/bin/
```


## 快速开始

```bash
# 初始化项目
dep init

# 添加仓库
dep add git@github.com:org/project-a.git
dep add git@github.com:org/project-b.git

# 锁定版本并提交锁文件
dep lock
git add dep.lock && git commit -m "lock versions 1.0.0.1"

# 提交配置（锁文件需 git push 到远程）
git push
```

其他开发者同步项目：

```bash
git pull
dep sync
```

此时所有仓库被同步到 `dep.lock` 记录的精确版本。

## 命令

| 命令 | 功能 |
|------|------|
| `dep init` | 在当前目录初始化项目 |
| `dep add <url>` | 添加并 Clone 仓库 |
| `dep remove <repo>` | 移除仓库（保留本地文件） |
| `dep lock` | 锁定所有仓库 HEAD 到 dep.lock |
| `dep sync` | 按 dep.lock 恢复所有仓库 |
| `dep status` | 显示异常仓库 |
| `dep exec -- <cmd>` | 在所有仓库执行命令 |
| `dep exec -p -- <cmd>` | 并行执行命令（`-p` 或 `--parallel`） |

## 文件

| 文件 | 说明 |
|------|------|
| `dep.toml` | 项目仓库列表 |
| `dep.lock` | 仓库版本锁定文件，需提交到版本控制 |

## 理念

- 仓库平级，保持独立 Git 历史
- Git Commit 是唯一版本标识
- 配置最小化，仅 7 条命令
- dep.lock 提交到版本控制，保障全团队一致

