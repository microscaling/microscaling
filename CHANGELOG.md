# 0.7.2

- Docker Labels changed from com.microscaling to org.label-schema

# 0.7.1

- Add monitors package allowing metrics to be sent to multiple locations.
- Add a Makefile with targets for building and releasing the Docker image.
- Add metadata labels with dynamic values populated by the Makefile.

# 0.7.0

- Add support for Marathon/Mesos as a scheduler
- Add support for scaling to maintain the length of an Azure queue

# 0.6.0

- Supports scaling to maintain the length of an NSQ queue

# 0.5.3

- Rename from force12io/force12 to microscaling/microscaling

# 0.5.2

- Remove container volumes when we remove containers
- Set the PublishAllPorts flag on (eventually this will be configurable)

# 0.5.1 

Get those pesky UTs running

# 0.5.0

- In the initial release we generated randomized demand locally at the client. Demand is now passed to the client from the server. 
- We're now using a web socket to communicate metrics to the server and receive demand instructions from the server.
- Containers are no longer hard-coded - we get container images from the server
- The client now pulls Docker images - you don't have to do this as a manual step yourself

# Initial version

- Generates random demand for priority1 locally
