package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var windowPairs = []string{
	"ALPHA=1",
	"BETA=2",
	"GAMMA=3",
	"DELTA=4",
	"EPSILON=5",
}

func TestWindow_NoneMode_Passthrough(t *testing.T) {
	w, err := NewWindower(WindowModeNone, 0, 0)
	require.NoError(t, err)
	got := w.Apply(windowPairs)
	assert.Equal(t, windowPairs, got)
}

func TestWindow_UnknownMode_Passthrough(t *testing.T) {
	w, err := NewWindower(WindowMode("bogus"), 0, 3)
	require.NoError(t, err)
	got := w.Apply(windowPairs)
	assert.Equal(t, windowPairs, got)
}

func TestWindow_HeadMode_ReturnsFirstN(t *testing.T) {
	w, err := NewWindower(WindowModeHead, 0, 3)
	require.NoError(t, err)
	got := w.Apply(windowPairs)
	assert.Equal(t, []string{"ALPHA=1", "BETA=2", "GAMMA=3"}, got)
}

func TestWindow_HeadMode_SizeExceedsLen_Passthrough(t *testing.T) {
	w, err := NewWindower(WindowModeHead, 0, 100)
	require.NoError(t, err)
	got := w.Apply(windowPairs)
	assert.Equal(t, windowPairs, got)
}

func TestWindow_TailMode_ReturnsLastN(t *testing.T) {
	w, err := NewWindower(WindowModeTail, 0, 2)
	require.NoError(t, err)
	got := w.Apply(windowPairs)
	assert.Equal(t, []string{"DELTA=4", "EPSILON=5"}, got)
}

func TestWindow_TailMode_ZeroSize_Passthrough(t *testing.T) {
	w, err := NewWindower(WindowModeTail, 0, 0)
	require.NoError(t, err)
	got := w.Apply(windowPairs)
	assert.Equal(t, windowPairs, got)
}

func TestWindow_SliceMode_OffsetAndSize(t *testing.T) {
	w, err := NewWindower(WindowModeSlice, 1, 3)
	require.NoError(t, err)
	got := w.Apply(windowPairs)
	assert.Equal(t, []string{"BETA=2", "GAMMA=3", "DELTA=4"}, got)
}

func TestWindow_SliceMode_OffsetBeyondEnd_ReturnsNil(t *testing.T) {
	w, err := NewWindower(WindowModeSlice, 99, 3)
	require.NoError(t, err)
	got := w.Apply(windowPairs)
	assert.Nil(t, got)
}

func TestWindow_SliceMode_ZeroSize_ReturnsToEnd(t *testing.T) {
	w, err := NewWindower(WindowModeSlice, 3, 0)
	require.NoError(t, err)
	got := w.Apply(windowPairs)
	assert.Equal(t, []string{"DELTA=4", "EPSILON=5"}, got)
}

func TestWindow_SliceMode_NegativeOffset_ReturnsError(t *testing.T) {
	_, err := NewWindower(WindowModeSlice, -1, 2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "offset")
}

func TestWindow_SliceMode_NegativeSize_ReturnsError(t *testing.T) {
	_, err := NewWindower(WindowModeSlice, 0, -1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "size")
}

func TestWindow_EmptyInput_ReturnsEmpty(t *testing.T) {
	w, err := NewWindower(WindowModeHead, 0, 3)
	require.NoError(t, err)
	got := w.Apply([]string{})
	assert.Empty(t, got)
}
