FROM node:10 AS build_static
COPY vue /usr/src
WORKDIR /usr/src
RUN npm install
RUN npm run build

FROM golang
COPY go /go/src
WORKDIR /go/src/main
RUN go get
RUN go build
COPY --from=build_static /usr/src/dist .
RUN cp config.prod.json /go/bin 
ENV SPOTIFY_ID $SPOTIFY_ID
ENV SPOTIFY_SECRET $SPOTIFY_SECRET
CMD /go/bin/main
EXPOSE 443