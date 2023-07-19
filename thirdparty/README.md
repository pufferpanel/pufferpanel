This is a collection of third party libraries and/or Golang files which require
alterations for misc purposes. Licenses of these libraries permit such modifications
and as such, are being done here.

Full disclosure of changes being applied:

fmt
---
- Source: https://go.googlesource.com/go/+/refs/tags/go1.20.6/src/fmt/
- License: https://go.googlesource.com/go/+/refs/tags/go1.20.6/LICENSE

Removed the io and os imports to allow for building for WASM without producing
excessively large builds

fmtsort
---
- Source: https://go.googlesource.com/go/+/refs/tags/go1.20.6/src/internal/fmtsort/
- License: https://go.googlesource.com/go/+/refs/tags/go1.20.6/LICENSE

This is an internal library in Go for usage by fmt, included in our package for 
consumption by altered fmt package

cel-go
------
- Source: https://github.com/google/cel-go/tree/v0.17.0
- License: https://github.com/google/cel-go/blob/v0.17.0/LICENSE

Change fmt usage to altered fmt package to allow for building for WASM without 
producing excessively large builds
