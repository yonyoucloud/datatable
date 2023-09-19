# 数据表格（datatable）

数据库`查询展示`工具。

### 1、构建前端

```sh
make install-fe build-fe
```

### 2、构建后端

```sh
make build-be
```

### 3、修改配置

```sh
cd ./deploy
cp -rp config-example.yaml config.yaml

# 修改 config.yaml
mysql.master 修改正确数据库链接信息
```

### 4、运行服务

```sh
npm run
```

### 5、生成前端配置文件，访问如下地址

```sh
# host 参数为 api 接口地址，可以设置为当前服务访问地址
curl http://${IP}:8889/api/v1/create/config/js?host=http://${IP}:8889/
```

### 6、系统入口地址

```sh
http://${IP}:8889/web
```