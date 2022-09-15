SHELL := /bin/bash

# ==============================================================================
# Testing running system
#
# For testing a simple query on the system. Don't forget to `make seed` first.
# curl --user "admin@example.com:gophers" http://localhost:3000/v1/users/token
# export TOKEN="COPY TOKEN STRING FROM LAST CALL"
# curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users/1/2
#
# For testing load on the service.
# go install github.com/rakyll/hey@latest
# hey -m GET -c 100 -n 10000 -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users/1/2
# hey -m GET -c 100 -n 10000 http://localhost:3000/v1/test
#
# Access metrics directly (4000) or through the sidecar (3001)
# go install github.com/divan/expvarmon@latest
# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
# expvarmon -ports=":3001" -endpoint="/metrics" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
#
# To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem
# ./users-admin genkey
#
# Testing Auth
# curl -il http://localhost:3000/v1/testauth
# curl -il -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/testauth
#
# Database Access
# dblab --host 0.0.0.0 --user postgres --db postgres --pass postgres --ssl disable --port 5432 --driver postgres
# pgcli postgres://postgres:postgres@localhost:5432/postgres
#
# Testing coverage.
# go test -coverprofile p.out
# go tool cover -html p.out
#
# Test debug endpoints.
# curl http://localhost:4000/debug/liveness
# curl http://localhost:4000/debug/readiness
#
# Running pgcli client for database.
# brew install pgcli
# pgcli postgresql://postgres:postgres@localhost
#
# Launch zipkin.
# http://localhost:9411/zipkin/

run:
	go run app/services/users-api/main.go | go run app/tooling/logfmt/main.go

admin:
	go run app/tooling/admin/main.go

# ==============================================================================
# Building containers

VERSION := 1.0

all: users-api posts-api comments-api email-api

users-api:
	docker build \
		-f zarf/docker/dockerfile.users-api \
		-t users-api-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

posts-api:
	docker build \
		-f zarf/docker/dockerfile.posts-api \
		-t posts-api-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

comments-api:
	docker build \
		-f zarf/docker/dockerfile.comments-api \
		-t comments-api-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

email-api:
	docker build \
		-f zarf/docker/dockerfile.email-api \
		-t email-api-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

email-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		business/data/email/email.proto
		
# ==============================================================================
# Running from within k8s/kind

KIND_CLUSTER := social-network-cluster

kind-up:
	kind create cluster \
		--image kindest/node:v1.21.1@sha256:fae9a58f17f18f06aeac9772ca8b5ac680ebbed985e266f711d936e91d113bad \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=services-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	cd zarf/k8s/kind/users/users-pod; kustomize edit set image users-api-image=users-api-amd64:$(VERSION)
	kind load docker-image users-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

	cd zarf/k8s/kind/posts/posts-pod; kustomize edit set image posts-api-image=posts-api-amd64:$(VERSION)
	kind load docker-image posts-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

	cd zarf/k8s/kind/comments/comments-pod; kustomize edit set image comments-api-image=comments-api-amd64:$(VERSION)
	kind load docker-image comments-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

	cd zarf/k8s/kind/email; kustomize edit set image email-api-image=email-api-amd64:$(VERSION)
	kind load docker-image email-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build zarf/k8s/kind/nats | kubectl apply -f -
	kubectl wait --namespace=services-system --timeout=240s --for=condition=Available deployment/nats-pod

	kustomize build zarf/k8s/kind/email | kubectl apply -f -
	
	kustomize build zarf/k8s/kind/users/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=240s --for=condition=Available deployment/users-database-pod
	kustomize build zarf/k8s/kind/users/zipkin-pod | kubectl apply -f -
	kubectl wait --namespace=zipkin-system --timeout=240s --for=condition=Available deployment/users-zipkin-pod
	kustomize build zarf/k8s/kind/users/users-pod | kubectl apply -f -

	kustomize build zarf/k8s/kind/posts/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=240s --for=condition=Available deployment/posts-database-pod
	kustomize build zarf/k8s/kind/posts/zipkin-pod | kubectl apply -f -
	kubectl wait --namespace=zipkin-system --timeout=240s --for=condition=Available deployment/posts-zipkin-pod
	kustomize build zarf/k8s/kind/posts/posts-pod | kubectl apply -f -

	kustomize build zarf/k8s/kind/comments/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=240s --for=condition=Available deployment/comments-database-pod
	kustomize build zarf/k8s/kind/comments/zipkin-pod | kubectl apply -f -
	kubectl wait --namespace=zipkin-system --timeout=240s --for=condition=Available deployment/comments-zipkin-pod
	kustomize build zarf/k8s/kind/comments/comments-pod | kubectl apply -f -

