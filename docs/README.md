# DNS设计和实现原理
dns服务是为了解决调度器[swan](https://github.com/Dataman-Cloud/swan)七层服务发现的问题。

dns目前具备两个功能:
+ 与代理服务器[swan-janitor](https://github.com/Dataman-Cloud/swan-janitor)配合使用,根据请求的域名,将请求的地址解析为proxy的地址(A).
+ 单独使用时,可以解析出每个应用实例的srv地址(srv).

## 详细设计:
+ 实现原理:

dns通过监听swan的实例更新(创建/删除)事件，将实例的ip:port0存储在srvs中,dig时根据url解析出task_id对应的ip:ports.  
当swan-janitor启动的时候会发送一个事件，将swan-janitor的服务ip通过事件的形式发送给dns，存储在As中, domain提前在dns config中配好。
dig时根据url解析出domain,然后找到domain对应的swan-janitor的ip.

+ 与swan的交互:

dns对外暴露一个叫RecordGeneratorChangeChan的channel,通过接收swan发过来的事件信息来更新维护应用与ip的

# DNS验证方法

+ 范域名支持(A)
```
#dig @DNS_ADDR 0.app1.user1.cluster1.DOMAIN
```
DNS_ADDR为dns服务的ip地址, 例如localhost, 127.0.0.1。

dns的范域名支持与代理服务器proxy配合使用,通过识别DOMAIN，dns可以解析到proxy的地址。
当请求到达proxy时，proxy通过请求的host地址将请求反向代理到实际的服务地址。

+ 服务支持(SRV)
```
#dig @DNS_ADDR 0.app1.user1.cluster1.DOMAIN srv
```
dns会解析出0.app1.user1.cluster1.DOMAIN服务实际的ip,port.

如果dig之后很久没有结果，请查看53端口是否被占用,命令为:
```
#lsof -i -P|grep 53
```
杀掉占用53端口的服务，然后再使用此dns服务。
