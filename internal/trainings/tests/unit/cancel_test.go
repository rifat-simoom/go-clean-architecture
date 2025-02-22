package unit_test

import (
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/src/domain/training"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTraining_Cancel(t *testing.T) {
	t.Parallel()
	tr := newExampleTraining(t)
	// it's always a good idea to ensure about pre-conditions in the test ;-)
	assert.False(t, tr.IsCanceled())

	err := tr.Cancel()
	require.NoError(t, err)
	assert.True(t, tr.IsCanceled())
}

func TestTraining_Cancel_already_canceled(t *testing.T) {
	t.Parallel()
	tr := newCanceledTraining(t)

	assert.EqualError(t, tr.Cancel(), training.ErrTrainingAlreadyCanceled.Error())
}
