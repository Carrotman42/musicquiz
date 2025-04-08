// TODO: rename keyboard-specific things with the name keyboard.
package oslevelinput

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"unsafe"
)

// See https://www.kernel.org/doc/Documentation/input/input.txt
// Sizes of the data may change with your version of linux. Send CLs.
type InputEvent struct {
	Seconds      uint32
	Microseconds int32
	Type         EventType
	Code         EventCode
	Value        uint32
}

func (ie InputEvent) String() string {
	if ie.Type == EV_KEY {
		desc := ""
		switch ie.Value {
		case 0:
			desc = "up"
		case 1:
			desc = "down"
		case 2:
			desc = "autorepeat"
		default:
			desc = fmt.Sprintf("unknown(%v)", ie.Value)
		}
		return fmt.Sprintf("%d.%d: %v %s", ie.Seconds, ie.Microseconds, ie.Code, desc)
	}
	return fmt.Sprintf("%d.%d: %v %v %v", ie.Seconds, ie.Microseconds, ie.Type, ie.Code, ie.Value)
}

// I found this hard-coded somewhere else, so it should not be changed. If
// the Sizeof InputEvent changes, there's a problem somewhere else!
const InputEventSize = 16

func init() {
	if sz := unsafe.Sizeof(InputEvent{}); sz != InputEventSize {
		panic(sz)
	}
}

// Open opens all files matching /dev/input/event* and sends events to the
// given channel.  Note that there will necessarily be scheduling latency for
// receiving events; [ReadEventFiles] can be used instead if that's a problem.
func Open() (<-chan InputEvent, error) {
	// TODO: pick a specific event source rather than all of them.
	fs, err := filepath.Glob("/dev/input/event*")
	if err != nil {
		return nil, err
	}
	ch := make(chan InputEvent, 10)
	cb := func(ev InputEvent) {
		ch <- ev
	}
	for _, f := range fs {
		go ReadEventFile(f, cb)
	}
	return ch, nil
}

// ReadEventFile calls cb for each new event from matching /dev/input/event*;
// it may be called in parallel. This function returns after background
// goroutines are started, assuming there is no error.
func ReadEventFiles(cb func(file string, event InputEvent)) error {
	fs, err := filepath.Glob("/dev/input/event*")
	if err != nil {
		return err
	}
	for _, f := range fs {
		go ReadEventFile(f, func(event InputEvent) { cb(f, event) })
	}
	return nil
}

// ReadEventFile calls cb for each new event from srcFile.
// Panics on error.
func ReadEventFile(srcFile string, cb func(InputEvent)) {
	f, err := os.Open(srcFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for {
		var buf [16]byte
		n, err := f.Read(buf[:])
		if err != nil {
			panic(err)
		} else if n != len(buf) {
			// Apparently we can rely on reads always completing in full.
			panic(fmt.Sprintf("only read %d, expected 16", n))
		}

		var event InputEvent
		if err := binary.Read(bytes.NewReader(buf[:]), binary.LittleEndian, &event); err != nil {
			panic(err)
		}
		cb(event)
	}
}

// OpenAllWrite opens all files matching /dev/input/event*, returning a
// mapping to help write events to them. Some valid writers may be returned
// along with an error. Use [OpenWrite] for just a single file.
func OpenAllWrite() (map[string]Writer, error) {
	//fs, err := filepath.Glob("/dev/input/event*")
	fs, err := filepath.Glob("/dev/input/**/*")
	if err != nil {
		return nil, err
	}
	ret := make(map[string]Writer, len(fs))
	var errs []error
	for _, f := range fs {
		wr, err := OpenWrite(f)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", f, err))
			continue
		}
		ret[f] = wr
	}
	if len(ret) == 0 {
		return nil, fmt.Errorf("OpenWrite: couldn't open any input events; errors: %v", errs)
	}
	if len(errs) > 0 {
		return ret, fmt.Errorf("OpenWrite: succeeded %d times, failed %d; errors: %v", len(ret), len(errs), errs)
	}
	return ret, nil
}

func OpenWrite(f string) (Writer, error) {
	file, err := os.OpenFile(f, os.O_WRONLY, 0)
	if err != nil {
		return Writer{}, err
	}
	return Writer{file}, nil
}

// With help from https://forums.raspberrypi.com/viewtopic.php?t=214540.
type Writer struct {
	f *os.File
}

func (wr Writer) Close() error {
	return wr.f.Close()
}

// Ignores errors
func (wr Writer) Hold(code EventCode) (release func()) {
	// SO MANY ERRORS TO CHECK!
	must(writeEvent(wr.f, EV_KEY, code, 1))       // down
	must(writeEvent(wr.f, EV_SYN, SYN_REPORT, 0)) // report
	return func() {
		must(writeEvent(wr.f, EV_KEY, code, 0))       // up
		must(writeEvent(wr.f, EV_SYN, SYN_REPORT, 0)) // report
	}
}

// Ignores errors
func (wr Writer) Keypress(code EventCode) {
	// SO MANY ERRORS TO CHECK!
	must(writeEvent(wr.f, EV_KEY, code, 1))       // down
	must(writeEvent(wr.f, EV_SYN, SYN_REPORT, 0)) // report
	must(writeEvent(wr.f, EV_KEY, code, 0))       // up
	must(writeEvent(wr.f, EV_SYN, SYN_REPORT, 0)) // report
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func writeEvent(wr *os.File, typ EventType, code EventCode, value uint32) error {
	event := InputEvent{
		Seconds:      0, // ignored
		Microseconds: 0, // ignored
		Type:         typ,
		Code:         code,
		Value:        value,
	}
	var buf [InputEventSize]byte
	if n, err := binary.Encode(buf[:], binary.LittleEndian, event); err != nil {
		panic(err)
	} else if n != InputEventSize {
		panic("Wrong written size!")
	}
	_, err := wr.Write(buf[:])
	return err
}

type MouseWriter struct {
	f *os.File
}

func OpenMouseWriter(devFile string) (*MouseWriter, error) {
	file, err := os.OpenFile(devFile, os.O_WRONLY, 0)
	if err != nil {
		return nil, err
	}
	return &MouseWriter{file}, nil
}

func (mw *MouseWriter) Close() error {
	return mw.f.Close()
}

func (mw *MouseWriter) Click(btn EventCode) {
	must(writeEvent(mw.f, EV_KEY, btn, 1))        // down
	must(writeEvent(mw.f, EV_SYN, SYN_REPORT, 0)) // report
	must(writeEvent(mw.f, EV_KEY, btn, 0))        // up
	must(writeEvent(mw.f, EV_SYN, SYN_REPORT, 0)) // report
}

func (mw *MouseWriter) SetCursorLocation(x, y uint32) {
	// TODO: doesn't seem to work!
	//must(writeEvent(mw.f, EV_ABS, ABS_X, x))       // set x
	//must(writeEvent(mw.f, EV_ABS, ABS_Y, y))       // set y
	//must(writeEvent(mw.f, EV_SYN, SYN_REPORT, 0)) // report

	mw.MoveCursor(-10000, -10000)
	mw.MoveCursor(int32(x), int32(y))
}

func (mw *MouseWriter) MoveCursor(x, y int32) {
	must(writeEvent(mw.f, EV_REL, REL_X, uint32(x))) // slide x
	must(writeEvent(mw.f, EV_REL, REL_Y, uint32(y))) // slide y
	must(writeEvent(mw.f, EV_SYN, SYN_REPORT, 0))    // report
}
