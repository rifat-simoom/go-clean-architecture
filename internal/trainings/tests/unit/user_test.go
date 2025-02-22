package unit_test

import (
	training2 "github.com/rifat-simoom/go-hexarch/internal/trainings/src/domain/training"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsUserAllowedToSeeTraining(t *testing.T) {
	t.Parallel()
	attendee1, err := training2.NewUser(uuid.New().String(), training2.Attendee)
	require.NoError(t, err)

	attendee2, err := training2.NewUser(uuid.New().String(), training2.Attendee)
	require.NoError(t, err)

	trainer, err := training2.NewUser(uuid.New().String(), training2.Trainer)
	require.NoError(t, err)

	testCases := []struct {
		Name              string
		CreateTraining    func(t *testing.T) *training2.Training
		User              training2.User
		ExpectedIsAllowed bool
	}{
		{
			Name: "attendees_training",
			CreateTraining: func(t *testing.T) *training2.Training {
				tr, err := training2.NewTraining(
					uuid.New().String(),
					attendee1.UUID(),
					"user name",
					time.Now(),
				)
				require.NoError(t, err)

				return tr
			},
			User:              attendee1,
			ExpectedIsAllowed: true,
		},
		{
			Name: "another_attendees_training",
			CreateTraining: func(t *testing.T) *training2.Training {
				tr, err := training2.NewTraining(
					uuid.New().String(),
					attendee1.UUID(),
					"user name",
					time.Now(),
				)
				require.NoError(t, err)

				return tr
			},
			User:              attendee2,
			ExpectedIsAllowed: false,
		},
		{
			Name: "trainer",
			CreateTraining: func(t *testing.T) *training2.Training {
				tr, err := training2.NewTraining(
					uuid.New().String(),
					attendee1.UUID(),
					"user name",
					time.Now(),
				)
				require.NoError(t, err)

				return tr
			},
			User:              trainer,
			ExpectedIsAllowed: true, // trainer have access to all trainings
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			tr := c.CreateTraining(t)

			err := training2.CanUserSeeTraining(c.User, *tr)

			if c.ExpectedIsAllowed {

			} else {
				assert.EqualError(
					t,
					err,
					training2.ForbiddenToSeeTrainingError{c.User.UUID(), tr.UserUUID()}.Error(),
				)
			}
		})
	}
}
