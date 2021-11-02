package main

const (
	ctxId = ctxIdType("id")
)

// ctxIdType is used to avoid collisions between packages using context
type ctxIdType string
