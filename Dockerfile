FROM golang:alpine as builder

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	bash \
	ca-certificates

COPY . /go/src/github.com/W1lkins/twitback

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		git \
		gcc \
		libc-dev \
		libgcc \
		make \
	&& cd /go/src/github.com/W1lkins/twitback \
	&& make static \
	&& mv twitback /usr/bin/twitback \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."

FROM alpine:latest

COPY --from=builder /usr/bin/twitback /usr/bin/twitback
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

RUN adduser -D -u 1000 user \
  && chown -R user /home/user

USER user

ENV USER user

WORKDIR /home/user

ENTRYPOINT [ "twitback" ]
CMD [ "--help" ]
