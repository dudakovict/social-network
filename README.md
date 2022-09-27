# Communication patterns in a microservice architecture
### Master's thesis
Author: Timon Dudaković

Mentor: doc. dr. sc. Nikola Tanković

### Abstract
One of the biggest challenges when migrating from an application based on a monolithic architecture to an application based on a microservice architecture is adopting changes present in communication paradigm. Converting from local method calls to unreliable cross-service synchronous and asynchronous calls adds a higher level of complexity and reduces efficiency in communication which violates performance in distributed systems. The challenges of designing and implementing a distributed system are well known, but the process is still long-lasting and complex. The solution presented in this thesis involves high levels of microservice isolation through the use of asynchronous communication patterns between the internal microservices. Although several microservices are communicating over synchronous communication protocols, they don't violate the integrity of the communication, and they maintain a certain level of isolation.

[Juraj Dobrila University of Pula](https://www.unipu.hr/)

[Faculty of informatics](https://fipu.unipu.hr/fipu)

[Master's thesis](https://dabar.srce.hr/islandora/object/unipu%3A7227)

## Requirements
- [Golang](https://go.dev/)
- [Docker](https://www.docker.com/)
- [Staticcheck](https://staticcheck.io/docs)
- [KIND](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [Kubernetes](https://kubernetes.io/docs/tasks/tools/)
- [Kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/)
