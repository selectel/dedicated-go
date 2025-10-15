// Package httptest provides helpers for testing HTTP clients.
//
// It offers:
//   - RoundTripFunc to stub http.RoundTripper
//   - NewFakeResponse to build *http.Response with a status/body
//   - NewFakeTransport to return a transport that yields a fixed response/error
//
// These utilities simplify unit tests by avoiding real network calls.
package httptest
