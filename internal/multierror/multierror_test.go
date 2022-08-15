package multierror

import (
	"errors"
	"fmt"
	"testing"
)

var errInternalSingleLevel = fmt.Errorf("internalErrorSingleLevel")
var errInternalDoubleLevel = fmt.Errorf("internalErrorDoubleLevel")
var errInternalTripleLevel = fmt.Errorf("internalErrorTripleLevel")

var multiError = MultiError{
	fmt.Errorf("top level error1 (%w)", errInternalSingleLevel),
	fmt.Errorf("top level error2"),
	fmt.Errorf("top level error3"),
}

const multiErrorString string = `top level error1 (internalErrorSingleLevel), top level error2, top level error3`

func TestMultiErrorError(t *testing.T) {
	t.Parallel()
	if multiError.Error() != multiErrorString {
		t.Fatalf("multiError.Error() is not multiErrorString, wanted (%v), got (%v)\n", multiErrorString, multiError.Error())
	}
}

func TestMultiErrorIs(t *testing.T) {
	t.Parallel()
	if !errors.Is(multiError, errInternalSingleLevel) {
		t.Fatal("errors.Is multiError interalErrorSingleLevel is false")
	}
	if errors.Is(multiError, fmt.Errorf("fake error")) {
		t.Fatal("errors.Is multiError fake error is true")
	}
}

func TestAppend(t *testing.T) {
	t.Parallel()
	err := Append(errInternalSingleLevel, errInternalDoubleLevel)
	_, ok := err.(MultiError)
	if !ok {
		t.Fatal("returned error is not a MultiError")
	}

	if !errors.Is(err, errInternalSingleLevel) {
		t.Fatalf("err.Is not a errInternalSingleLevel instead a (%v)\n", err)
	}

	if !errors.Is(err, errInternalDoubleLevel) {
		t.Fatalf("err.Is not a errInternalDoubleLevel instead a (%v)\n", err)
	}

	err = Append(err, errInternalTripleLevel)
	if !errors.Is(err, errInternalTripleLevel) {
		t.Fatal("err.Is not a errInternalTripleLevel after appending")
	}

	err = nil
	err = Append(err, errInternalSingleLevel)
	if !errors.Is(err, errInternalSingleLevel) {
		t.Fatal("err.Is not a errInternalSingleLevel")
	}

	err2 := Append(err, nil)
	if err != err2 {
		t.Fatal("err2 != err")
	}
}
