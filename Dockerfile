FROM alpine:latest

EXPOSE 8080 8090 8081 8091

COPY ./build/linux/simple_tiktok /app/simple_tiktok
COPY ./apidocs/swagger.json /app/apidocs/swagger.json
COPY ./upload /upload

WORKDIR /app

CMD /app/simple_tiktok