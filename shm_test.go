package shm

import (
	"errors"
	"github.com/hashicorp/go-multierror"
	"io"
	"reflect"
	"syscall"
	"testing"
)

func TestSharedMemory_ReadAt(t *testing.T) {
	t.Parallel()
	const size = 512
	shm := New(0, size)
	err := shm.Open()
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		if err == nil {
			return
		}
		err2 := shm.Close()
		if err2 != nil {
			err = multierror.Append(err, err2)
			t.Error(err)
		}
	}()

	for i := 0; i < size; i++ {
		shm.b[i] = 1
	}

	sr := io.NewSectionReader(shm, 0, size)
	bytes, err := io.ReadAll(sr)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(shm.b, bytes) {
		t.Error("not DeepEqual")
		return
	}

	_, err = shm.ReadAt(bytes, size)
	if err != ErrorGivenSliceTooBig {
		t.Errorf("error not ErrorGivenSliceTooBig is (%v)\n", err)
		return
	}

	err = shm.Close()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSharedMemory_Open(t *testing.T) {
	t.Parallel()
	shm := New(0, 512)
	err := shm.Open()
	if err != nil {
		t.Error(err)
		return
	}
	err2 := shm.Close()
	if err2 != nil {
		err = multierror.Append(err, err2)
		t.Error(err)
		return
	}

	shm = New(0, 0)
	err = shm.Open()
	if err == nil {
		err = shm.Close()
		if err != nil {
			t.Error(err)
			return
		}
	}
	if !errors.Is(err, syscall.Errno(22)) {
		t.Errorf("err type (%t)", err)
		return
	}

	shm = New(-1, 0)
	err = shm.Open()
	if err == nil {
		err = shm.Close()
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func TestSharedMemory_ID(t *testing.T) {
	t.Parallel()
	shm := New(0, 256)
	err := shm.Open()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err2 := shm.Close()
		if err2 != nil {
			err = multierror.Append(err, err2)
			t.Error(err)
		}
	}()
	if shm.id != shm.ID() {
		t.Fail()
		return
	}
}

func TestSharedMemory_Size(t *testing.T) {
	t.Parallel()
	shm := New(0, 256)
	err := shm.Open()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		err2 := shm.Close()
		if err2 != nil {
			err = multierror.Append(err, err2)
			t.Error(err)
		}
	}()
	if shm.size != shm.Size() {
		t.Fail()
		return
	}
}

func TestSharedMemory_Close(t *testing.T) {
	t.Parallel()
	shm := New(-1, 256)
	err := shm.Close()
	if err == nil {
		t.Error("err == nil")
		return
	}
}

func TestSharedMemory_WriteAt(t *testing.T) {
	t.Parallel()
	const size = 512
	shm := New(0, size)
	err := shm.Open()
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		if err == nil {
			return
		}
		err2 := shm.Close()
		if err2 != nil {
			err = multierror.Append(err, err2)
			t.Error(err)
		}
	}()

	bCopy := make([]byte, size)

	for i := 0; i < size; i++ {
		shm.b[i] = 1
	}

	_, err = shm.WriteAt(bCopy, 0)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(shm.b, bCopy) {
		t.Error("shm: reflect.DeepEqual returned false for shm.b and bCopy")
		return
	}

	_, err = shm.WriteAt(bCopy, size)
	if err != ErrorGivenSliceTooBig {
		t.Errorf("error not ErrorGivenSliceTooBig is (%v)\n", err)
		return
	}

	err = shm.Close()
	if err != nil {
		t.Error(err)
		return
	}
}
