FROM centos:7
MAINTAINER tglmm
ADD nginx.repo /etc/yum.repos.d/
RUN yum -y install nginx
RUN echo "this docker of nginx" > /usr/share/nginx/html/index.html
EXPOSE 80
ENTRYPOINT ["/usr/sbin/nginx"]
