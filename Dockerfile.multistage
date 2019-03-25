FROM golang:1.10.3-alpine as builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN apk add --update bash git curl && \
    (curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh) && \
    git config --global http.https://gopkg.in.followRedirects true

WORKDIR /go/src/github.com/chickenzord/kube-rule
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -v --vendor-only

COPY cmd/ cmd/
COPY pkg/ pkg/
RUN go build -a -o /bin/kuberule-manager ./cmd/manager


FROM alpine:latest
COPY --from=builder /bin/kuberule-manager /bin/kuberule-manager
CMD ["kuberule-manager"]
