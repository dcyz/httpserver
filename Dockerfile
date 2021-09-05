FROM centos

LABEL owner="dcyz.kascas"
# RUN yum clean packages && echo "y" | yum install -y make
COPY . /opt/httpserver
WORKDIR /opt/httpserver
EXPOSE 443
ENTRYPOINT ["./httpserver"]