# httpserver使用指南

## 配置文件

在项目根文件夹下新建`.conf`文件夹，在该文件夹内net.json和db.json的配置如下：

**net.json**：

```json
{
  "IP": "192.168.43.234",
  "Port": 443
}
```

**db.json**：

> Tables的内部结构为{表名: {字段名: 类型, 字段名: 类型 ...}, 表名: {...} ...}

```json
{
  "SQL": "mysql",
  "IP": "192.168.43.123",
  "Port": 3306,
  "User": "root",
  "Passwd": "12345678",
  "DBName": "myDB",
  "MaxOpenConns": 100,
  "MaxIdleConns": 100,
  "Tables": {
    
    "table1": {
      "column1": "TEXT",
      "column2": "INT"
    },
    "table2": {
      "column1": "BLOB"
    }
  }
}
```

## 运行httpserver

在根目录下执行命令`make clean`和`make run`即可运行服务器