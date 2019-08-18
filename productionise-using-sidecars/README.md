# Productionise any Application within Kubernetes using Sidecar Containers

# Table of Contents

- [Productionise any Application within Kubernetes using Sidecar Containers](#productionise-any-application-within-kubernetes-using-sidecar-containers)
- [Table of Contents](#table-of-contents)
  - [Abstract](#abstract)
  - [Demo](#demo)
    - [Commands](#commands)
    - [Improvements](#improvements)
      - [Monitoring - Sidecar Pattern](#monitoring---sidecar-pattern)
      - [Logging - Adapter Pattern](#logging---adapter-pattern)
      - [Circuit Breaker - Ambassador Pattern](#circuit-breaker---ambassador-pattern)

## Abstract

The sidecar pattern in Kubernetes allows you to add other functionality
alongside an application for things like monitoring, TLS termination,
circuit breaking etc. This talk will demonstrate how to easily improve
the security and reliability of a service that you don’t have access
to the code base.

## Demo

### Commands

```
▶ make help

 Choose a command run:

  install_helm             Installs the Helm Tiller on Cluster
  install_prometheus       Installs the Prometheus Operator
  apply_prometheus_rules   Applies recording and alerting rules
  install_nginx            Installs the Nginx Ingress

  forward_prometheus       Port Forwards Prometheus Server
  open_grafana             Open Grafana in Browser
  images                   creates dependant docker images
  create_load              create some load on the service

  initial                  installs the initial version of the services
  v1                       installs the v1 version of the services
  v2                       installs the v2 version of the services
  v3                       installs the v3 version of the services
  reset_demo               deletes all demo resources

  help                     prints this help message
```

This demo code is centered around 2 services, a frontend and a dependant API.

![initial arch][./docs/initial-architecture.png]

To start deploy the dependant services:

```
$ make install_prometheus
$ make install_nginx
```

Then you'll want to deploy the initial version of the services:

```
make initial
```

Navigate to [http://localhost](http://localhost) to see the application.

![frontend app ok][./docs/frontend-app-ok.png]

Refresh a couple times and you will see a 500 failure

![frontend app failure][./docs/frontend-app-fail.png]

### Improvements

Now that we have the app up and running we can start to make improvements.

#### Monitoring - Sidecar Pattern

The first thing we want to do is to monitor both the front end and API. To do
this without modifying the application code we are going to need a proxy to put
in front of the services.

We will use the sidecar pattern and use haproxy as v2 exposes stats in the
prometheus format so they can be scraped. Here are the steps we will
need to do and the iumprovements can be seen in [manifests/v1](./manifests/v1/)

- Create a configMap to hold the required haproxy configuration
- Add the proxy container to the deployment manifest
- Update the service to point to the container port of the proxy
- Create a serviceMonitor that scrapes the proxy's metrics endpoint

#### Logging - Adapter Pattern

Now that we have some metrics we also want to be able to get the event data
from the API that is being emitted to a log file within the contianer. As
Kubernetes reads logs from STDOUT this is not ideal and we will need to create
an adapter to adapt the log format to a structured logging format that can then
be inserted into something like elasticsearch.

#### Circuit Breaker - Ambassador Pattern

Lastly we are going to add an Ambassador contianer that will handle the
communication between the frontend and the API. The extra feature that this
proxy is going to give us is that if it notices an increased error rate with the
external service it will trip its circuit breaker and return a canned response
to the frontend and reduce the load being made to the API.
