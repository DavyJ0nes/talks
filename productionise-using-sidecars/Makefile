.DEFAULT_GOAL := help

# -- Variables --
CURRENT_CONTEXT = $(shell kubectl config current-context)

.PHONY: all
all: help

## install_helm: Installs the Helm Tiller on Cluster
.PHONY: install_helm
install_helm:
ifneq ("$(CURRENT_CONTEXT)","docker-desktop")
	$(call red, "## Not on local cluster.")
	@exit 1
endif
	$(call blue, "# Installing Helm to Cluster...")
	@helm init

## install_prometheus: Installs the Prometheus Operator
.PHONY: install_prometheus
install_prometheus:
ifneq ("$(CURRENT_CONTEXT)","docker-desktop")
	$(call red, "## Not on local cluster.")
	@exit 1
endif
	$(call blue, "# Installing Prometheus Operator to Cluster...")
	@helm install --name monitoring stable/prometheus-operator

## apply_prometheus_rules: Applies recording and alerting rules
.PHONY: apply_prometheus_rules
apply_prometheus_rules:
ifneq ("$(CURRENT_CONTEXT)","docker-desktop")
	$(call red, "## Not on local cluster.")
	@exit 1
endif
	$(call blue, "# Applying recording and alerting rules...")
	@kubectl apply -f core/prometheus-rules/rules.yaml
	@kubectl apply -f core/prometheus-rules/alert-rules.yaml

## install_nginx: Installs the Nginx Ingress
.PHONY: install_nginx
install_nginx:
ifneq ("$(CURRENT_CONTEXT)","docker-desktop")
	$(call red, "## Not on local cluster.")
	@exit 1
endif
	$(call blue, "# Installing Nginx Ingress to Cluster...")
	@kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/mandatory.yaml
	@kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/cloud-generic.yaml

## install_cert_manager: Installs Cert Manager
.PHONY: install_cert_manager
install_cert_manager:
	@kubectl apply -f https://raw.githubusercontent.com/jetstack/cert-manager/release-0.9/deploy/manifests/00-crds.yaml
	@kubectl create namespace cert-manager
	@kubectl label namespace cert-manager certmanager.k8s.io/disable-validation=true
	@helm repo add jetstack https://charts.jetstack.io
	@helm repo update
	@helm install \
  	--name cert-manager \
  	--namespace cert-manager \
  	--version v0.9.0 \
  	jetstack/cert-manager
	@kubectl create -f core/cert-manager/issuer.yaml

## forward_prometheus: Port Forwards Prometheus Server
.PHONY: forward_prometheus
forward_prometheus:
	$(call blue, "# Port Forwarding Prometheus Server...")
	@kubectl port-forward prometheus-monitoring-prometheus-oper-prometheus-0 9090:9090

## open_grafana: Open Grafana in Browser
.PHONY: open_grafana
open_grafana:
	$(call blue, "# Opening Grafana...")
	@open http://localhost:$(kubectl describe svc monitoring-grafana | grep "^NodePort:" | awk '{print $NF}' | tr -d "/TCP")

## images: creates dependany docker images
.PHONY: images
images:
	$(call blue, "# Creating and/or pulling images ...")
	@cd services/feelgood-api && make image
	@cd services/frontend && docker build -t davyj0nes/gophercon-frontend:v0.4.0 .
	@docker pull haproxy:2.0.4
	@docker pull alpine:3.10
	@docker pull busybox:1.28
	@docker pull davyj0nes:breaker:v0.1.0

## initial: installs the inital version of the services
.PHONY: initial
initial:
	$(call blue, "# Creating inital version ...")
	@kubectl apply -f manifests/initial/api/
	@kubectl apply -f manifests/initial/frontend/

## v1: installs the v1 version of the services
.PHONY: v1
v1:
	$(call blue, "# Creating v1 version ...")
	@kubectl apply -f manifests/v1/api/
	@kubectl apply -f manifests/v1/frontend/

## v2: installs the v2 version of the services
.PHONY: v2
v2:
	$(call blue, "# Creating v2 version ...")
	@kubectl apply -f manifests/v2/api/
	@kubectl apply -f manifests/v2/frontend/

## v3: installs the v3 version of the services
.PHONY: v3
v3:
	$(call blue, "# Creating v3 version ...")
	@kubectl apply -f manifests/v3/api/
	@kubectl apply -f manifests/v3/frontend/

## reset_demo: deletes all demo resources
.PHONY: reset_demo
reset_demo:
	$(call blue, "# deleting demo resources ...")
	@kubectl delete all -l release=gophercon

## create_load: create some load on the service
.PHONY: create_load
create_load:
	$(call blue, "# creating load on service ...")
	hey -z 30m -c 10 http://localhost

## help: prints this help message
.PHONY: help
help: Makefile
	@echo
	$(call blue, " Choose a command run:")
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

# -- Helper Functions --
define blue
	@tput setaf 4
	@echo $1
	@tput sgr0
endef

define red
	@tput setaf 1
	@echo $1
	@tput sgr0
endef

define yellow
	@tput setaf 3
	@echo $1
	@tput sgr0
endef
