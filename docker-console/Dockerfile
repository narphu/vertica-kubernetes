ARG OS_VERSION="7.9.2009"
FROM centos:centos${OS_VERSION}

ARG VERTICA_RPM="vertica-console-x86_64.RHEL6.latest.rpm"

ADD ./packages/${VERTICA_RPM} /tmp/

RUN set -x \
  && yum install -y net-tools \
  && yum localinstall -q -y /tmp/${VERTICA_RPM} \
  && yum clean all

EXPOSE 5450

VOLUME /opt/vconsole/config

RUN /usr/sbin/groupadd -r verticadba \
  && /usr/sbin/useradd -r -m -s /bin/bash -g verticadba dbadmin

ADD ./docker-entrypoint.sh /root/docker-entrypoint.sh

ENTRYPOINT ["bash", "/root/docker-entrypoint.sh"]
