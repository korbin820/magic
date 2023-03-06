# magic

# zookeeper

+ 创建客户端并获取值，zookeeper.NewManager().Get("/config/test")

+ 客户端连接zk后，默认只存储、监听/connection、/config节点及其子节点，同时监听事件