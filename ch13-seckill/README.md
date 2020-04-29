### 简介
配套《Go语言高并发与微服务实战》使用，具体讲解见图书13章

### 依赖基础组件
- redis
- zookeeper
- git仓库
- consul

#### 部署
- 1 部署 consul 
参考书籍第六章6.3小节内容，安装部署 consul
- 2 部署 Redis,Zookeeper,MySQL。
参考对应组件的官方部署文档，安装完MySQL后，可以导入主目录下的seckill.sql
- 3 新建git repo
可以参考 https://gitee.com/cloud-source/config-repo 创建对应项目的文件，修改Redis，MySQL，Zookeeper等组件的配置
- 4 部署 Config-Service
参考书籍第八章8.3.1小节 在ch8-config文件夹下有 config-service项目，
在yml文件中配置对应的git项目地址和consul地址，构建并运行Java程序，将config-service注册到consul上
- 5 修改bootstrap文件
修改各个项目中的bootstrap.yml文件discover相关的consul地址和config-service的相关配置

