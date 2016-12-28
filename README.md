# Microscaling Engine

Our Microscaling Engine provides automation, resilience and efficiency for microservice architectures. You can use our [Microscaling-in-a-Box](https://app.microscaling.com) site to
experiment with microscaling. Or visit [microscaling.com](https://microscaling.com) to find out more about our product and Microscaling Systems.

### Docker Image
[![](https://images.microbadger.com/badges/image/microscaling/microscaling.svg)](http://microbadger.com/images/microscaling/microscaling "Get your own image badge on microbadger.com") [![](https://images.microbadger.com/badges/version/microscaling/microscaling.svg)](http://microbadger.com/images/microscaling/microscaling "Get your own version badge on microbadger.com") [![](https://images.microbadger.com/badges/commit/microscaling/microscaling.svg)](http://microbadger.com/images/microscaling/microscaling "Get your own commit badge on microbadger.com") [![](https://images.microbadger.com/badges/license/microscaling/microscaling.svg)](http://microbadger.com/images/microscaling/microscaling "Get your own license badge on microbadger.com")

### Build
[![Build Status](https://travis-ci.org/microscaling/microscaling.svg?branch=master)](https://travis-ci.org/microscaling/microscaling) 

Go 1.4 1.5 1.6 

Microscaling Engine is under development, so we're not making any promises about forward compatibility, and we wouldn't advise running it on production machines yet. But if you're keen to get it into production we'd love to hear from you.

## Schedulers

Microscaling Engine will integrate with all the popular container schedulers. Currently we support

* Docker API
* Marathon 

Support for more schedulers is coming soon. Let us know if there is a particular scheduler you wish us to support.

## Metrics

Currently we support scaling a queue to maintain a target length. Support for more metrics is coming soon.

### Queue Types

* [NSQ](http://nsq.io) - see this [blog post](http://blog.microscaling.com/2016/04/microscaling-with-nsq-queue.html) for more details.
* Azure storage queues - this [blog post](http://blog.microscaling.com/2016/05/microscaling-marathon-with-dcos-on.html) describes using the Azure queue as the metric while running microscaled tasks on DC/OS.

Support for more message queues is coming soon. Let us know if there is a particular queue you wish us to integrate with.

## Running

The easiest way to run Microscaling-in-a-box is to [follow the instructions](http://app.microscaling.com). The `docker run` command
pulls the latest image of this code from [Docker hub](https://hub.docker.com/u/microscaling/microscaling).

## Running with label-based config

Get scaling parameters from your image metadata by configuring them with the following labels:
- com.microscaling.is-scalable
- com.microscaling.priority
- com.microscaling.max-delta
- com.microscaling.min-containers
- com.microscaling.max-containers

Download the compose file and add the following environment variable to the environment settings for the microscaling image:
```MSS_CONFIG=LABEL```

## Building from source

If you want to build and run your own version locally:

- Clone this repo
- Build your own version of the Docker image `DOCKER_IMAGE=<your-image> make build`
- Specify `-it <your-image>` instead of `-it microscaling/microscaling:latest` on `docker run` so that it picks up your version of the image

## Licensing

Microscaling Engine is licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/microscaling/microscaling/blob/master/LICENSE) for the full license text.

## Contributing

We'd love to get contributions from you! Please see [CONTRIBUTING.md](https://github.com/microscaling/microscaling/blob/master/CONTRIBUTING.md) for more details.

## Contact Us

We'd love to hear from you at [hello@microscaling.com](mailto:hello@microscaling.com) or on Twitter at [@microscaling](http://twitter.com/microscaling). 
And we welcome new [issues](https://github.com/microscaling/microscaling/issues) or pull requests!
