package shm

import (
	"bytes"
	"errors"
	"github.com/hashicorp/go-multierror"
	"io"
	"syscall"
	"testing"
)

func FillSlice(p []byte, b byte) {
	for i := range p {
		p[i] = b
	}
}

func TestFillSlice(t *testing.T) {
	t.Parallel()
	p := make([]byte, 64)
	FillSlice(p, 1)
	for _, b := range p {
		if b != 1 {
			t.Fatal("b not 1")
		}
	}
}

func CompareSlices(p []byte, q []byte) bool {
	if len(p) != len(q) {
		return false
	}
	for i, b := range p {
		if b != q[i] {
			return false
		}
	}
	return true
}

func TestCompareSlices(t *testing.T) {
	t.Parallel()

	p := make([]byte, 64)
	FillSlice(p, 1)

	q := make([]byte, 64)
	FillSlice(q, 0)

	if CompareSlices(p, q) {
		t.Fatal("p and q are equal")
	}

	FillSlice(q, 1)
	if !CompareSlices(p, q) {
		t.Fatal("p and q are not equal")
	}
}

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

	FillSlice(shm.b, 1)

	sr := io.NewSectionReader(&shm, 0, size)
	b, err := io.ReadAll(sr)
	if err != nil {
		t.Error(err)
		return
	}

	if !CompareSlices(shm.b, b) {
		t.Error("shm.b and b are not equal in content or length")
		return
	}

	_, err = shm.ReadAt(b, size)
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

	FillSlice(shm.b, 1)

	_, err = shm.WriteAt(bCopy, 0)
	if err != nil {
		t.Error(err)
		return
	}

	if !CompareSlices(shm.b, bCopy) {
		t.Error("shm.b and bCopy are not equal in content and/or size")
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

func TestSharedMemory_Read(t *testing.T) {
	t.Parallel()
	const size = 512
	shm := New(0, size)
	err := shm.Open()
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		err = shm.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	FillSlice(shm.b, 1)

	b, err := io.ReadAll(&shm)
	if err != nil {
		t.Error(err)
		return
	}

	if !CompareSlices(b, shm.b) {
		t.Error("b is not equal to shm.b in content and/or length")
		return
	}

	b = make([]byte, size+1)
	if _, err = shm.Read(b); err != ErrorGivenSliceBiggerThanData {
		t.Errorf("shm.Read(b) did not return ErrorGivenSliceBiggerThanData")
		return
	}
}

func TestSharedMemory_Write(t *testing.T) {
	t.Parallel()
	const size = 512
	shm := New(0, size)
	err := shm.Open()
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		err = shm.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	p := make([]byte, size)
	FillSlice(p, 1)

	written, err := io.Copy(&shm, bytes.NewReader(p))
	if err != nil || written != int64(len(p)) {
		t.Errorf("err (%v) written (%v)\n", err, written)
		return
	}

	if !CompareSlices(p, shm.b) {
		t.Error("p and shm.b are not equal in content and/or length")
		return
	}

	p = make([]byte, size+1)
	if _, err = shm.Write(p); err != ErrorGivenSliceBiggerThanData {
		t.Errorf("shm.Write(p) did not return ErrorGivenSliceBiggerThanData")
	}
}
