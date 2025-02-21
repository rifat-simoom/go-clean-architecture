module github.com/rifat-simoom/go-clean-architecture/internal/c4

go 1.16

require (
	github.com/rifat-simoom/go-clean-architecture/internal/trainer v0.0.0-00010101000000-000000000000
	github.com/rifat-simoom/go-clean-architecture/internal/trainings v0.0.0-00010101000000-000000000000
	github.com/krzysztofreczek/go-structurizr v0.1.2
)

replace (
	github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel => ../../internal/shared_kernel/
	github.com/rifat-simoom/go-clean-architecture/internal/trainer => ../../internal/trainer/
	github.com/rifat-simoom/go-clean-architecture/internal/trainings => ../../internal/trainings/
)
