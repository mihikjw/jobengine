FROM golang AS builder
WORKDIR /code
COPY . .
RUN make

#---

FROM ubuntu:latest
WORKDIR /code
RUN mkdir /jobengine
COPY --from=builder /code/bin/jobengine .
ENTRYPOINT ["./jobengine"]