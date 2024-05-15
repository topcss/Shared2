# Shared2 - 最简单的局域网文件共享服务 🌐

🎉 欢迎使用Shared2！我们提供简单易用的局域网文件共享服务。

🚀 无需登录，快速启动，立即共享文件。

🔍 尝试过 Jupyter Lab、Nextcloud 等服务？我们做得更简单，更实用。

🤝 局域网的协作，基于信任，无需复杂的账户系统。

⚙️ 只需在一个系统上启动服务，指定文件夹，设置大小，即可开始。

🌍 支持主流操作系统：Windows，Amd64 Linux 和 Arm64 Linux。

## 使用 🚀

1. 📥 下载 `shared2.exe` 或 `shared2`。
2. 📂 放入 PATH 环境变量中。
3. 🎯 启动服务，指定文件夹，设置大小。

``` shell
// 显示帮助
shared2 -h
 
// 设置权限（仅限 Linux）
chmod 777 shared2

// 启动服务，指定端口和限制数据目录大小
./shared2 -p 9002 -s 1G &

// 查找进程 ID (PID)
lsof -i :9002

// 终止进程
kill -9 <PID>
```

## 构建 🛠️

### Windows
🖥️ Windows 用户？一行命令，即可构建：
```
go build -ldflags="-s -w" -o shared2.exe
```

### Linux
🐧 Linux 用户？设置环境变量，运行构建命令：
```
SET CGO_ENABLED=0
set GOARCH=amd64
set GOOS=linux
go build -ldflags="-s -w" -o shared2
```

## 贡献 🙌

👀 有改进 Shared2 的建议？或者发现错误？欢迎开启 issue！

👐 我们欢迎所有帮助我们改进 Shared2 的贡献。

## 许可证 📝

[MIT](https://choosealicense.com/licenses/mit/)

📚 请注意，参与本项目即表示你同意遵守我们的贡献者行为准则。