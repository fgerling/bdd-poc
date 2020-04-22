FROM registry.opensuse.org/opensuse/leap:latest
RUN zypper -n install go git
WORKDIR /app
COPY . .
RUN go get -v -u github.com/cucumber/godog/cmd/godog

ENV PATH="/root/go/bin:${PATH}"

RUN godog -o tester features

ENTRYPOINT ["/app/tester"]
CMD ["/app/features/hello/hello_world.feature"]
