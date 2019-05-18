rticonnextdds-connector-go
=======
[![GoDoc](https://godoc.org/github.com/rticommunity/rticonnextdds-connector-go?status.png)](https://godoc.org/github.com/rticommunity/rticonnextdds-connector-go) [![Build Status](https://travis-ci.org/rticommunity/rticonnextdds-connector-go.svg?branch=master)](https://travis-ci.org/rticommunity/rticonnextdds-connector-go) [![Coverage](https://codecov.io/gh/rticommunity/rticonnextdds-connector-go/branch/master/graph/badge.svg)](https://codecov.io/gh/rticommunity/rticonnextdds-connector-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/rticommunity/rticonnextdds-connector-go)](https://goreportcard.com/report/github.com/rticommunity/rticonnextdds-connector-go)

### RTI Connector for Connext DDS
*RTI Connector* for Connext DDS is a quick and easy way to access the power and
functionality of [RTI Connext DDS](http://www.rti.com/products/index.html).
It is based on [XML-Based Application Creation](https://community.rti.com/static/documentation/connext-dds/5.3.1/doc/manuals/connext_dds/xml_application_creation/RTI_ConnextDDS_CoreLibraries_XML_AppCreation_GettingStarted.pdf) and Dynamic Data.

*Connector* was created by the RTI Research Group to quickly and easily develop demos
and proofs of concept. It can be useful for anybody that needs
a quick way to develop an application communicating over the Connext Databus.
Thanks to the binding with multiple programming languages, you can integrate
with many other available technologies.

The *Connector* library is provided in binary form for [select architectures](https://github.com/rticommunity/rticonnextdds-connector/tree/master/lib). Language bindings and examples are provided in source format.

Go *Connector* leverages [cgo](https://golang.org/cmd/cgo) to call its C library;
this detail is hidden in a Go wrapper. RTI tested with Go v.12, v1.11, v1.10, and v1.9.

### Getting started
Be sure you have Go installed and set your go workspace ($GOPATH). Then run:

``` bash
$ go get github.com/rticommunity/rticonnextdds-connector-go
```

Check that the Go *Connector* repository is cloned properly at the following location:
$GOPATH/src/github.com/rticommunity/rticonnextdds-connector-go.

Then take a look at [this document](examples/README.md) to build and run examples.

### Platform support
Go *Connector* builds its library for few [select architectures](https://github.com/rticommunity/rticonnextdds-connector/tree/master/lib). If you need another architecture, please contact your RTI account manager or sales@rti.com.

If you want to check the version of the libraries you can run the following command:

``` bash
strings librtiddsconnector.so | grep BUILD
```

### Threading model
The *Connector* Native API does not yet implement any mechanism for thread safety. Originally, the *Connector* native code was built to work with *RTI Prototyper* and Lua. That was a single-threaded loop. RTI then introduced support for JavaScript, Python, and Go. For now, you are responsible for protecting calls to *Connector*. Thread safety
may be implemented in the future.

### Support
*Connector* is an experimental RTI product. If you have questions, use the [RTI Community Forum](https://community.rti.com/forums/technical-questions).

### Documentation
The best way to get started with *Connector* is to look at the
examples; you will see that it is very easy to use.

See [an overview of the API](https://godoc.org/github.com/rticommunity/rticonnextdds-connector-go).

### License
With the sole exception of the contents of the "examples" subdirectory, all use of this product is subject to the RTI Software License Agreement included at the top level of this repository. Files within the "examples" subdirectory are licensed as marked within the file.

This software is an experimental (aka "pre-production") product. The Software is provided "as is", with no warranty of any type, including any warranty for fitness for any purpose. RTI is under no obligation to maintain or support the Software. RTI shall not be liable for any incidental or consequential damages arising out of the use or inability to use the software.
