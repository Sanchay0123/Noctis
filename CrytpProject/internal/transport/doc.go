// Package transport manages peer TCP connections and strict length-prefixed framing.
// It enforces maximum frame size limits prior to allocation and handles low-level I/O.
// It operates purely on framed bytes and has no access to or knowledge of application plaintext.
package transport
