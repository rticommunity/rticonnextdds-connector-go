//go:build !static
// +build !static

package rti

/*
#cgo CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include
#cgo linux,amd64 CFLAGS: -DRTI_UNIX -DRTI_LINUX -DRTI_64BIT
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/linux-x64 -lrtiddsconnector -ldl -lm -lpthread -lrt
#cgo linux,arm64 CFLAGS: -DRTI_UNIX -DRTI_LINUX -DRTI_64BIT
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/linux-arm64 -lrtiddsconnector -ldl -lm -lpthread -lrt
#cgo linux,arm CFLAGS: -DRTI_UNIX -DRTI_LINUX
#cgo linux,arm LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/linux-arm -lrtiddsconnector -ldl -lm -lpthread -lrt
#cgo darwin,amd64 CFLAGS: -DRTI_UNIX -DRTI_DARWIN -DRTI_DARWIN10 -DRTI_64BIT -m64
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/osx-x64 -lrtiddsconnector -ldl -lm -lpthread
#cgo darwin,arm64 CFLAGS: -DRTI_UNIX -DRTI_DARWIN -DRTI_DARWIN10 -DRTI_64BIT
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/osx-arm64 -lrtiddsconnector -ldl -lm -lpthread
#cgo windows,amd64 CFLAGS: -DWIN32_LEAN_AND_MEAN -DWIN32 -D_WINDOWS -DRTI_WIN32 -DNDEBUG -DRTI_64BIT
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/rticonnextdds-connector/lib/win-x64 -lrtiddsconnector -lws2_32 -ladvapi32 -luser32 -lwinmm -lnetapi32 -lversion -lkernel32

#include "rticonnextdds-connector.h"
#include <stdlib.h>
*/
import "C"
