# 构建发布版本

## 前提条件

- Go 1.22+
- Git

## 构建命令

```bash
# 克隆仓库
git clone https://github.com/your-username/uda.git
cd uda

# 构建
go build -ldflags="-s -w" -o uda .

# 或使用构建脚本
./scripts/build.sh
```

## 一键安装

### Linux/macOS

```bash
# 方法 1: 使用安装脚本
curl -sSf https://raw.githubusercontent.com/your-username/uda/main/install.sh | sh

# 方法 2: 直接下载二进制
curl -fsSL https://github.com/your-username/uda/releases/latest/download/uda-linux-x86_64 -o ~/.local/bin/uda
chmod +x ~/.local/bin/uda
```

### Windows

```powershell
# 使用 PowerShell
irm https://github.com/your-username/uda/releases/latest/download/uda-windows-x86_64.exe -o $env:LOCALAPPDATA\uda\uda.exe
```

## 安装后配置

```bash
# 添加到 PATH
export PATH="$HOME/.local/bin:$PATH"

# 初始化 shell（添加到 ~/.bashrc 或 ~/.zshrc）
eval "$(uda init bash)"

# 安装 uv
uda self install

# 创建环境
uda create myenv --python 3.11

# 激活环境
uda activate myenv
```

## 使用镜像安装

如果网络访问 GitHub 困难，可以设置镜像：

```bash
# 使用环境变量
export UV_MIRROR=https://pypi.tuna.tsinghua.edu.cn

# 或配置文件
echo 'mirror = { url = "https://pypi.tuna.tsinghua.edu.cn" }' > ~/.uda/config.toml
```
