package journal

// #cgo pkg-config: libsystemd
// #include <systemd/sd-journal.h>
// #include <errno.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

type Journal struct {
	c *C.sd_journal
}

type Flag int

const (
	FlagSystem Flag = C.SD_JOURNAL_LOCAL_ONLY | C.SD_JOURNAL_SYSTEM
	FlagUser   Flag = C.SD_JOURNAL_LOCAL_ONLY | C.SD_JOURNAL_CURRENT_USER
)

func New(flag Flag) (*Journal, error) {
	var j = &Journal{}
	if ret := int(C.sd_journal_open(&j.c, C.int(flag))); ret < 0 {
		return nil, fmt.Errorf("failed to open journal : %s", errnoTo(ret))
	}

	// set threshold to unlimited
	if ret := int(C.sd_journal_set_data_threshold(j.c, 0)); ret < 0 {
		return nil, fmt.Errorf("failed to set data threshold : %s", errnoTo(ret))
	}
	return j, nil
}

func (j *Journal) Close() {
	C.sd_journal_close(j.c)
}

func (j *Journal) SeekNext() (int, error) {
	rc := int(C.sd_journal_next(j.c))
	if rc < 0 {
		return rc, fmt.Errorf("failed seek next : %s", errnoTo(rc))
	}
	return rc, nil
}

func (j *Journal) GetCursor() (string, error) {
	var cursor *C.char
	if rc := int(C.sd_journal_get_cursor(j.c, &cursor)); rc < 0 {
		return "", fmt.Errorf("failed get cursor : %s", errnoTo(rc))
	}
	return C.GoString(cursor), nil
}

func (j *Journal) AddMatch(v string) (err error) {
	cstr := C.CString(v)
	defer C.free(unsafe.Pointer(cstr))

	if rc := int(C.sd_journal_add_match(j.c, unsafe.Pointer(cstr), 0)); rc < 0 {
		err = fmt.Errorf("failed add match : %s", errnoTo(rc))
	}
	return
}

func (j *Journal) AddDisjunction() (err error) {
	if rc := int(C.sd_journal_add_disjunction(j.c)); rc < 0 {
		err = fmt.Errorf("failed add disjunction : %s", errnoTo(rc))
	}
	return
}

func (j *Journal) AddConjunction() (err error) {
	if rc := int(C.sd_journal_add_conjunction(j.c)); rc < 0 {
		err = fmt.Errorf("failed add conjunction : %s", errnoTo(rc))
	}
	return
}

func (j *Journal) TestCursor(v string) (err error) {
	cstr := C.CString(v)
	defer C.free(unsafe.Pointer(cstr))

	if rc := int(C.sd_journal_test_cursor(j.c, cstr)); rc < 0 {
		err = fmt.Errorf("failed to test the cursor : %s", errnoTo(rc))
	}
	return
}

func (j *Journal) SeekCursor(v string) (err error) {
	cstr := C.CString(v)
	defer C.free(unsafe.Pointer(cstr))

	if rc := int(C.sd_journal_seek_cursor(j.c, cstr)); rc < 0 {
		err = fmt.Errorf("failed to seek to the cursor : %s", errnoTo(rc))
	}
	return
}

func (j *Journal) Getdata(name string) (string, error) {
	var cstr *C.char = C.CString(name)
	defer C.free(unsafe.Pointer(cstr))

	var (
		s unsafe.Pointer
		n C.size_t
		e int
	)

	e = int(C.sd_journal_get_data(j.c, cstr, &s, &n))

	if e == -C.ENOENT {
		return "", nil
	}

	if e < 0 {
		return "", fmt.Errorf("failed to get %s: %s", name, errnoTo(e))
	}

	ret := C.GoStringN((*C.char)(s), C.int(n))

	if after, ok := strings.CutPrefix(ret, name+"="); ok {
		return after, nil
	}

	return ret, nil
}

/* see sd_journal_get_data(3) */
var errMap = map[int]string{
	C.EINVAL:          "One of the required parameters is NULL or invalid",
	C.ECHILD:          "The journal object was created in a different process, library or module instance",
	C.EADDRNOTAVAIL:   "The read pointer is not positioned at a valid entry",
	C.ENOENT:          "The current entry does not include the specified field",
	C.ENOMEM:          "Memory allocation failed",
	C.ENOBUFS:         "A compressed entry is too large",
	C.E2BIG:           "The data field is too large for this computer architecture",
	C.EPROTONOSUPPORT: "The journal is compressed with an unsupported method or the journal uses an unsupported feature",
	C.EBADMSG:         "The journal is corrupted",
	C.EIO:             "An I/O error was reported by the kernel",
}

func errnoTo(rc int) string {
	if s := errMap[-1*rc]; s != "" {
		return s
	}

	return syscall.Errno(-rc).Error()
}
