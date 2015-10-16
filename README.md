# Force12 - Microscaling-in-a-box

This simple demo runs on your own machine and makes it easy to experiment with microscaling and see containers running locally 
on your machine. Full instructions for running the demo are at [app.force12.io](http://app.force12.io). 

[![Build Status](https://api.travis-ci.org/force12io/force12.svg)](https://travis-ci.org/force12io/force12) Go 1.4 1.5

## Building and running from source

The easiest way to run Microscaling-in-a-box is to [follow the instructions](http://app.force12.io). The `docker run` command 
pulls the latest image of this code from [Docker hub](https://hub.docker.com/u/force12io/force12). 

[![](https://badge.imagelayers.io/force12io/force12:latest.svg)](https://imagelayers.io/?images=force12io/force12:latest 'Get your own badge on imagelayers.io')

If you want to build and run your own version locally:

- Clone this repo
- Build the code as a linux executable (since it runs inside a linux container): 
`GOOS=linux go build -o force12 .`
- Build your own version of the container image, and give it a tag:
`docker build -t <your-tag-name> .`
- Specify `-it <your-tag-name>` instead of `-it force12io/force12:latest` on `docker run` so that it picks up your version of the image

## Roadmap

We've got lots of ideas for improvements - here are a few headlines:

- We're planning to add support for a number of different schedulers. If there's one you'd particularly like to see us support, please let us know.
- We can improve performance by parallelizing requests to the Docker remote API.
- In this demo we simply randomize demand, but we'll add support for real-work demand models.
- The state API is currently hard-coded to our demo priority1 & priority2 tasks. We'll make it much more generic to support any task more complex scenarios.

## Contact Us

We'd love to hear from you at [hello@force12.io](mailto:hello@force12.io) or on Twitter at [@force12io](http://twitter.com/force12io). 
And we welcome new [issues](https://github.com/force12io/force12/issues) or pull requests!
