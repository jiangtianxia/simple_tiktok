
FROM nginx:latest

# 暴露端口号
EXPOSE 80 8000 8080

# # 删除原有配置文件，创建静态资源文件夹和动态资源文件夹
# RUN rm /etc/nginx/conf.d/default.conf \
# && mkdir -p /usr/share/nginx/html/static \
# && mkdir -p /usr/share/nginx/html/media

# # 设置media文件夹用户和文件组为linux默认www-data,并给予可读可执行权限
# # 否则用户上传的图片无法正确显示
# RUN chown -R www-data:www-data /usr/share/nginx/html/media \
# && chmod -R 775 /usr/share/nginx/html/media

# # 创建前端目录,并导入网页
# RUN mkdir -p /usr/share/nginx/html/treehole
# COPY ./html/ /usr/share/nginx/html/

# 添加配置文件
ADD ./conf.d/ /etc/nginx/conf.d/

# 关闭守护模式
CMD ["nginx", "-g", "daemon off;"]