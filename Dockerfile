FROM centos:7
MAINTAINER tglmm
RUN rm -f /etc/yum.repos.d/* 
RUN yum -y install nginx
ADD CentOS-Base.repo /etc/yum.repos.d/
ADD nginx.repo /etc/yum.repos.d/
RUN echo "this docker of nginx" > /usr/share/nginx/html/index.html
EXPOSE 80
ENTRYPOINT ["/usr/sbin/nginx"]
