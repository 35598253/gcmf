# Gcmf 使用文档

# 安装


## 直接安装
```
go get -u -v github.com/35598253/gcmf
```

## 推荐使用 go.mod:
```
require github.com/35598253/gcmf latest
```


# 程序编译

通过如下方式编译得到 Admin 可执行程序:
```
# 切换到 Admin 根目录
// 部署 X86
env GOOS=linux go build admin
// M1
env GOOS=linux GOARCH=amd64 go build admin
//Win
env GOOS=windows GOARCH=amd64 go build admin
```
编译成功后会在当前目录生成一个 admin 的可执行文件，默认是没有任何模块的纯后台管理系统，只包含基本设置及会员管理


# `模块` 编译

通过 `go build -tags "module"` 命令生成含有模块的可执行程序:
```
➜   go build -tags "content" // 含有内容模块
➜   go build -tags "content mall" // 含有内容、商城模块

```
