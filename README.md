### 说明

修改内容：

1. 👌mysql gorm 强制删除（关闭软连接标志）
2. 👌关闭页面定时轮询未读消息（jwt 鉴权有问题，一直在读数据库）
3. 👌首页可以删除 twitter， 不用进入详细页再删
4. 👌视频定格1s
5. (todo) 增加mysql执行日志配置开关

### 开发

1. 启动 web
   修改 web/.env host 地址
   pnpm run dev

2. 启动 api
   go run main.go serve


icon https://ionic.io/ionicons