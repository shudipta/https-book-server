//go:generate go-enum -f=failurepolicy.go --lower --flag
package types

// FailurePolicy defines how errors from the docker registry are handled
// allowed values are Ignore or Fail. Defaults to Ignore.
// ENUM(
// Ignore
// Fail
// )
type FailurePolicy int32
