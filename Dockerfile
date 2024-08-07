FROM centos:8
LABEL maintainer="734839030@qq.com"
#  docker build --build-arg MODULE_NAME -t [image:version] .
ARG MODULE_NAME
WORKDIR /data/${MODULE_NAME}

# for shell to choose environment
ENV IN_CONTAINER=true
# just src content
COPY target/${MODULE_NAME} /data/${MODULE_NAME}/

# 方便docker -P参数
#EXPOSE 8080
ENTRYPOINT ["bin/docker-entrypoint.sh"]
CMD ["start"]