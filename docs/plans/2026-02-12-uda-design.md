# UDA 设计文档

## 项目概述

**项目名称**: uda
**项目类型**: Python 环境管理工具
**核心特性**: 结合 Conda 的全局环境管理和 UV 的快速安装能力
**目标用户**: 喜欢折腾、需要在任意目录快速切换 Python 环境的开发者

## 技术栈

- **语言**: Go
- **核心依赖**: uv (通过 CLI 调用)
- **CLI 框架**: urfave/cli 或 cobra

## 核心架构

```
┌─────────────────────────────────────────┐
│              uda CLI (Go)               │
├─────────────────────────────────────────┤
│  命令层 (cobra/urfave)                  │
├─────────────────────────────────────────┤
│  UV 调用层                               │
│  - uv 下载/更新 (自动镜像)               │
│  - uv python install                    │
│  - uv venv                              │
│  - uv pip install                       │
├─────────────────────────────────────────┤
│  镜像管理                                │
│  - 内置常用镜像列表                       │
│  - 自动检测/切换                         │
├─────────────────────────────────────────┤
│  环境存储 (~/.uda/)                      │
│  - uv (uv 二进制)                       │
│  - envs/<env_name>/                    │
│  - cache/                               │
│  - config.toml                          │
└─────────────────────────────────────────┘
```

## 命令设计

| 命令 | 功能 |
|------|------|
| `uda create <name> --python 3.11` | 创建环境 + 自动下载 Python |
| `uda activate <name>` | 激活环境 |
| `uda deactivate` | 退出环境 |
| `uda list` | 列出所有环境 |
| `uda remove <name>` | 删除环境 |
| `uda install <packages>` | 安装包 |
| `uda run <cmd>` | 在环境中运行 |
| `uda self install` | 安装/更新 uv (自动镜像) |
| `uda init [shell]` | 初始化 shell 插件 |

## 镜像策略

- **内置镜像**: 常用国内镜像 (清华、阿里等)
- **自动检测**: 先尝试官方源，失败后自动切换镜像
- **手动配置**: 支持 `UV_MIRROR` 环境变量或 `~/.uda/config.toml`

## 激活方式

```bash
# 初始化 shell
eval "$(uda init zsh)"

# 之后可以像 Conda 一样使用
uda activate myenv
conda deactivate  # 兼容 Conda 别名
```

## 数据存储

- **根目录**: `~/.uda/`
- **环境目录**: `~/.uda/envs/<env_name>/`
- **uv 二进制**: `~/.uda/uv`
- **配置文件**: `~/.uda/config.toml`

## 实现优先级

1. **MVP**:
   - `uda create` - 创建环境
   - `uda list` - 列出环境
   - `uda remove` - 删除环境
   - `uda self install` - 安装 uv

2. **核心功能**:
   - `uda activate` / `uda deactivate`
   - `uda init` - shell 集成

3. **包管理**:
   - `uda install`
   - `uda run`

4. **增强功能**:
   - Python 版本管理
   - 镜像自动切换
   - 环境克隆/导出
