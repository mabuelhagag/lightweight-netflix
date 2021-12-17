FROM golang:1.17.5-bullseye


RUN addgroup --system go && \
    adduser --system --ingroup go --home /opt/project go

COPY --chown=go:go ./compose/entrypoint /entrypoint
RUN sed -i 's/\r//' /entrypoint \
    && chmod +x /entrypoint


ADD --chown=go:go https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

COPY --chown=go:go . /opt/project

USER go
WORKDIR /opt/project

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080
ENTRYPOINT ["/entrypoint"]