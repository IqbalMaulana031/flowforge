package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPaginationNormalizesValues(t *testing.T) {
	pagination := NewPagination(0, 200)

	require.Equal(t, 1, pagination.Page)
	require.Equal(t, 100, pagination.Limit)
	require.Equal(t, 0, pagination.Offset())
}
