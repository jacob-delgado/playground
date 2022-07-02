# Start by building the application.
FROM golang:1.18-bullseye as build

WORKDIR /go/src/github.com/jacob-delgado/playground
ADD . /go/src/github.com/jacob-delgado/playground

RUN GOOS=linux go build ./cmd/playground

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian11
COPY --from=build /go/src/github.com/jacob-delgado/playground/playground /
CMD ["/playground"]
