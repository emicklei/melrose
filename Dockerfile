FROM golang:1.14-stretch as build

WORKDIR /app

# get external dependencies
RUN apt-get update
RUN apt-cache policy libportmidi0
RUN apt-cache policy libportmidi-dev
RUN apt-get -y install libportmidi-dev

# allow caching of vendor layer
COPY . .
RUN go build -mod=vendor -o melrose ./cmd/melrose
RUN chmod +x melrose

#FROM scratch
#COPY --from=build /app/melrose /
# docker run -it --device /dev/snd melrose-docker /bin/sh
CMD ["/app/melrose"]