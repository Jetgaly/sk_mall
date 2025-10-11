#!/bin/bash

# 数据库连接信息
DB_USER="jet"
DB_PASS="Qwe@123456"
DB_NAME="sk_merchant_db"
DB_HOST="127.0.0.1"
DB_PORT="3306"
SQL_FILE="/home/jjet/projects/GoProjects/sk_mall/scripts/sql/sk_merchant_db/create_tables.sql"

if [ ! -f "$SQL_FILE" ]; then
    echo "sql file does not exist :$SQL_FILE"
    exit 1
fi

echo "starting to connect DB and exec sql script"
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASS $DB_NAME < $SQL_FILE 

if [ $? -eq 0 ]; then
  echo "create merchant tables successfully"
else
  echo "create merchant tables unsucessfully"
  exit 1
fi