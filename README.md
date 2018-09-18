rticonnextdds-connector-go
=======

[![Build Status](https://travis-ci.org/kyoungho/rticonnextdds-connector-go.svg?branch=master)]https://travis-ci.org/kyoungho/rticonnextdds-connector-go)

### RTI Connector for Connext DDS
RTI Connector for Connext DDS is a quick and easy way to access the power and
functionality of [RTI Connext DDS](http://www.rti.com/products/index.html).
It is based on [XML App Creation](https://community.rti.com/static/documentation/connext-dds/5.3.1/doc/manuals/connext_dds/xml_application_creation/RTI_ConnextDDS_CoreLibraries_XML_AppCreation_GettingStarted.pdf) and Dynamic Data.

RTI Connector was created by the RTI Research Group to quickly and easily develop demos
and proof of concept. We think that it can be useful for anybody that needs
a quick way to script tests and interact with DDS using different scripting languages.

It can be used to quickly create tests for your distributed system and, thanks
to the binding with scripting languages and the use of XML, to easily integrate
with tons of other available technologies.

The RTI Connector library is provided in binary form for selected architectures. Scripting language bindings and examples are provided in source format.

For **Go Connector**, we leveraged [cgo](https://golang.org/cmd/cgo) to call our C library, but we try to hide
that from you using a nice Go wrapper. We tested with Go v1.10.1.

### Getting started
Be sure you have Go installed and set your go workspace ($GOPATH). Then run:

``` bash
$ go get github.com/rticommunity/rticonnextdds-connector-go
```

Check out the Go Connector repository is cloned properly at the following location.
$GOPATH/src/github.com/rticommunity/rticonnextdds-connector-go

Then, take a look at [this document](examples/README.md) to build and run examples!

### Platform support
We are building our library for few architectures only. Check them out [here](https://github.com/rticommunity/rticonnextdds-connector/tree/master/lib). If you need another architecture.

If you want to check the version of the libraries you can run the following command:

``` bash
strings librtiddsconnector.so | grep BUILD
```

### Threading model
The RTI Connext DDS Connector Native API do not yet implement any mechanism for thread safety. Originally the Connector native code was built to work with RTI DDS Prototyper and Lua. That was a single threaded loop. We then introduced support for javascript, python, and Go. For now the responsibility of protecting calls to the Connector are left to the user. This may change in the future.

### Support
This is an experimental RTI product. As such we do offer support through the [RTI Community Forum](https://community.rti.com/forums/technical-questions) where fellow users and RTI engineers can help you.
We'd love your feedback.

### Documentation
We do not have much documentation yet. But we promise you: if you look at the
examples you'll see that is very easy to use our connector.

For an overview of the API, check this [page](examples/README.md).

### License
With the sole exception of the contents of the "examples" subdirectory, all use of this product is subject to the RTI Software License Agreement included at the top level of this repository. Files within the "examples" subdirectory are licensed as marked within the file.

This software is an experimental (aka "pre-production") product. The Software is provided "as is", with no warranty of any type, including any warranty for fitness for any purpose. RTI is under no obligation to maintain or support the Software. RTI shall not be liable for any incidental or consequential damages arising out of the use or inability to use the software.
