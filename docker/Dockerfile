ARG GO_VERSION=1.23.2
ARG ALPINE_VERSION=3.20


###
# Builder image 
###
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

# Install packages
RUN apk add --no-cache git make gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod .
COPY go.sum .

# Cache dependencies
RUN go mod download -x

# Copy the rest of source files
COPY .git .
COPY cmd ./cmd
COPY pkg ./pkg
COPY Makefile .

# Build the binary
RUN make build

###
# Runner image
###

FROM alpine:${ALPINE_VERSION}

# Copy binary from builder image
COPY --from=builder /app/build/* /bin/

# Set user directory and details
ARG USER="obfsproxy"
ARG GROUP=${USER}
ARG UID=1000
ARG GID=${UID}
ARG HOME_DIR="/home/${USER}"

# Add user to use in the container
RUN addgroup --system ${GROUP} --gid ${GID} \
    && adduser ${USER} --uid ${UID} -G ${GROUP} --system --home ${HOME_DIR} --shell /bin/bash

# Set working directory & user
RUN mkdir -p ${HOME_DIR} && chown -R ${USER}:${GROUP} ${HOME_DIR}
WORKDIR ${HOME_DIR}
USER ${USER}

# Set entrypoint
CMD ["obfsproxy"]
