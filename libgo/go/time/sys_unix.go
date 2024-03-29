// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin freebsd linux netbsd openbsd

package time

import "syscall"

// for testing: whatever interrupts a sleep
func interrupt() {
	syscall.Kill(syscall.Getpid(), syscall.SIGCHLD)
}

// readFile reads and returns the content of the named file.
// It is a trivial implementation of ioutil.ReadFile, reimplemented
// here to avoid depending on io/ioutil or os.
func readFile(name string) ([]byte, error) {
	f, err := syscall.Open(name, syscall.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.Close(f)
	var (
		buf [4096]byte
		ret []byte
		n   int
	)
	for {
		n, err = syscall.Read(f, buf[:])
		if n > 0 {
			ret = append(ret, buf[:n]...)
		}
		if n == 0 || err != nil {
			break
		}
	}
	return ret, err
}
