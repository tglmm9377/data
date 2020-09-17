FROM centos:7
MAINTAINER tglmm
# RUN rm -f /etc/yum.repos.d/* 
# COPY CentOS-Base.repo /etc/yum.repos.d/
# COPY nginx.repo /etc/yum.repos.d/
# RUN yum -y install nginx
# RUN echo "this docker of nginx" > /usr/share/nginx/html/index.html
RUN echo `date` >> /opt/ip.txt
# EXPOSE 80
# ENTRYPOINT ["/usr/sbin/nginx"]
