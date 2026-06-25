<p align="center">
  <img src="https://img.shields.io/badge/go-1.21%2B-blue">
  <img src="https://img.shields.io/badge/license-MIT-green">
</p>

# DEP — Git Multi-Repository Project Manager

DEP 管理一组独立 Git 仓库，将其视为一个项目。`dep sync` 可将本地仓库恢复到锁文件指定的精确版本，实现项目级源码可复现。

[English](README.md)

## 安装

依赖：Go 1.21+、Git

### 方式一：从 Releases 下载（推荐）

从 [Releases 页面](https://github.com/gaamingcan/DEP/releases) 下载对应平台的预编译二进制：

```bash
# 示例：Linux amd64
curl -LO https://github.com/gaamingcan/DEP/releases/download/v0.2.0/dep-0.2.0-linux-amd64.tar.gz
tar xzf dep-0.2.0-linux-amd64.tar.gz
sudo cp dep /usr/local/bin/
```

### 方式二：从源码构建

```bash
git clone https://github.com/gaamingcan/DEP && cd DEP
go build -o dep .
sudo cp dep /usr/local/bin/
```

## 快速开始

```bash
# 初始化项目
dep init

# 添加仓库
dep add git@github.com:org/project-a.git
dep add git@github.com:org/project-b.git

# 锁定 HEAD 并提交锁文件
dep lock
git add dep.lock && git commit -m "lock versions"

# 推送到远程
git push
```

其他开发者同步项目：

```bash
git pull
dep sync
```

所有仓库已被恢复到 `dep.lock` 记录的精确版本。

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
| `dep.lock` | 项目管理文件（仓库列表 + 锁定版本），需提交到版本控制 |

## 理念

- 仓库平级，保持独立 Git 历史
- Git Commit 是唯一版本标识
- 配置最小化，仅 7 条命令
- dep.lock 提交到版本控制，保障全团队一致
