# Komunikacijski obrasci u mikroservisnoj arhitekturi
### Diplomski rad
Autor: Timon Dudaković

Mentor: doc. dr. sc. Nikola Tanković

### Sažetak
Jedan od najvećih izazova pri prijelazu s aplikacije temeljene na monolitnoj arhitekturi na aplikaciju temeljenu na mikroservisnoj arhitekturi je usvajanje promjena prisutnih u komunikacijskoj paradigmi. Prijelaz iz poziva lokalnih metoda u nepouzdane sinkrone i asinkrone pozive između servisa dodaje višu razinu složenosti i smanjuje učinkovitost komunikacije što narušava performanse u raspodijeljenim sustavima. Izazovi dizajniranja i implementacije raspodijeljenog sustava su dobro poznati, ali je proces još uvijek dugotrajan i složen. Rješenje predstavljeno u ovom diplomskom radu uključuje visoke razine izolacije mikroservisa korištenjem asinkronih komunikacijskih obrazaca između internih mikroservisa. Iako postoji niz mikroservisa koji komuniciraju preko sinkronih komunikacijskih protokola, oni ne narušavaju integritet komunikacije i održavaju određenu razinu izolacije.

[Sveučilište Jurja Dobrile u Puli](https://www.unipu.hr/)

[Fakultet informatike u Puli](https://fipu.unipu.hr/fipu)

# Communication patterns in a microservice architecture
### Master's thesis
Author: Timon Dudaković

Mentor: doc. dr. sc. Nikola Tanković

### Abstract
One of the biggest challenges when migrating from an application based on a monolothic arhictecture to an application based on a microservice architecture is adopting changes present in communication paradigm. Converting from local method calls to unreliable cross-service synchronous and asynchronous calls adds a higher level of complexity and reduces efficiency in communication which violates peformance in distributed systems. The challenges of designing and implementing a distributed system are well known, but the process is still long lasting and complex. The solution presented in this thesis involves high levels of microservice isolation through the use of asynchronous communication patterns between the internal microservices. Although there is a number of microservices communicating over synchronous communication protocols, they don't violate the integrity of the communication, and they maintain a certain level of isolation.

[Juraj Dobrila University of Pula](https://www.unipu.hr/)

[Faculty of informatics](https://fipu.unipu.hr/fipu)

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

Apply the configuration files to the Kind environment: ```make kind-apply```
