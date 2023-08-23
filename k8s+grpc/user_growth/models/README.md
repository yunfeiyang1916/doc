# 安装xorm，生成数据模型(./models)
````
# 注意：这里的目录，需要根据本地的xorm安装目录，以及本地的user_growth项目目录做修改。
cd /var/www/go/src/gitea.com/xorm/reverse/example
reverse -f mysql-usergrowth.yml
cp ../models/user_growth/models.go ~/Documents/imooc/user_growth/models
# 注意：生成后的models内容，主键id是int类型，但是文件中是string，需要手动修改一下
````