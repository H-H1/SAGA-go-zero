#!/usr/bin/env bash

# 使用方法：
# ./genModel.sh usercenter user
# ./genModel.sh usercenter user_auth
# 再将./genModel下的文件剪切到对应服务的model目录里面，记得改package


#生成的表名
tables=testtable
#表生成的genmodel目录
modeldir=./gen

# 数据库配置
host=127.0.0.1
port=33069
dbname=t1111
username=root
passwd=PXDN93VRKUm8TeE7


# echo "开始创建库：$dbname 的表：$2"
# goctl model mysql datasource -url="${username}:${passwd}@tcp(${host}:${port})/${dbname}" -table="${tables}"  -dir="${modeldir}" -cache=true --style=goZero


goctl model mysql datasource -url="root:PXDN93VRKUm8TeE7@tcp(127.0.0.1:33069)/dtm_gozero" -table="stock"  -dir="./gen"  --style=goZero
