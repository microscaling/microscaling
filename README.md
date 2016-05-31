# Microscaling Engine

Our Microscaling Engine provides automation, resilience and efficiency for microservice architectures. You can use our [Microscaling-in-a-Box](https://app.microscaling.com) site to
experiment with microscaling. Or visit [microscaling.com](https://microscaling.com) to find out more about our product and Microscaling Systems.

[![Build Status](https://travis-ci.org/microscaling/microscaling.svg?branch=master)](https://travis-ci.org/microscaling/microscaling) Go 1.4 1.5 1.6
[![](https://badge.imagelayers.io/microscaling/microscaling:latest.svg)](https://imagelayers.io/?images=microscaling/microscaling:latest 'Get your own badge on imagelayers.io')

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

## Building and running from source

The easiest way to run Microscaling-in-a-box is to [follow the instructions](http://app.microscaling.com). The `docker run` command 
pulls the latest image of this code from [Docker hub](https://hub.docker.com/u/microscaling/microscaling). 

If you want to build and run your own version locally:

- Clone this repo
- Build the code as a linux executable (since it runs inside a linux container): 
`GOOS=linux go build -o microscaling .`
- Build your own version of the container image, and give it a tag:
`docker build -t <your-tag-name> .`
- Specify `-it <your-tag-name>` instead of `-it microscaling/microscaling:latest` on `docker run` so that it picks up your version of the image

## Licensing

Microscaling Engine is licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/microscaling/microscaling/blob/master/LICENSE) for the full license text.

## Contributing

We'd love to get contributions from you! Please see [CONTRIBUTING.md](https://github.com/microscaling/microscaling/blob/master/CONTRIBUTING.md) for more details.

## Contact Us

We'd love to hear from you at [hello@microscaling.com](mailto:hello@microscaling.com) or on Twitter at [@microscaling](http://twitter.com/microscaling). 
And we welcome new [issues](https://github.com/microscaling/microscaling/issues) or pull requests!
