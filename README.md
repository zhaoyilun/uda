# UDA

UDA（Unix-like Python Environment Manager）是一个用 Go 实现的轻量级 Python 环境管理工具，核心目标是用 CLI 封装 `uv`，提供类似 Conda 的环境管理体验。

## 目录
- [快速开始](#快速开始)
- [安装与配置](#安装与配置)
- [命令参考](#命令参考)
- [项目结构](#项目结构)
- [环境与配置](#环境与配置)
- [开发与验证](#开发与验证)

## 快速开始

```bash
# 克隆项目
git clone https://github.com/zhaoyilun/uda.git
cd uda

# 编译
./scripts/build.sh

# 初始化 shell
eval "$(./uda init bash)"

# 创建并激活环境
uda create myenv --python 3.11
uda activate myenv
```

## 安装与配置

- 先决条件：Go 1.22+（源码编译）、Git。
- 初始化时会创建 `~/.uda`、`~/.uda/envs`、`~/.uda/cache`。
- 安装 uv（首次推荐）：
```bash
uda self install
```
- 配置镜像：可通过环境变量或配置文件设置
```bash
export UV_MIRROR=https://pypi.tuna.tsinghua.edu.cn
```
或写入 `~/.uda/config.toml`。

## 命令参考

```bash
uda create <name> [--python 3.11]   # 创建环境
uda list                             # 列出环境
uda remove <name>                    # 删除环境
uda activate <name>                  # 激活环境（输出 shell 片段）
uda deactivate                       # 退出环境
uda install --env <name> pkg1 pkg2   # 安装依赖
uda run --env <name> <command>       # 在指定环境执行命令
uda run <command>                   # 未指定 env 时使用当前环境
uda self install                     # 安装/更新 uv
uda init [bash|zsh|fish]            # 输出 shell 集成脚本
```

## 项目结构

- `cmd/`：CLI 命令定义
- `internal/config/`：配置路径与目录初始化
- `internal/env/`：环境创建/查询/删除
- `internal/uv/`：uv 下载与调用
- `internal/shell/`：`init/activate/deactivate` 脚本
- `docs/`：设计、计划与完整文档
- `scripts/build.sh`：发布构建脚本

## 开发与验证

建议验证流程：

```bash
go test ./...            # 运行测试（当前仓库若无测试可正常为无测试输出）
go build -o uda .         # 验证源码可编译
./uda --help             # 验证 CLI 命令列表
```

## 贡献

- 遵循仓库提交规范：`feat/fix/chore/docs` 等前缀。
- 新增逻辑建议补充 `_test.go`。
- 如有行为变更，请在 PR 中附带命令输出。
