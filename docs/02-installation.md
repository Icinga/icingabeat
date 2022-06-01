# Installation

## Packages
Packages are available on [packages.icinga.com](https://packages.icinga.com).

Depending on your distribution and version you need to run one of the following
commands:

#### Debian

``` shell
wget -O - https://packages.icinga.com/icinga.key | apt-key add -
echo 'deb http://packages.icinga.com/debian icinga-stretch main' > etc/apt/sources.list.d/icinga.list
```

``` shell
apt-get update
apt-get install icingabeat
```

#### Ubuntu

``` shell
wget -O - https://packages.icinga.com/icinga.key | apt-key add -
echo 'deb http://packages.icinga.com/ubuntu icinga-xenial main' > etc/apt/sources.list.d/icinga.list
```

``` shell
apt-get update
apt-get install icingabeat
```

#### CentOS

``` shell
yum install epel-release
rpm --import https://packages.icinga.com/icinga.key
yum install https://packages.icinga.com/epel/icinga-rpm-release-7-latest.noarch.rpm
```

``` shell
yum install icingabeat
```

#### RHEL

``` shell
yum install https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
rpm --import https://packages.icinga.com/icinga.key
yum install https://packages.icinga.com/epel/icinga-rpm-release-7-latest.noarch.rpm
```

``` shell
yum install icingabeat
```

#### SLES

``` shell
rpm --import https://packages.icinga.com/icinga.key
zypper ar https://packages.icinga.com/SUSE/ICINGA-release.repo
zypper ref
```

``` shell
zypper install icingabeat
```

### Run
Make sure you have configured Icingabeat properly before starting it. Use one
of the following commands to start Icingabeat:

* `service icingabeat start` or
* `systemctl start icingabeat` or
* `/etc/init.d/icingabeat start`

## Dashboards
We have dashboards prepared that you can use when getting started with
Icingabeat. They are meant to give you some inspiration before you start
exploring the data by yourself.

Starting with icingabeat v7.17.4 you have to download and import the dashboards manually.

Download and upack `dashboards.zip` from the [latest release](https://github.com/Icinga/icingabeat/releases/latest) page.

Use Kibana's Import functionality to upload the `*.ndjson` files which will
import a bunch of saved objects, including dashboards and single visualizations.

## Manual Installation

Download and install a package or tarball from the
[latest release](https://github.com/Icinga/icingabeat/releases/latest) page.

## Development
Please follow [README.md](https://github.com/icinga/icingabeat/README.md) for
instructions about how to build icingabeat.
