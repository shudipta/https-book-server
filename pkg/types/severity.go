//go:generate go-enum -f=severity.go --lower --flag
package types

// ref: https://github.com/coreos/clair/blob/a5b3e747a0bffd13929bcbea5bd1d6495cf9f64b/database/severity.go#L27

// Severity defines a standard scale for measuring the severity of a vulnerability.
// ENUM(
// Defcon1
// Critical
// High
// Medium
// Low
// Negligible
// Unknown
// )
type Severity int32
