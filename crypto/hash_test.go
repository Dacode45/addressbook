package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Dacode45/addressbook/crypto"
)

func Test_Hash(t *testing.T) {
	t.Run("Can generate and compare hashes", should_be_able_to_hash)
	t.Run("Can compare unequal hashes", should_fail_with_different_hashes)
	t.Run("Generates a new salt every time", should_generate_unique_salts)
}

func should_be_able_to_hash(t *testing.T) {
	c := crypto.Hash{}
	input := "testInput"

	generatedHash, generatedError := c.Generate(input)
	compareErr := c.Compare(generatedHash, input)

	assert.NoError(t, generatedError, "Error generating Hash")
	assert.NotEqual(t, generatedHash, input, "Hash is the same as input")
	assert.NoError(t, compareErr, "Comparison should be successful")
}

func should_fail_with_different_hashes(t *testing.T) {
	c := crypto.Hash{}
	input := "testInput"
	compare := "testCompare"

	generatedHash, generatedError := c.Generate(input)
	compareErr := c.Compare(generatedHash, compare)

	assert.NoError(t, generatedError, "Error generating Hash")
	assert.NotEqual(t, generatedHash, input, "Hash is the same as input")
	assert.Error(t, compareErr, "Comparison should not be successful")
}

func should_generate_unique_salts(t *testing.T) {
	c := crypto.Hash{}
	input := "testInput"

	generatedHash, generatedError := c.Generate(input)
	generatedHash2, generatedError2 := c.Generate(input)

	assert.NoError(t, generatedError, "Error generating Hash")
	assert.NoError(t, generatedError2, "Error generating Hash")
	assert.NotEqual(t, generatedHash, generatedHash2, "Same salt generated")
}
