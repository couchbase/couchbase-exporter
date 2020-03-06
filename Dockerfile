FROM scratch

# Add licenses and help file
COPY LICENSE /licenses/LICENSE.txt
COPY README.md /help.1
COPY pkg grafana prometheus /
COPY bin/linux/couchbase-exporter /

ARG PROD_VERSION
ARG PROD_BUILD
ARG OS_BUILD

LABEL name="scratch/couchbase-exporter" \
      vendor="Couchbase" \
      version="${PROD_VERSION}" \
      openshift_build="${OS_BUILD}" \
      exporter_build="${PROD_BUILD}" \
      release="Latest" \
      summary="Couchbase Exporter ${PROD_VERSION}" \
      description="Couchbase Exporter ${PROD_VERSION}" \
      architecture="x86_64" \
      run="docker run --rm couchbase-exporter registry.connect.redhat.com/couchbase/exporter:${PROD_VERSION}-${OS_BUILD} --help"

ENTRYPOINT ["/couchbase-exporter"]
