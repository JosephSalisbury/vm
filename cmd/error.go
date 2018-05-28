package cmd

import "errors"

var (
	// NoVMError is the error returned if there are no VMs.
	NoVMError = errors.New("no vm")
	// MultipleVMError is the error returned if there are multiple VMs,
	// and we need to specify just one.
	MultipleVMError = errors.New("multiple vm")
)
