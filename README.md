# 数据准备

通过 docker 下载数据库，并且运行起来
```
docker pull mariadb:5.5.64-trusty
docker run --name mariadb -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mariadb
```

通过连接工具去连接数据库，接下来创建 schema 并且在建表
```
create schema atec2022;

create table total_energy
      (
        id           int auto_increment
        primary key,
        gmt_create   datetime    null,
        gmt_modified datetime    null,
        user_id      varchar(64) null,
        total_energy int         null,
        constraint total_energy_pk
        unique (user_id)
      );
create table to_collect_energy
      (
        id                int auto_increment
        primary key,
        gmt_create        timestamp   null,
        gmt_modified      timestamp   null,
        user_id           varchar(64) null,
        to_collect_energy int         null,
        status            varchar(32) null
      );
```


# 测试参数

原生 net/http 性能 13.8 分左右。

空跑 fasthttp 最好是 22s, 大概 39.41 分
delete + insert 方案 99.39 分（被否定了）

