# Fondu Demo

This demo shows you how to set up a test Fondu server using Vagrant to see how it all comes together.

## Installation

First, make sure you have [Vagrant](http://www.vagrantup.com) installed.

Then, you need have a base box to install from. I used `precise64`:

```bash
vagrant box add precise64 http://files.vagrantup.com/precise64.box
```

Note that if you use a different base box, you'll need to change it in the `Vagrantfile`.

Then inside this directory, just run `vagrant up` and your VM will be running.

## Usage

If you run `vagrant ssh`, you can start up Fondu inside Vagrant:

```bash
fondu > fondu.log 2>&1
```

Here I just start the server and redirect any output.

Then from your terminal (outside Vagrant), you can install a package that will hit Fondu:

```bash
sudo PIP_INDEX_URL=http://192.168.50.4:3638/simple/ pip install amqp
```

You can see that Fondu downloaded the package and served it from `/home/vagrant/data/amqp`. If you try uninstalling and re-installing the package again, it will just server the cached file.
