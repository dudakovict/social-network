# Communication patterns in a microservice architecture
### Master's thesis
Author: Timon Dudaković
Mentor: doc. dr. sc. Nikola Tanković

## Requirements
- [Golang](https://go.dev/)
- [Docker](https://www.docker.com/)
- [Staticcheck](https://staticcheck.io/docs)
- [KIND](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [Kubernetes](https://kubernetes.io/docs/tasks/tools/)
- [Kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/)


## Installation
### KIND
MacOS: ```brew install kind```

Windows: ```choco install kind```

### Kustomize
MacOS: ```brew install kustomize```

Windows: ```choco install kustomize```

## Configuration
Run ```kubectl create secret generic smtp-credentials --from-literal=email --from_literal=email_app_password``` replacing **email** and **email_app_password** with valid data.

To configure your app password go to your gmail account and head over to security.
There you can create an Email app password that is located below your 2FA checkbox.

## Running the project
Build the Docker images: ```make all```

Spin up a local Kind environment: ```make kind-up```

Load the Docker images into the Kind environment: ```make kind-load```

Apply the configuration files to the Kind environment: ```make kind-apply`
