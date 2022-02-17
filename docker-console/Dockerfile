ARG OS_VERSION="7.9.2009"
FROM centos:centos${OS_VERSION}

ARG VERTICA_RPM="vertica-console-x86_64.RHEL6.latest.rpm"

ADD ./packages/${VERTICA_RPM} /tmp/

RUN set -x \
  && yum install -y net-tools \
  && yum localinstall -q -y /tmp/${VERTICA_RPM} \
  && yum clean all

EXPOSE 5450

# We will setup a PV for the /config and /mcdb directories.  We need to save
# the contents of those so that we have something to seed when we start the
# container the first time.
RUN set -x \
  && mkdir -p /tmp/pv-seed \
  && cp -r /opt/vconsole/config /tmp/pv-seed \
  && cp -r /opt/vconsole/mcdb /tmp/pv-seed

RUN /usr/sbin/groupadd -r verticadba \
  && /usr/sbin/useradd -r -m -s /bin/bash -g verticadba dbadmin

ENV PATH "$PATH:/opt/vconsole/vendor/OpenJDK/Zulu/java/1.8u212/jre/bin"

ADD ./docker-entrypoint.sh /root/docker-entrypoint.sh
ADD ./packages/MCClient.jar /jars/MCClient.jar

ENTRYPOINT ["bash", "/root/docker-entrypoint.sh"]