kind-services-delete:
	kustomize build zarf/k8s/kind/users/users-pod | kubectl delete -f -
	kustomize build zarf/k8s/kind/users/zipkin-pod | kubectl delete -f -
	kustomize build zarf/k8s/kind/users/database-pod | kubectl delete -f -

kind-zipkin-delete:
	kustomize build zarf/k8s/kind/users/zipkin-pod | kubectl delete -f -
	kustomize build zarf/k8s/kind/posts/zipkin-pod | kubectl delete -f -
	kustomize build zarf/k8s/kind/comments/zipkin-pod | kubectl delete -f -

kind-databases-delete:
	kustomize build zarf/k8s/kind/users/database-pod | kubectl delete -f -
	kustomize build zarf/k8s/kind/posts/database-pod | kubectl delete -f -
	kustomize build zarf/k8s/kind/comments/database-pod | kubectl delete -f -

kind-restart:
	kubectl rollout restart deployment users-pod
	kubectl rollout restart deployment posts-pod 
	kubectl rollout restart deployment comments-pod
	kubectl rollout restart deployment email-pod

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-logs:
	kubectl logs -l app=users --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go

kind-logs-users:
	kubectl logs -l app=users --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go -service=USERS-API

kind-logs-posts:
	kubectl logs -l app=posts --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go -service=POSTS-API

kind-logs-comments:
	kubectl logs -l app=comments --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go -service=COMMENTS-API

kind-logs-email:
	kubectl logs -l app=email --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go -service=EMAIL-API

kind-logs-db:
	kubectl logs -l app=database --namespace=database-system --all-containers=true -f --tail=100

kind-logs-zipkin:
	kubectl logs -l app=zipkin --namespace=zipkin-system --all-containers=true -f --tail=100

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-services:
	kubectl get pods -o wide --watch

kind-status-databases:
	kubectl get pods -o wide --watch --namespace=database-system

kind-status-zipkin:
	kubectl get pods -o wide --watch --namespace=zipkin-system

kind-describe:
	kubectl describe nodes
	kubectl describe svc
	kubectl describe pod -l app=users

kind-describe-deployment:
	kubectl describe deployment users-pod

kind-describe-replicaset:
	kubectl get rs
	kubectl describe rs -l app=users

kind-events:
	kubectl get ev --sort-by metadata.creationTimestamp

kind-events-warn:
	kubectl get ev --field-selector type=Warning --sort-by metadata.creationTimestamp

kind-context-services:
	kubectl config set-context --current --namespace=services-system

kind-context-databases:
	kubectl config set-context --current --namespace=database-system

kind-context-zipkin:
	kubectl config set-context --current --namespace=zipkin-system

kind-shell:
	kubectl exec -it $(shell kubectl get pods | grep sales | cut -c1-26) --container users-api -- /bin/sh

kind-database:
	# ./admin --db-disable-tls=1 migrate
	# ./admin --db-disable-tls=1 seed

# ==============================================================================
# Administration

migrate:
	go run app/tooling/sales-admin/main.go migrate

seed: migrate
	go run app/tooling/sales-admin/main.go seed

# ==============================================================================
# Running tests within the local computer

test:
	go test ./... -count=1
	staticcheck -checks=all ./...

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

# ==============================================================================
# Docker support

docker-down:
	docker rm -f $(shell docker ps -aq)

docker-clean:
	docker system prune -f	

docker-kind-logs:
	docker logs -f $(KIND_CLUSTER)-control-plane