
# 0.5.0

- In the initial release we generated randomized demand locally at the client. Demand is now passed to the client from the server. 
- We're now using a web socket to communicate metrics to the server and receive demand instructions from the server.
- Containers are no longer hard-coded - we get container images from the server
- The client now pulls Docker images - you don't have to do this as a manual step yourself

# Initial version

- Generates random demand for priority1 locally
