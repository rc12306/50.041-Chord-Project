package main

import "testing"

// prototype to demonstrate "testing" with scan_ring.go
func TestScanRing(t *testing.T) {
  ipInRing := CheckRing()
  // check if ipInRing is empty
  if len(ipInRing) == 0 {
    t.Errorf("No nodes found in chord ring )-:")
  }
}

// TODO:
// Debug errors found when using
// go test -run TestScanRing
