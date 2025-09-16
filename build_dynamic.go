//go:build !static
// +build !static

package rti

// #cgo linux,amd64 CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_UNIX -DRTI_LINUX -DRTI_64BIT
// #cgo linux,amd64 LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/linux-x64 -lrtiddsconnector -ldl -lm -lpthread -lrt -Wl,-rpath,${SRCDIR}/rticonnextdds-connector/lib/linux-x64
// #cgo linux,arm64 CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_UNIX -DRTI_LINUX -DRTI_64BIT
// #cgo linux,arm64 LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/linux-arm64 -lrtiddsconnector -ldl -lm -lpthread -lrt -Wl,-rpath,${SRCDIR}/rticonnextdds-connector/lib/linux-arm64
// #cgo darwin,amd64 CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_UNIX -DRTI_DARWIN -DRTI_DARWIN10 -DRTI_64BIT -m64
// #cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/osx-x64 -lrtiddsconnector -ldl -lm -lpthread
// #cgo darwin,arm64 CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_UNIX -DRTI_DARWIN -DRTI_DARWIN10 -DRTI_64BIT
// #cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/osx-arm64 -lrtiddsconnector -ldl -lm -lpthread
import "C"
