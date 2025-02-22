package application

import (
	command2 "github.com/rifat-simoom/go-clean-architecture/internal/trainings/application/command"
	query2 "github.com/rifat-simoom/go-clean-architecture/internal/trainings/application/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	ApproveTrainingReschedule command2.ApproveTrainingRescheduleHandler
	CancelTraining            command2.CancelTrainingHandler
	RejectTrainingReschedule  command2.RejectTrainingRescheduleHandler
	RescheduleTraining        command2.RescheduleTrainingHandler
	RequestTrainingReschedule command2.RequestTrainingRescheduleHandler
	ScheduleTraining          command2.ScheduleTrainingHandler
}

type Queries struct {
	AllTrainings     query2.AllTrainingsHandler
	TrainingsForUser query2.TrainingsForUserHandler
}
