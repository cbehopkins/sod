package sod

import (
	"errors"
	"strconv"
)

type Value int

func (v Value) String() string {
	return strconv.Itoa(int(v))
}

const PuzzleSize = 9

var ErrUnfoundValue = errors.New("Value not found in array")

// Hidden Twin(Triplet?) Processing
// In a group Look for Twins
// A twin is a pair of numbers that appear in cell n and m
// but do not appear in any other cells in that group
// In that case any other numbers in cells n and m can be removed.
// This will produce a set of naked Twins
// Increase the difficulty every time we use this

// Naked Twin processing
// First Find the Naked twins
// In that case Eliminate the numbers that formt he naked twins
// From all other cells in the group
// Increase the difficulty for having used this

// X-Wing processing
// Look for a pair of numbers that appears twice in a row
// Record the co-ords this happens at
// If the same pair of numbers appear twice again in another row
// Check to see if the co-ords for this happen in the same columns
// Should this happen we have an X-Wing
// In which case For the row an column it is safe to remove both
// Numbers from cells not in the X-Wing
// Really increase the difficulty for having used this

// Puzzle Generation
// Start from a blank grid
// Pick a random co-ord and pick one of the possible values to
// set this cell to
// Run the solver & is it solved
// Repeat until solved.
