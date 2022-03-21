package shm

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"golang.org/x/sys/unix"
)

//SharedMemory is a Wrapper around SysVShm, and is not safe for concurrent access
type SharedMemory struct {
	id   int
	size int
	b    []byte
}

//New SharedMemory using id, if id is 0, make one there
func New(id int, size int) *SharedMemory {
	return &SharedMemory{id: id, size: size}
}

var ErrorGivenSliceTooBig = fmt.Errorf("shm: slice of bytes provided too big to write")

func (sm *SharedMemory) ReadAt(p []byte, off int64) (n int, err error) {
	if int64(len(p))+off > int64(sm.size) {
		return 0, ErrorGivenSliceTooBig
	}
	copy(p, sm.b[off:len(p)+int(off)])
	return len(p), nil
}

func (sm *SharedMemory) WriteAt(p []byte, off int64) (n int, err error) {
	if int64(len(p))+off > int64(sm.size) {
		return 0, ErrorGivenSliceTooBig
	}
	copy(sm.b[off:len(p)+int(off)], p)
	return len(p), nil
}

//Open the SharedMemory (shm) segment, will return nil if called while already open
func (sm *SharedMemory) Open() error {
	var err error
	if sm.id == 0 {
		sm.id, err = unix.SysvShmGet(unix.IPC_PRIVATE, sm.size, unix.IPC_CREAT|unix.IPC_EXCL|0o600)
		if err != nil {
			return fmt.Errorf("sysvshmget failed: (%w)", err)
		}
	}

	sm.b, err = unix.SysvShmAttach(sm.id, 0, 0)
	if err != nil {
		err2 := sm.Close()
		if err2 != nil {
			return multierror.Append(err, err2)
		}
	}
	return nil
}

//Close *MUST* be called after SharedMemory is finished being used
func (sm *SharedMemory) Close() error {
	_, err := unix.SysvShmCtl(sm.id, unix.IPC_RMID, nil)
	if err != nil {
		return fmt.Errorf("sysvshmctl close failed (%w)", err)
	}
	return nil
}

func (sm *SharedMemory) ID() int {
	return sm.id
}

func (sm *SharedMemory) Size() int {
	return sm.size
}
