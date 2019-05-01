package xerrorz

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"golang.org/x/xerrors"
)

func TestErrQueue0(t *testing.T) {
	e1 := errors.New("e1")

	eq := ErrQueue{Curr: e1, Queue: nil} // Queue is nil

	errStr := fmt.Sprintf("%+v", eq)
	errLines := strings.Split(errStr, "\n")

	if len(errLines) != 1 {
		t.Errorf("Invalid length: %d\n", len(errLines))
	}

	if !strings.Contains(errLines[0], "e1") {
		t.Errorf("l.1 should contain \"e1\"\n")
	}
}

func TestErrQueue1(t *testing.T) {
	e1 := errors.New("e1")

	eq := ErrQueue{Curr: e1, Queue: []*error{}} // Empty queue

	errStr := fmt.Sprintf("%+v", eq)
	errLines := strings.Split(errStr, "\n")

	if len(errLines) != 1 {
		t.Errorf("Invalid length: %d\n", len(errLines))
	}

	if !strings.Contains(errLines[0], "e1") {
		t.Errorf("l.1 should contain \"e1\"\n")
	}
}

func TestErrQueue2(t *testing.T) {
	e1 := errors.New("e1")
	e2 := errors.New("e2")

	eq := ErrQueue{Curr: e1, Queue: []*error{&e2}} // Queue has a single element

	errStr := fmt.Sprintf("%+v", eq)
	errLines := strings.Split(errStr, "\n")

	if len(errLines) != 2 {
		t.Errorf("Invalid length: %d\n", len(errLines))
	}

	if !strings.Contains(errLines[0], "e1") {
		t.Errorf("l.1 should contain \"e1\"\n")
	}

	if !strings.Contains(errLines[1], "e2") {
		t.Errorf("l.2 should contain \"e2\"\n")
	}
}

func TestErrQueue3(t *testing.T) {
	e1 := errors.New("e1")
	e2 := errors.New("e2")
	e3 := errors.New("e3")

	eq := ErrQueue{Curr: e1, Queue: []*error{&e2, &e3}} // Queue has multiple elements

	errStr := fmt.Sprintf("%+v", eq)
	errLines := strings.Split(errStr, "\n")

	if len(errLines) != 3 {
		t.Errorf("Invalid length: %d\n", len(errLines))
	}

	if !strings.Contains(errLines[0], "e1") {
		t.Errorf("l.1 should contain \"e1\"\n")
	}

	if !strings.Contains(errLines[1], "e2") {
		t.Errorf("l.2 should contain \"e2\"\n")
	}

	if !strings.Contains(errLines[2], "e3") {
		t.Errorf("l.3 should contain \"e3\"\n")
	}
}

func TestErrQueue4(t *testing.T) {
	// Nested errors by xerrors.Wrapper
	e1 := xerrors.Errorf("e1: %w", io.ErrClosedPipe)
	e2 := xerrors.Errorf("e2: %w", e1)
	e3 := xerrors.Errorf("e3: %w", e2)
	f1 := xerrors.Errorf("f1: %w", io.ErrNoProgress)
	g1 := xerrors.Errorf("g1")

	eq := ErrQueue{Curr: e3, Queue: []*error{&f1, &g1}}

	errStr := fmt.Sprintf("%+v", eq)
	errLines := strings.Split(errStr, "\n")

	if len(errLines) != 17 {
		t.Errorf("Invalid length: %d\n", len(errLines))
	}

	if !strings.Contains(errLines[0], "e3") {
		t.Errorf("l.1 should contain \"e3\"\n")
	}

	if !strings.Contains(errLines[3], "e2") {
		t.Errorf("l.4 should contain \"e2\"\n")
	}

	if !strings.Contains(errLines[6], "e1") {
		t.Errorf("l.7 should contain \"e1\"\n")
	}

	if !strings.Contains(errLines[9], io.ErrClosedPipe.Error()) {
		t.Errorf("l.10 should contain \"%s\"\n", io.ErrClosedPipe.Error())
	}

	if !strings.Contains(errLines[10], "f1") {
		t.Errorf("l.11 should contain \"f1\"\n")
	}

	if !strings.Contains(errLines[13], io.ErrNoProgress.Error()) {
		t.Errorf("l.14 should contain \"%s\"\n", io.ErrNoProgress.Error())
	}

	if !strings.Contains(errLines[14], "g1") {
		t.Errorf("l.15 should contain \"g1\"\n")
	}
}
