package zfs

// #cgo CFLAGS: -I/usr/include/libzfs
// #cgo CFLAGS: -I/usr/include/libspl
// #cgo CFLAGS: -I/usr/include/x86_64-linux-gnu
// #cgo CFLAGS: -DHAVE_IOCTL_IN_SYS_IOCTL_H=1
// #cgo LDFLAGS: -lzfs
// #include <libzfs.h>
import "C"

import (
	"errors"
	"fmt"
	"strconv"
)

var libzfs *C.libzfs_handle_t

// Zfs is a pointer to libzfs handler
type Zfs struct {
	zfs *C.zfs_handle_t
}

func init() {
	libzfs = C.libzfs_init()
	if libzfs == nil {
		panic("Can't initialize libzfs")
	}
}

// OpenDataset opens dataset for further operations; Returns Zfs handle
func OpenDataset(dataset string) (handle *Zfs, err error) {
	handle = nil
	h := new(Zfs)

	zfs := C.zfs_open(libzfs, C.CString(dataset), C.ZFS_TYPE_FILESYSTEM)
	if zfs == nil {
		err = fmt.Errorf("Can't dataset %v", dataset)
	} else {
		h.zfs = zfs
		handle = h
	}
	return
}

// GetUserQuota returns user quota is bytes
func (z *Zfs) GetUserQuota(uid uint64) (quota uint64, err error) {
	prop := C.CString(fmt.Sprintf("userquota@%v", uid))
	var buf *C.char
	buf = (*C.char)(C.calloc(C.ZFS_MAXPROPLEN, 1))

	ret := C.zfs_prop_get_userquota(z.zfs, prop, buf, C.ZFS_MAXPROPLEN, C.B_TRUE)
	if ret != 0 {
		err = fmt.Errorf("Can't get user quota for uid %v", uid)
	} else {
		quota, _ = strconv.ParseUint(C.GoString(buf), 0, 64)
	}

	return
}

// GetUserUsed returns user usage
func (z *Zfs) GetUserUsed(uid uint64) (quota uint64, err error) {
	prop := C.CString(fmt.Sprintf("userused@%v", uid))
	var buf *C.char
	buf = (*C.char)(C.calloc(C.ZFS_MAXPROPLEN, 1))

	ret := C.zfs_prop_get_userquota(z.zfs, prop, buf, C.ZFS_MAXPROPLEN, C.B_TRUE)
	if ret != 0 {
		err = fmt.Errorf("Can't get user quota for uid %v", uid)
	} else {
		quota, _ = strconv.ParseUint(C.GoString(buf), 0, 64)
	}

	return
}

// GetGroupQuota returns group quota
func (z *Zfs) GetGroupQuota(gid uint64) (quota uint64, err error) {
	prop := C.CString(fmt.Sprintf("groupquota@%v", gid))
	var buf *C.char
	buf = (*C.char)(C.calloc(C.ZFS_MAXPROPLEN, 1))

	ret := C.zfs_prop_get_userquota(z.zfs, prop, buf, C.ZFS_MAXPROPLEN, C.B_TRUE)
	if ret != 0 {
		err = fmt.Errorf("Can't get group quota for gid %v", gid)
	} else {
		quota, _ = strconv.ParseUint(C.GoString(buf), 0, 64)
	}

	return
}

// GetGroupUsed returns group usage
func (z *Zfs) GetGroupUsed(gid uint64) (quota uint64, err error) {
	prop := C.CString(fmt.Sprintf("groupused@%v", gid))
	var buf *C.char
	buf = (*C.char)(C.calloc(C.ZFS_MAXPROPLEN, 1))

	ret := C.zfs_prop_get_userquota(z.zfs, prop, buf, C.ZFS_MAXPROPLEN, C.B_TRUE)
	if ret != 0 {
		err = fmt.Errof("Can't get group used for gid %v", gid)
	} else {
		quota, _ = strconv.ParseUint(C.GoString(buf), 0, 64)
	}

	return
}

// SetGroupQuota sets group quota
func (z *Zfs) SetGroupQuota(gid uint64, quota string) error {
	prop := C.CString(fmt.Sprintf("groupquota@%v", gid))

	ret := C.zfs_prop_set(z.zfs, prop, C.CString(quota))
	if ret != 0 {
		return fmt.Errorf("Unable to set quota '%v' for gid %v", quota, gid)
	}

	return nil
}

// SetUserQuota sets user quota
func (z *Zfs) SetUserQuota(uid uint64, quota string) error {
	prop := C.CString(fmt.Sprintf("userquota@%v", uid))

	ret := C.zfs_prop_set(z.zfs, prop, C.CString(quota))
	if ret != 0 {
		return fmt.Errorf("Unable to set quota '%v' for uid %v", quota, uid)
	}

	return nil
}
