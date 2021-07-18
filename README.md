# httpserver使用指南

## 配置文件

在项目根文件夹下新建`.conf`文件夹，在该文件夹内net.json和db.json的配置如下：

**net.json**：

```json
{
  "BindIP": "xxx.xxx.xxx.xxx",
  "PrivateIP": "xxx.xxx.xxx.xxx",
  "PublicIP": "xxx.xxx.xxx.xxx",
  "Port": 443
}
```

**db.json**：

```json
{
  "SQL": "mysql",
  "IP": "xxx.xxx.xxx.xxx",
  "Port": 3306,
  "User": "xxxx",
  "Passwd": "xxx",
  "DBName": "xxxxxx",
  "MaxOpenConns": 100,
  "MaxIdleConns": 100,
  "Tables": {
    "authtable": {
      "uid": "INT AUTO_INCREMENT PRIMARY KEY",
      "user": "TEXT NOT NULL",
      "passwd": "TEXT NOT NULL"
    }
  }
}
```

## 运行httpserver

- 在根目录下执行命令`make clean`和`make run`在本地运行服务器
- 在根目录下执行命令`sh docker-build.sh`构建Docker Image，然后执行`sh docker-run.sh`在Docker容器内运行服务器