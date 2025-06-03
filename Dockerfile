# checkov:skip=CKV_DOCKER_3
# checkov:skip=CKV_GHA_7
# trivy:ignore:AVD-DS-0002
# trivy:ignore:AVD-DS-0026

# CKV_DOCKER_3 # Ensure that a user for the container has been created
# CKV_GHA_7    # The build output cannot be affected by user parameters...
# AVD-DS-0002 (HIGH): Specify at least 1 USER command in Dockerfile with non-root user as argument
# AVD-DS-0026 (LOW):  Add HEALTHCHECK instruction in your Dockerfile

FROM scratch
COPY compat /
ENTRYPOINT ["/compat"]
