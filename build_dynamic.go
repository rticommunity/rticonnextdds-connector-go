// +build !static

package rti

// #cgo darwin CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_UNIX -DRTI_DARWIN -DRTI_DARWIN10 -DRTI_64BIT -m64
// #cgo darwin LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/x64Darwin16clang8.0 -lrtiddsconnector -ldl -lm -lpthread
import "C"
