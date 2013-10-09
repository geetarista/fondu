# Fondu ![Build Status](https://travis-ci.org/geetarista/fondu.png)

Fondu is a [CheeseShop](http://wiki.python.org/moin/CheeseShop) server powered by [Go](http://golang.org/).

## Features

* Cache packages from [PyPi](http://pypi.python.org/) in case of failure
* Install packages faster inside your network
* Store private packages
* easy_install/pip compliant
* setup.py register/upload compatible

## Installation

Just [download the binary for your platform](https://github.com/geetarista/fondu/releases).

Please not that so far the only platforms I've tested on are darwin-amd64 and linux-amd64. If you have any platform-specific problems, please file an issue so I can investigate.

## Usage

To start Fondu, all you have to do use the `fondu` command.

## Configure

You can pass Fondu the `-f` flag which tells it what config file to use.

```shell
fondu -f /etc/fondu.conf
```

An example configuration file looks like this (defaults shown):

```ini
[fondu]
data_directory = data
port = 3638
pypi_mirror = http://pypi.python.org
```

## Uploading

For registering/uploading packages, you need to tell distutils where to look in `pypirc`:

```ini
[distutils]
index-servers =
  fondu

[fondu]
username = foo
password = bar
repository = http://your-host:3638/
```

Or you can just set it when uploading inline:

```bash
python setup.py sdist upload -r fondu
python setup.py sdist upload -r http://your-host:3638/
```

Or as an environment variable:

```bash
export PIP_INDEX_URL=http://your-host:3638/simple/
```

Note that username and password can be passed to pip, but fondu does not support authentication at this time.

## Demo

See the [demo](https://github.com/geetarista/fondu/tree/master/demo) directory for information on how to use [Vagrant](http://www.vagrantup.com/) to run Fondu locally.

## License

MIT. See [LICENSE](https://github.com/geetarista/fondu/blob/master/LICENSE).
