# DEP - Git Multi-Repository Project Manager

## 1. 项目目标

DEP 是一个参考 **uv** 设计理念的 Git 多仓库项目管理工具。

一个项目由多个独立 Git 仓库组成，DEP 提供统一的项目管理能力，包括：

* 管理仓库集合
* 锁定仓库版本（Git Commit）
* 同步整个项目源码
* 查看项目状态
* 在多个仓库执行统一命令

DEP 不管理仓库内部依赖，不关心项目使用的编程语言。

---

# 2. 设计原则

* 一个项目由多个独立 Git 仓库组成
* 仓库必须平级组织
* 每个仓库保持独立 Git 历史
* Git Commit 是唯一版本标识
* 一个项目仅包含一个配置文件
* 一个项目仅包含一个锁文件
* 配置最小化
* 命令最小化
* 行为可预测
* 项目源码可完全复现
* CLI 风格参考 `uv`

---

# 3. 项目结构

```text
workspace/
├── dep.toml
├── dep.lock
├── project-a/
│   └── .git
├── project-b/
│   └── .git
└── project-c/
    └── .git
```

约束：

* 所有仓库必须位于项目根目录
* 仓库之间禁止嵌套
* 本地目录名称由仓库地址自动推导
* 同一项目内目录名称必须唯一

---

# 4. 配置文件（dep.toml）

用于描述项目包含的仓库。

```toml
[[repo]]
url = "git@github.com:org/project-a.git"

[[repo]]
url = "https://github.com/org/project-b.git"
```

字段说明：

| 字段  | 说明       |
| --- | -------- |
| url | Git 仓库地址 |

约束：

* 仓库地址必须以 `.git` 结尾
* 支持 SSH 与 HTTPS
* 不保存仓库名称
* 不保存本地路径
* 不保存分支
* 所有仓库按 URL 字符串排序保存

---

## 仓库地址

仓库地址是仓库的唯一标识。

要求：

* 必须以 `.git` 结尾
* 保留原始格式
* 保留大小写
* 不允许自动补全
* 不允许自动转换

目录名称推导规则：

* 取 URL 最后一个路径段
* 去掉 `.git`

例如：

```
git@github.com:org/project-a.git
                     │
                     └── project-a
```

---

# 5. 锁文件（dep.lock）

记录项目所有仓库对应的源码版本。

```toml
version = 1

[[repo]]
url = "git@github.com:org/project-a.git"
commit = "f8d32418a64b1f4e61d2d5e7e17a3d52fce3d9d2"

[[repo]]
url = "https://github.com/org/project-b.git"
commit = "92ad33e6e7c3c77d9d2f7f0f7e2e8a4c1fd91244"
```

字段说明：

| 字段      | 说明                       |
| ------- | ------------------------ |
| version | 锁文件版本                    |
| url     | Git 仓库地址                 |
| commit  | 完整 Git Commit Hash（40 位） |

## 约束

- Commit 必须为完整 Hash（40 位）
- 所有仓库按 URL 字符串排序保存
- `dep.lock` 由 `dep lock` 命令重建，`dep remove` 命令可删除其中的条目
- `dep lock` 每次执行均重建整个锁文件，不进行增量更新
- `dep add` 与 `dep sync` 不修改锁文件

---

# 6. CLI

## 初始化项目

```bash
dep init
```

生成：

* `dep.toml`
* `dep.lock`

---

## 添加仓库

```bash
dep add <url>
```

示例：

```bash
dep add git@github.com:org/project-a.git
```

执行：

1. 校验 URL 格式
2. 推导目录名称
3. 若 URL 已在 `dep.toml` 中：
   * 本地目录存在 → 报错退出（URL 已存在）
   * 本地目录不存在 → 跳至步骤 4（补克隆）
4. 检测目录冲突
5. Clone 仓库
6. 更新 `dep.toml`

URL 校验：

* 必须以 `.git` 结尾
* 必须能够推导出目录名称

异常：

* URL 不合法
* URL 已在 `dep.toml` 中且本地目录也存在
* 推导目录已存在且不是 Git 仓库
* 推导目录已存在但属于其他仓库

---

## 删除仓库

支持目录名称：

```bash
dep remove project-a
```

或完整 URL：

```bash
dep remove git@github.com:org/project-a.git
```

执行：

* 删除配置
* 删除锁记录
* 默认保留本地仓库

仓库解析顺序：

1. 本地目录名称
2. 完整 URL

URL 必须与配置完全一致。

---

## 锁定版本

```bash
dep lock
```

执行：

* 检查所有仓库状态
* 读取所有仓库 HEAD
* 重建 `dep.lock`

要求：

* 所有仓库工作区必须干净（Clean）
* 任意仓库 Dirty 时立即退出
* 不访问远程仓库
* 不修改仓库内容

---

## 同步项目

```bash
dep sync
```

执行：

* 以 `dep.lock` 为唯一数据源，遍历所有锁记录
* 若仓库在锁文件中但不在 `dep.toml` 中，自动补入 `dep.toml`
* Clone 缺失仓库
* Fetch Git 对象
* Checkout 到锁文件指定 Commit

要求：

* 不升级版本
* 不修改锁文件
* 保证所有仓库与锁文件一致

异常：

* 工作区 Dirty
* Checkout 失败
* 不自动覆盖本地修改

---

## 查看状态

```bash
dep status
```

仅输出异常仓库。

状态：

| 状态      | 说明              |
| ------- | --------------- |
| missing | 仓库不存在           |
| invalid | 不是有效 Git 仓库     |
| dirty   | 工作区存在未提交修改      |
| drifted | 当前 HEAD 与锁文件不一致 |

全部正常：

```text
all synced
```

---

## 执行命令

全部仓库：

```bash
dep exec -- git status
```

指定仓库：

```bash
dep exec project-a -- git status
```

或

```bash
dep exec git@github.com:org/project-a.git -- git status
```

默认串行执行。

并行执行：

```bash
dep exec --parallel -- git status
dep exec -p -- git status
```

仓库解析顺序：

1. 本地目录名称
2. 完整 URL

---

# 7. 命令列表

| 命令                  | 功能    |
| ------------------- | ----- |
| `dep init`          | 初始化项目 |
| `dep add <url>`     | 添加仓库  |
| `dep remove <repo>` | 删除仓库  |
| `dep lock`          | 更新锁文件 |
| `dep sync`          | 同步项目  |
| `dep status`        | 查看状态  |
| `dep exec [flags]`  | 执行命令（`-p` 并行） |

---

# 8. 工作流程

初始化：

```bash
dep init
```

添加仓库：

```bash
dep add git@github.com:org/project-a.git
dep add git@github.com:org/project-b.git
```

开发：

```bash
cd project-a
git pull
git commit
```

锁定版本：

```bash
dep lock
```

其他开发者同步：

```bash
git pull
dep sync
```

恢复整个项目到 `dep.lock` 指定的源码版本。

---

# 9. 设计特点

* 单一配置文件（`dep.toml`）
* 单一锁文件（`dep.lock`）
* Git Commit 作为唯一版本标识
* Git 仓库地址作为唯一仓库标识
* 多仓库平级组织
* 仓库保持独立 Git 历史
* 支持任意编程语言
* 配置最小化
* 命令最小化
* 行为可预测
* 项目源码可完全复现
