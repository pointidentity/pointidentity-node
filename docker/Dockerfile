###############################################################
###        STAGE 1: Build node binary pre-requisites        ###
###############################################################

FROM golang:1.22-alpine AS builder

# Install minimum necessary dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers findutils
RUN apk update && apk add --no-cache $PACKAGES

# Set working directory for the build
WORKDIR /go/src/github.com/pointidentity/pointidentity-node

# Add source files
COPY . .

# Make node binary
RUN make tidy && make build-linux

###############################################################
###             STAGE 2: Build pointidentity-node image             ###
###############################################################

FROM alpine:3.18 AS runner

# Install pre-requisites
RUN apk update && apk add --no-cache bash ca-certificates

# Copy over binaries from the build-env
COPY --from=builder /go/src/github.com/pointidentity/pointidentity-node/build/pointidentity-noded /bin/pointidentity-noded

# Set user directory and details
ARG HOME_DIR="/home/pointidentity"
ARG USER="pointidentity"
ARG GROUP=${USER}
ARG UID=1000
ARG GID=${UID}

# Add pointidentity user to use in the container
RUN addgroup --system ${GROUP} --gid ${GID} \
    && adduser ${USER} --uid ${UID} -G ${GROUP} --system --home ${HOME_DIR} --shell /bin/bash

# Set entrypoint script
COPY --chown=${USER}:${GROUP} docker/entrypoint.sh /bin/node-start
RUN chmod +rx /bin/node-start

# Set working directory & bash defaults
WORKDIR ${HOME_DIR}
USER ${USER}
SHELL ["/bin/bash", "-euo", "pipefail", "-c"]

EXPOSE 26656 26657 26660 1317 9090 9091

# Set plain vanilla default entrypoint/command
CMD [ "pointidentity-noded" ]
