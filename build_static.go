// +build static

package rti

// #cgo linux,amd64 CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_UNIX -DRTI_LINUX -DRTI_64BIT
// #cgo linux,amd64 LDFLAGS: -L${SRCDIR}/static_lib/x64Linux2.6gcc4.4.5 -lrtiddsconnectorz -lluaz -lnddscz -lnddscorez -ldl -lnsl -lm -lpthread -lrt
// #cgo darwin CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/rticonnextdds-connector/include -DRTI_UNIX -DRTI_DARWIN -DRTI_DARWIN10 -DRTI_64BIT -DRTI_STATIC -m64
// #cgo darwin LDFLAGS: -L${SRCDIR}/static_lib/x64Darwin16clang8.0 -lrtiddsconnectorz -lluaz -lnddscz -lnddscorez -ldl -lm -lpthread
import "C"
