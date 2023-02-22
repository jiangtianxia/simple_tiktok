<p align="center">
    <img src="images/title.jpeg" alt="Logo">
<p>
  <p align="center">
    <br />
    <a href="https://github.com/jiangtianxia/simple_tiktok.git">查看本项目的文档 »</a>
    <br />
    <a href="https://ozilj01ufe.feishu.cn/base/bascnXyGKEcJi7vOkiVQVlPvoWc?table=tblOc75EZYVXCah0&view=vewMnpNgGD">查看团队工作文档 »</a>
    <br />
  </p>
  </p>

## **项目介绍**

实现极简版抖音服务端

## **目录**
- [上手指南](#上手指南)
    - [开发前的配置要求](#开发前的配置要求)
    - [安装步骤](#安装步骤)
    - [演示界面](#演示界面)
    - [演示视频](#演示视频)
- [文件目录说明](#文件目录说明)
- [开发的整体设计](#开发的整体设计)
   - [服务模块设计](#服务模块设计)
   - [整体的架构图](#整体的架构图)
   - [数据库的设计](#数据库的设计)
   - [Redis架构的设计](#Redis架构的设计)
   - [RocketMQ架构的设计](#RocketMQ架构的设计)
- [使用技术](#使用技术)
- [性能测试](#性能测试)
- [总结与反思](#总结与反思)
   - [目前仍存在的问题](#目前仍存在的问题)
   - [已识别出的优化项](#已识别出的优化项)
   - [架构演进的可能性](#架构演进的可能性)
   - [项目过程中的反思与总结](#项目过程中的反思与总结)
- [参与开源项目](#如何参与开源项目)
- [版本控制](#版本控制)
- [贡献者](#贡献者)
- [鸣谢](#鸣谢)

### **上手指南**

#### 配置要求

1. go 1.13(最低)
2. MySQL(配置文件位于config包中)
3. Redis、RocketMQ环境
4. [抖声客户端app](https://pan.baidu.com/s/194g4bi9ETFWiXEgPM5qDng?pwd=jtiu)

#### 安装步骤
1. 下载源码
2. 配置SSH、FTP、Redis、静态服务器地址等相关参数
3. 启动服务
4. 在客户端app配置相关地址服务端地址

```sh
git clone https://github.com/jiangtianxia/simple_tiktok.git
```
#### 演示界面
**基础功能演示**

//图片文件

**拓展功能演示**

//图片文件

**设置服务端地址**

//图片文件

#### 演示视频

//视频文件

### **文件目录说明**
//目录结构注释

```
simple_tiktok
├── config
├── controller
├── dao
│   ├── mysql
│   └── redis
├── docker
│   ├── mysql
│   ├── nginx
│   ├── redis
│   └── rocketmq
├── docs
├── images
├── logger
├── middlewares
├── models
├── rocketmq
├── router
├── service
├── sql
├── test
├── upload
├── utils
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
└── README.md
```

### **开发的整体设计**

#### 整体的架构图

<p align="center">
    <img src="images/overall_structure_diagram.jpeg">
<p>

#### 数据库的设计

<p align="center">
    <img src="images/sql.jpeg">
<p>

#### Redis架构的设计

<p align="center">
    <img src="images/redis.jpeg">
<p>

#### RocketMQ架构的设计

//图片文件

#### 服务模块设计

###### 用户模块的设计
用户模块包括用户注册、用户登录和用户信息获取三个部分，详情：[用户模块设计说明](https://ozilj01ufe.feishu.cn/docx/P4Asd72jsoQTvQxcDAhcVmlEnJg) 。

###### 点赞模块的设计
点赞模块包括xx。详情：[点赞模块设计说明](https://ozilj01ufe.feishu.cn/docx/Y7Ejd0UyioCGYRxUKMAcoLvXnWf) 。

###### 评论模块的设计
评论模块包括xx。详情：[评论模块设计说明](https://bytedancecampus1.feishu.cn/docx/VhzHd95ccoU0LexA37ic5nNTn4t) 。

###### 视频模块的设计
视频模块包括xx。详情：[视频模块设计说明](https://ozilj01ufe.feishu.cn/docx/PwjidjFklopXYfxCblpcMLa9nud)。

###### 消息模块的设计
消息模块包括xx。详情：[消息模块设计说明](https://ozilj01ufe.feishu.cn/docx/V23UdBH5boF5Dwxjw9qcPMZbnfh) 。

###### 关系模块的设计
关注模块包括xx。详情：[关系模块设计说明](https://ozilj01ufe.feishu.cn/docx/TRald6KOJoJ9tHxaDNBclhSMnQg) 。

### **使用技术**
框架相关：
- [Gin](https://gin-gonic.com/docs/)
- [Gorm](https://gorm.io/docs/)

服务器相关：

//

中间件相关：
- [Redis](https://redis.io/docs/)
- [RocketMQ](https://rocketmq.apache.org/)

数据库：
- [MySQL](https://dev.mysql.com/doc/)

### **性能测试**

//

### **总结与反思**

#### 目前仍存在的问题

//

#### 已识别出的优化项

//

#### 架构演进的可能性

目前项目采用的是单体架构在一定程度上可以满足应用程序的需求，但随着应用程序的不断发展和业务的不断扩展，单体架构也会出现一些问题，比如：

扩展困难：单体架构将整个应用程序封装在一个进程中，很难将应用程序的各个部分进行分离和独立部署，从而影响应用程序的可扩展性。

可维护性下降：单体架构通常是一个庞大的代码库，随着业务逻辑的增加，代码的规模不断扩大，维护难度逐渐增大。

故障扩散：单体架构中的一个故障往往会影响整个应用程序的正常运行，难以快速定位和修复。

未来架构演进的可能性包括：

微服务架构：将单体架构中的各个功能模块拆分成独立的服务，每个服务都运行在自己的进程中，从而实现服务的独立部署和扩展，提高应用程序的可扩展性和可维护性。

容器化架构：将应用程序和它的所有依赖项打包成容器，这些容器可以在任何地方运行，从而提高应用程序的可移植性和部署效率。

无服务架构：将应用程序的功能抽象成函数，通过云服务提供商的无服务器计算平台来运行这些函数，无需关心底层的服务器和基础设施，从而实现弹性扩展和更高的可用性。

总的来说，未来架构演进的目标是提高应用程序的可扩展性、可维护性和可用性，通过使用适合当前业务需求的架构来实现。

#### 项目过程中的反思与总结

1. 安全性有待加强，还需要注意项目的安全性。对于用户上传的视频、图片等敏感数据，需要进行安全的存储和传输，避免数据泄露等安全问题。
 同时，还需要考虑到网络安全、服务器安全等方面的问题，采取相应的措施保证系统的安全性。
2. 在项目开发中，需要注重代码的可读性、可维护性和可扩展性。合理的代码结构、注释和命名，可以让代码更易于理解和维护。在开发过程中，需要考虑到应用可能的扩展和变更，设计可扩展的架构。

### **参与开源项目**

贡献使开源社区成为一个学习、激励和创造的绝佳场所。你所作的任何贡献都是**非常令人感谢**的。

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### **版本控制**

该项目使用Git进行版本管理。您可以在repository参看当前可用版本。

我们欢迎其他的贡献者参与此项目，在这之前，您需要遵循[Git 分支管理规范](https://ypbg9olvt2.feishu.cn/docs/doccnTMRmh7YgMwL2PgZ5moWUsd)和[注释规范](https://juejin.cn/post/7096881555246678046)。

### **贡献者**

- 江泽彬 https://github.com/jiangtianxia
- 裴君辉 https://github.com/paradoxskin
- 刘昕 https://github.com/TomiokapEace
- 张啸宇 https://github.com/TemplarX-boop
- 段欣悦 https://github.com/CynthiaZzzzzzzzz
- 徐龙 https://github.com/longxu0509

您也可以查阅仓库为该项目做出贡献的开发者。

### **鸣谢**

- [字节跳动后端青训营](https://youthcamp.bytedance.com/)
