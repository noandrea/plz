############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
ARG DOCKER_TAG=0.0.0
ARG DATA_URL=https://opendata.arcgis.com/datasets/273bf4ae7f6a460fbf3000d73f7b2f76_0.csv?outSR=%7B%22latestWkid%22%3A3857%2C%22wkid%22%3A102100%7D
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache curl
# checkout the project 
WORKDIR /builder
COPY . .
# Fetch dependencies.
RUN go get -d -v
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /plz -ldflags="-s -w -extldflags \"-static\" -X main.Version=$DOCKER_TAG"
# download the location file 
RUN curl -L $DATA_URL -o /tmp/data.csv
# build the database
RUN /plz massage --in /tmp/data.csv --out /data.json
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable + data
COPY --from=builder /data.json /
COPY --from=builder /plz /
# Run the whole shebang.
ENTRYPOINT [ "/plz" ]
CMD [ "serve"]
