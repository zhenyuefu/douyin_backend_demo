# 极简抖音后端demo-青训营

## 抖音项目服务端简单示例

项目介绍：https://rdvksaec6w.feishu.cn/docx/doxcnyx7XqSMkzZiJDeVk1pGZfe

项目需要安装FFmpeg

需要添加环境变量:`VIDEO_URL_PREFIX`会作为视频url的前缀保存在数据库中

项目启动命令：

```shell
go build && ./simple-demo
```

也可以使用dockerfile来启动项目

```shell
docker build -t douyin .
```

docker容器使用8080端口，会默认使用host的3306端口的mysql服务，需要映射`app/public`目录为容器的资源保存目录

可以使用docker

### 功能说明

* 用户登录数据保存在mysql中，mysql的地址请修改constants包中的MySQLDefaultDSN
* 视频上传后会保存到本地 public 目录中