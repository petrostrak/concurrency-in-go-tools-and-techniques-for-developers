package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

type MyError struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

func wrapError(err error, messagef string, msgArgs ...interface{}) MyError {
	return MyError{

		// Here we store the error we're wrapping. We always want to be able to get back to
		// the lowest-level error in case we need to investigate what happened.
		Inner:   err,
		Message: fmt.Sprintf(messagef, msgArgs...),

		// This line of code takes note of the stack trace when the error was created. A more
		// sophisticated error type might elide the stack-frame from wrapError.
		StackTrace: string(debug.Stack()),

		// Here we create a catch-all for storing miscelleneous information. This is where we
		// might store the concurrent ID, a hash of the stack trace, or other contextual info
		// that might help in diagnosing the error.
		Misc: map[string]interface{}{},
	}
}

func (err MyError) Error() string {
	return err.Message
}

// "low-level" module

type LowLevelErr struct {
	error
}

func isGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {

		// Here we wrap the raw error from calling os.Stat with a customized error. In this
		// case we are OK with the message comming out of this error, and so we won't mask it.
		return false, LowLevelErr{wrapError(err, err.Error())}
	}

	// Checks if file is executable by its owner (we use bitmask 0100)
	return info.Mode().Perm()&0100 == 100, nil
}

// "intermediate" module

type IntermediateErr struct {
	error
}

func runJob(id string) error {
	const jobBinPath = "bad/job/binary"
	isExecutable, err := isGloballyExec(jobBinPath)
	if err != nil {
		return err
	} else if !isExecutable {
		return wrapError(nil, "job binary is not executable")
	}

	// Here we are passing on errors from the lowlevel module. Because of our architectural
	// decision to consider errors passed on from other modules without wrapping them in our
	// own type bugs, this will cause us issues later.
	return exec.Command(jobBinPath, "--id="+id).Run()
}

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))

	// Here we log out the full error in case someone needs to dig into what happened.
	log.Printf("%#v", err)
	fmt.Printf("[%v] %v", key, message)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	if err := runJob("1"); err != nil {
		msg := "There was an unexpected issue; please report this as a bug."

		// Here we check to see if the error is of the expected type. If it is, we know
		// it's a weel-crafted error, and we can simply pass its message on to the user.
		if _, ok := err.(IntermediateErr); ok {
			msg = err.Error()
		}

		// On this line we bind the log and error message together with an ID of 1. We could
		// easily make this increase monotonically, or use a GUID to ensure a unique ID.
		handleError(1, err, msg)
	}
}
