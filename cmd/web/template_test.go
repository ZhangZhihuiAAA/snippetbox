package main

import (
	"snippetbox/internal/assert"
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
    tests := []struct {
        name     string
        input    time.Time
        expected string
    }{
        {
            name:     "UTC",
            input:    time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC),
            expected: "17 Mar 2024 at 10:15",
        },
        {
            name:     "Empty",
            input:    time.Time{},
            expected: "",
        },
        {
            name:     "CET",
            input:    time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
            expected: "17 Mar 2024 at 09:15",
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            hd := humanDate(tc.input)

            // Use the new assert.Equal() helper to compare the actual and expected values.
            assert.Equal(t, hd, tc.expected)
        })
    }
}
