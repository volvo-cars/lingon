FROM gcr.io/distroless/static-debian11:nonroot-amd64
# ARG BIN=/usr/local/bin
# ARG BUF_VERSION=1.24.0
# ARG KO_VERSION=0.14.1

# RUN apt-get update && apt-get install -y -qq curl

# # Install buf
# RUN curl -sSfL \
#     "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(uname -s)-$(uname -m)" \
#     -o "${BIN}/buf" && \
#     chmod +x "${BIN}/buf"

# # Install ko
# RUN curl -sSfL \
#     https://github.com/ko-build/ko/releases/download/v${KO_VERSION}/ko_${KO_VERSION}_Linux_amd64.tar.gz > ko.tar.gz

FROM golang:1.20 AS build-stage
ARG BUF_VERSION=1.24.0
ARG KO_VERSION=0.14.1

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

# Install ko
RUN go install \
    github.com/google/ko@v${KO_VERSION}
# Install buf
RUN go install \
    github.com/bufbuild/buf/cmd/buf@v${BUF_VERSION}

FROM gcr.io/distroless/static-debian11:nonroot-amd64

COPY --from=build-stage /go/bin/linux_amd64/ko /usr/local/bin/ko
COPY --from=build-stage /go/bin/linux_amd64/buf /usr/local/bin/buf