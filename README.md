# beanstalk-cli

A full functional beanstalkd CLI for [beanstalkd](https://github.com/beanstalkd/beanstalkd).

## Features

 * Full featured CLI to implement each command of [beanstalkd protocol](https://github.com/beanstalkd/beanstalkd/blob/master/doc/protocol.txt).
 * Connect to beanstalkd server in tcp mode or unix socket mode.
 * Output can be formatted to [JSON](https://json.org/).

## Example

 * Put a job
```
$ beanstalk put -a tcp://127.0.0.1:11300 foobar
id:	2
```
 * Reserve a job
```
$ beanstalk reserve -a tcp://127.0.0.1:11300
id:  	1
body:	2222
```
 * Bury a job
```
$ beanstalk-cli bury -a tcp://127.0.0.1:11300 2
id:  	2
body:	foobar
```
 * Delete a job
```
$ beanstalk-cli delete -a tcp://127.0.0.1:11300 2
```
 * Connect to a server listening on unix socket
```
$ beanstalkd -l unix:///tmp/beanstalkd.sock &
$ beanstalk-cli stats -a unix:///tmp/beanstalkd.sock
```

## Development

This CLI tools is wrtten in [golang](https://golang.org), based on the official [go-beanstalk](https://github.com/beanstalkd/go-beanstalk) lib.
