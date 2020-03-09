FROM registry.access.redhat.com/rhel7

# Install latest security updates
RUN yum repolist --disablerepo="*" && \
    yum-config-manager --enable rhel-7-server-rpms rhel-7-server-rt-rpms && \
    yum-config-manager --enable rhel-7-server-optional-rpms \
        rhel-7-server-extras-rpms rhel-server-rhscl-7-rpms && \
    yum -y update-minimal --security --sec-severity=Important --sec-severity=Critical \
        --setopt=tsflags=nodocs \
    yum clean all

# Add licenses and help file
RUN mkdir /license
COPY LICENSE /licenses/LICENSE.txt
COPY README.md /help.1

ARG PROD_VERSION
ARG PROD_BUILD
ARG OS_BUILD

# Install Couchbase Exporter
COPY bin/linux/couchbase-exporter /usr/local/bin/couchbase-exporter

LABEL name="rhel7/couchbase-exporter" \
      vendor="Couchbase" \
      version="${PROD_VERSION}" \
      openshift_build="${OS_BUILD}" \
      exporter_build="${PROD_BUILD}" \
      release="Latest" \
      summary="Couchbase Exporter ${PROD_VERSION}" \
      description="Couchbase Exporter ${PROD_VERSION}" \
      architecture="x86_64" \
      run="docker run --rm couchbase-exporter registry.connect.redhat.com/couchbase/exporter:${PROD_VERSION}-${OS_BUILD} --help"

ENTRYPOINT ["couchbase-exporter"]