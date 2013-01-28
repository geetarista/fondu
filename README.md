# Fondu ![Build Status](https://travis-ci.org/geetarista/fondu.png)

Fondu is a [CheeseShop](http://wiki.python.org/moin/CheeseShop) server powered by [Go](http://golang.org/).

## Features

* Cache packages from [PyPi](http://pypi.python.org/) in case of failure
* Install packages faster inside your network
* Store private packages
* easy install/pip compliant
* setup.py register|upload compatible

## Installation

[Install Go](http://golang.org/doc/install).

Install Fondu:

```shell
go install github.com/geetarista/fondu
```

## Usage

To start Fondu, all you have to do use the `fondu` command.

For registering/uploading, you need to tell distutils in ~/.pypirc:

```ini
[distutils]
index-servers =
  fondu

[fondu]
username = foo
password = bar
repository = http://your-host:3638/
```

Note that username and password must be passed to pip, but fondu does not support authentication at this time.

## Configure

You can pass Fondu the `-f` flag which tells it what config file to use.

```shell
fondu -f /etc/fondu.conf
```

An example configuration file looks like this (defaults shown):

```ini
[fondu]
data_directory = /data/fondu
port = 3638
pypi_mirror = http://pypi.python.org
```

## Thanks

Based on original Python version by [Mitchell Hashimoto](http://mitchellh.com).

[Kiip](http://kiip.me), for the opportunity to build this.

## License

MIT. See [LICENSE](https://github.com/geetarista/fondu/blob/master/LICENSE).
