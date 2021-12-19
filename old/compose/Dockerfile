# https://blog.jetbrains.com/go/2020/05/08/running-go-applications-using-docker-compose-in-goland/

FROM golang:1.17.5-bullseye

COPY ./compose/entrypoint /entrypoint
RUN sed -i 's/\r//' /entrypoint \
    && chmod +x /entrypoint

ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Insall Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /dockerdev
COPY ./go.mod ./go.sum ./
RUN go mod download -x

ADD . /dockerdev

ENTRYPOINT ["/entrypoint"]