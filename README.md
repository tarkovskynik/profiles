
# Service from Profiles

Use command: git clone https://github.com/tarkovskynik/profiles.git

**Database.**

For our project we use a relational database MySQL, you can run your database in a Docker container.

use this command: 
```
docker run -d -p 3306:3306 --name mysql-test -e MYSQL_ROOT_PASSWORD=1234 mysql/mysql-server

docker exec -it $(docker ps -q -f name='mysql-test') bash -l;

mysql -uroot -p (enter the password 1234)

CREATE USER 'testdb' IDENTIFIED BY '12345';
GRANT ALL ON *.* TO 'testdb';
CREATE DATABASE dbtest;
FLUSH PRIVILEGES;
```

**Config**

The application depends on config values that are located in the main.go, you can write your path.

**Authorization** :

In Postman you choose "Authorization" menu, type "API KEY"

Key: "Api-key"
Value: "ffff-2918-xcas" or "www-dfq92-sqfwf"

operations:
```
GET - "/profile?username=(enter username)" - get profile by username

GET "/profile" - get all profiles
```