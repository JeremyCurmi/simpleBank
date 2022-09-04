package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	password := "test"
	hashPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.True(t, CheckPasswordHash(password, string(hashPassword)))
}
