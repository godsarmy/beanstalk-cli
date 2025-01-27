<p align="center">
    <a href="https://github.com/godsarmy/beanstalk-cli/releases"><img src="https://img.shields.io/github/downloads/godsarmy/beanstalk-cli/total.svg" alt="Downloads"></a>
    <a href="https://github.com/godsarmy/beanstalk-cli/blob/master/LICENSE"><img src="https://img.shields.io/github/license/mashape/apistatus.svg" alt="Licenses"></a>
    <a href="https://github.com/godsarmy/beanstalk-cli/releases"><img src="https://img.shields.io/github/release/godsarmy/beanstalk-cli.svg?label=Release" alt="Release"></a>
</p>

## Overview

beanstalk-cli: A Powerful Command-Line Interface for [beanstalkd](https://github.com/beanstalkd/beanstalkd) work queue.

## Features

`beanstalk-cli` provides a comprehensive set of commands for managing your Beanstalkd queues directly from your terminal.  This makes it ideal for scripting, automation, debugging, and integrating with your CI/CD pipelines.

 * **Cross-platform support:** macOS/Linux/Windows 32/64-bit
 * **Full Beanstalkd Protocol Support:** Implement every command of the [beanstalkd protocol](https://github.com/beanstalkd/beanstalkd/blob/master/doc/protocol.txt), giving you complete control over your queues.
 * **TCP and Unix Socket Connections:** Connect to Beanstalkd servers using either TCP or Unix sockets, providing flexibility for different deployment environments.
 * **JSON Output:** Format output as [JSON](https://json.org/) for easy parsing and integration with other tools and scripts, enabling powerful automation workflows.
 * **Easy Job Management:** Put, reserve, bury, delete, and inspect jobs with simple and intuitive commands.
 * **Queue Statistics:** Monitor queue performance and health by retrieving detailed statistics.
 * **Tube Management:** Use tubes to organize your jobs and prioritize processing.  This is crucial for complex applications with different job types.

## Example

 * Put a job
```sh
$ beanstalk put -a tcp://127.0.0.1:11300 foobar
id:	2
```
 * Reserve a job
```sh
$ beanstalk reserve -a tcp://127.0.0.1:11300
id:  	1
body:	2222
```
 * Bury a job
```sh
$ beanstalk-cli bury -a tcp://127.0.0.1:11300 2
id:  	2
body:	foobar
```
 * Delete a job
```sh
$ beanstalk-cli delete -a tcp://127.0.0.1:11300 2
```
 * Show server stats
 ```sh
$ beanstalk-cli stats tcp://127.0.0.1:11300
 ```
 * Connect to a server with address defined in environment variable
 ```sh
$ export BS_ADDRESS=tcp://127.0.0.1:11300
$ beanstalk-cli stats-tube default
 ```
 * Connect to a server listening on unix socket
```sh
$ beanstalkd -l unix:///tmp/beanstalkd.sock &
$ beanstalk-cli stats -a unix:///tmp/beanstalkd.sock
```

## Installation

[Precompiled binaries](https://github.com/godsarmy/beanstalk-cli/releases) for supported operating systems are available.

## Development

This CLI tools is wrtten in [golang](https://golang.org), based on the official [go-beanstalk](https://github.com/beanstalkd/go-beanstalk) lib.
