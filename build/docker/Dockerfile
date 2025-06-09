FROM node:24-bookworm AS node-builder
WORKDIR /workspace
ENV NODE_ENV=production
ENV ENABLE_EXPERIMENTAL_COREPACK=1
ENV YARN_VERSION=4.5.0
COPY ./app/ /workspace/
RUN corepack enable && corepack prepare
RUN yarn install --immutable
RUN yarn build

FROM golang:1.24.3 AS builder
WORKDIR /workspace
ENV CGO_ENABLED=0
ENV GOOS=linux
COPY . .
COPY --from=node-builder /workspace/build/ /workspace/pkg/webclient/build/
RUN ls -al pkg/webclient
RUN go build -a -installsuffix cgo -o runtime_binary ./cmd/server

# Production stage
FROM scratch

# Build args
ARG COMMIT_SHA
ENV COMMIT_SHA="${COMMIT_SHA:-No COMMIT_SHA}"
ENV WEBCLIENT_APP_BUILD_DIR="/build"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /workspace/runtime_binary ./
COPY --from=node-builder /workspace/build/ /build/
EXPOSE 8080
ENTRYPOINT ["./runtime_binary"]
