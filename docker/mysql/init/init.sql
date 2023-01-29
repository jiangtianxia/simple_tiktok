root@localhost init # cat init.sql 
use mysql;
-- 授权命令，使得该用户可以被任何ip连接
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'tiktokadmin';
ALTER USER 'test'@'%' IDENTIFIED WITH mysql_native_password BY '674092';