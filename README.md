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

### Docker

The Docker image is available at [noandrea/plz](https://hub.docker.com/repository/docker/noandrea/plz), and can be run with

```sh
docker run -p 2007:2007 noandrea/plz
```

The image is built on [scratch](https://hub.docker.com/_/scratch), the image size is ~9.3mb:

[![asciicast](https://asciinema.org/a/350213.svg)](https://asciinema.org/a/350213)

### Manual setup

There are 3 steps to setup the service:

1. Download the dataset linked above:

```sh
curl -L https://opendata.arcgis.com/datasets/273bf4ae7f6a460fbf3000d73f7b2f76_0.csv?outSR=%7B%22latestWkid%22%3A3857%2C%22wkid%22%3A102100%7D -o /tmp/src.csv
```

2. Massage the dataset to produce an optimized json to be served via the Rest API

```sh
plz massage --input /tmp/src.csv --output rest.json
```

3. Run the Rest API service

```sh
plz serve --data rest.json
```

[![asciicast](https://asciinema.org/a/350219.svg)](https://asciinema.org/a/350219)

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
