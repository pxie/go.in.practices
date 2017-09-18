# 使用paho mqtt Go客户端

## 搭建环境
### 连接eclipse测试mqtt服务器，`tcp://iot.eclipse.org:1883`
### 连接本地的docker环境
```
docker pull toke/mosquitto
docker run -ti -p 1883:1883 -p 9001:9001 toke/mosquitto
```

## mqtt客户端使用
1. 先创建mqtt.ClientOptions用来设置连接相关的属性