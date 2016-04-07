# Microscaling-in-a-box

This simple demo runs on your own machine and makes it easy to experiment with microscaling and see containers running locally 
on your machine. Full instructions for running the demo are at [app.microscaling.com](https://app.microscaling.com). 

[![Build Status](https://api.travis-ci.org/microscaling/microcaling.svg)](https://travis-ci.org/microscaling/microscaling) Go 1.4 1.5

## Building and running from source

The easiest way to run Microscaling-in-a-box is to [follow the instructions](http://app.microscaling.com). The `docker run` command 
pulls the latest image of this code from [Docker hub](https://hub.docker.com/u/microscaling/microscaling). 

[![](https://badge.imagelayers.io/microscaling/microscaling:latest.svg)](https://imagelayers.io/?images=microscaling/microscaling:latest 'Get your own badge on imagelayers.io')

If you want to build and run your own version locally:

- Clone this repo
- Build the code as a linux executable (since it runs inside a linux container): 
`GOOS=linux go build -o microscaling .`
- Build your own version of the container image, and give it a tag:
`docker build -t <your-tag-name> .`
- Specify `-it <your-tag-name>` instead of `-it microscaling/microscaling:latest` on `docker run` so that it picks up your version of the image

## Roadmap

We've got lots of ideas for improvements - here are a few headlines:

- We're planning to add support for a number of different schedulers. If there's one you'd particularly like to see us support, please let us know.
- We can improve performance by parallelizing requests to the Docker remote API.
- In this demo we simply randomize demand, but we'll add support for real-work demand models.

## Contact Us

We'd love to hear from you at [hello@microscaling.com](mailto:hello@microscaling.com) or on Twitter at [@microscaling](http://twitter.com/microscaling). 
And we welcome new [issues](https://github.com/microscaling/microscaling/issues) or pull requests!
