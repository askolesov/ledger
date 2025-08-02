package model_v1

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransaction_Validate(t *testing.T) {
	tr := Transaction{
		Comment: "test",
	}

	err := tr.Validate()
	require.Error(t, err)
}
