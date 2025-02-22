package application

import (
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/src/application/command"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/src/application/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	ApproveTrainingReschedule command.ApproveTrainingRescheduleHandler
	CancelTraining            command.CancelTrainingHandler
	RejectTrainingReschedule  command.RejectTrainingRescheduleHandler
	RescheduleTraining        command.RescheduleTrainingHandler
	RequestTrainingReschedule command.RequestTrainingRescheduleHandler
	ScheduleTraining          command.ScheduleTrainingHandler
}

type Queries struct {
	AllTrainings     query.AllTrainingsHandler
	TrainingsForUser query.TrainingsForUserHandler
}
