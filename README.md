# 字节青训营大项目

本项目为你摆你🐴呢队字节青训营结营项目 抖音精简版后端程序



## 项目介绍

本项目使用gRPC框架实现微服务架构的抖音精简版后端程序



## 项目分工

* [@OrionLi](https://github.com/OrionLi)：用户关系相关接口
* [@ygxiaobai111](https://github.com/ygxiaobai111)：用户注册、登录，用户信息接口
* [@OriTsuruHime](https://github.com/OriTsuruHime)：视频流、投稿、发布列表接口
* [@Raqtpie](https://github.com/Raqtpie)：点赞操作、聊天接口
* [@Howe](https://github.com/Cinta-i29)：点赞列表、评论接口

> 运行前需要安装JDK8及以上用于部署Nacos

> 运行前需要ffmpeg



## 启动项目

```bash
cd your_directory/douyin-backend
bash ./start.sh
```

> sh文件中为绝对路径，请自行更改相关路径
>
> 安装jdk环境命令为nix包管理，如为yum或apt请自行更换
>
> 启动nacos的命令请自行更换路径