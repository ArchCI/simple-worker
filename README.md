# Simple worker for ArchCI

## Introduction

Simple-worker is the easy-to-deploy worker to run tests in containers.

## Usage

```
sudo docker run -d --net=host archci/simple-docker
```

* `MYSQL_SERVER` is optional to set address of MySQL(DEFAULT: "")
* `MYSQL_USERNAME` is optional to set MySQL username(DEFAULT: root)
* `MYSQL_PASSWORD` is optional to set user's password(DEFAULT: root)
* `MYSQL_DATABASE` is optional to set MySQL database(DEFAULT: mysql)
* `REDIS_SERVER` is optional to set address of redis(DEFAULT: 127.0.0.1:6379)