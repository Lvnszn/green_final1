#FROM python:3.8-alpine3.12
FROM centos:7
RUN mkdir -p /home/admin/atec_project
WORKDIR /home/admin/atec_project
COPY run.sh /home/admin/atec_project/run.sh
COPY server /home/admin/atec_project/server
RUN chmod +x /home/admin/atec_project/server