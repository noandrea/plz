# PLZ Rest API

A simple rest API that exposes data related to zip codes and buildings in Berlin.

The source dataset is published by [Esri](https://www.esri.de/de-de/home)
and is available [here](https://opendata-esri-de.opendata.arcgis.com/datasets/273bf4ae7f6a460fbf3000d73f7b2f76_0).

## Motivations

Built as an exercise

## Rest API

PLZ provides the following endpoints:

- `/status` returns the status of the service
- `/zip/buildings` returns the number of buildings aggregated by zip code
- `/zip/buildings/history` returns the number of buildings aggregated by zip code and year
- `/zip/buildings/:code` returns the number of buildings for a specific zip code
- `/zip/buildings/:code/history` returns the number of buildings aggregated by zip code and year for a specific zip code

## Usage

There are 2 ways to run the PLZ api service: using [Docker](#docker)(recommended) or via [manual setup](#manual-setup).

> For more details about the command line options, use the command `plz --help`

### Docker

The Docker image is available at [noandrea/plz](https://hub.docker.com/repository/docker/noandrea/plz), and can be run with

```sh
docker run -p 2007:2007 noandrea/plz
```

The image is built on [scratch](https://hub.docker.com/_/scratch), the image size is ~9.3mb:

[![asciicast](https://asciinema.org/a/350265.svg)](https://asciinema.org/a/350265?autoplay=1)

### Manual setup

Those are the steps to setup the service:

1. Install `plz`

```
go get github.com/noandrea/plz
```

**OR**

Download the latest from the [release page](https://github.com/noandrea/plz/releases)

2. Download the dataset linked above:

```sh
curl -L https://opendata.arcgis.com/datasets/273bf4ae7f6a460fbf3000d73f7b2f76_0.csv?outSR=%7B%22latestWkid%22%3A3857%2C%22wkid%22%3A102100%7D -o /tmp/src.csv
```

3. Massage the dataset to produce an optimized json to be served via the Rest API

```sh
plz massage --input /tmp/src.csv --output rest.json
```

4. Run the Rest API service

```sh
plz serve --data rest.json
```

[![asciicast](https://asciinema.org/a/350262.svg)](https://asciinema.org/a/350262?t=63&autoplay=1)

## Examples

### Docker compose

`docker-compose.yaml` example

```yaml
version: '3'
services:
  plz:
    container_name: plz
    image: noandrea/plz:latest
    ports:
    - 2007:2007

```


### K8s

Kubernetes configuration example:

```yaml
---
# Deployment
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    app: plz
  name: plz
spec:
  replicas: 1
  revisionHistoryLimit: 3
  selector:
    matchLabels:
      app: plz
  template:
    metadata:
      labels:
        app: plz
    spec:
      containers:
      - env:
        image: noandrea/plz:latest
        imagePullPolicy: Always
        name: plz
        ports:
        - name: http
          containerPort: 2007
        livenessProbe:
          httpGet:
            path: /status
            port: 2007
---
# Service
# the service for the above deployment
apiVersion: v1
kind: Service
metadata:
  name: plz-service
spec:
  type: ClusterIP
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: http
  selector:
    app: plz

```
