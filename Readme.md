# 说明

朋友需要一个简单的上传文件和下载文件的功能，使用 ftp 配置起来
比较麻烦，匿名用户还不可以上传，所以整体效率低下，有了自己服务器
希望能够简单部署起来。

本软件包含了前后端简单的实现
![image](https://user-images.githubusercontent.com/28471249/140696394-3c83479a-c7c0-4d7e-b3d4-243cc8572e44.png)


# 软件构建
## 编译二进制
```
 GOOS=linux GOARCH=amd64 go build -o file file.go
```
## 查看 docker file 

```
➜  awesomeProject git:(master) ✗ cat dockerfile 
FROM busybox
ADD  file /

ENTRYPOINT ["/file"]
```

## 打包镜像 

```
 docker build -t file:latest   .
```

# 软件部署
使用自己打包的镜像或者直接使用我已经发布的镜像

```
mkdir -p file && cd file && docker run -d -v`pwd`:/var/file  -e 'PUB_HOST=xx.xx.xx.xx' -p 9000:9000 quasimodo7017/dev
```
