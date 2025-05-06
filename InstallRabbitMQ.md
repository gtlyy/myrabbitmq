# Install RabbitMQ

## 1. Install
In Fedora:
```
google下载：rabbitmq-server-3.13.2-1.el8.noarch.rpm
然后: sudo dnf install ./rabbitmq-server-3.13.2-1.el8.noarch.rpm
最新版本请前往： https://www.rabbitmq.com/install-rpm.html#downloads

# ps: 这样装：sudo dnf install rabbitmq-server ，不知道为什么启动不了。
```

## 2. Enable Plugins
```
# 这样通过 js 可以访问 web_stomp这样
sudo rabbitmq-plugins enable rabbitmq_management rabbitmq_web_stomp rabbitmq_web_stomp_examples
# add to /etc/rabbitmq/rabbitmq.conf:  web_stomp.tcp.port = 15674
# Ps： 还可以通过web访问：http://127.0.0.1:15672
```

## 3. Start server
```
sudo systemctl start rabbitmq-server
sudo systemctl enable rabbitmq-server
```

## 4. Add user
```
sudo rabbitmqctl add_user userdw01 pq328hu7
# sudo rabbitmqctl change_password userdw01 pq328hu7
sudo rabbitmqctl set_user_tags userdw01 administrator
sudo rabbitmqctl  set_permissions -p / userdw01 '.*' '.*' '.*'
```

## 5. firewall
```
# 获取firewalld当前使用的区域名字
firewallZone=`sudo firewall-cmd --list-all | grep active | cut -d\( -f1`
echo $firewallZone
sudo firewall-cmd --permanent --zone=$firewallZone --add-port=15674/tcp
sudo firewall-cmd --reload
```