FROM golang:1.12

ENV PORT "8080"
ENV MONGO_DB_DATABASE "form3"
ENV MONGO_DB_URI "mongodb://localhost"

WORKDIR /go/src/app
COPY . .

RUN go get .  && go install -v .

ENTRYPOINT ["app"]
