module github.com/rifat-simoom/go-hexarch/internal/c4

go 1.16

require (
	github.com/rifat-simoom/go-hexarch/internal/trainer v0.0.0-00010101000000-000000000000
	github.com/rifat-simoom/go-hexarch/internal/trainings v0.0.0-00010101000000-000000000000
	github.com/krzysztofreczek/go-structurizr v0.1.2
)

replace (
	github.com/rifat-simoom/go-hexarch/internal/shared_kernel => ../../internal/shared_kernel/
	github.com/rifat-simoom/go-hexarch/internal/trainer => ../../internal/trainer/
	github.com/rifat-simoom/go-hexarch/internal/trainings => ../../internal/trainings/
)
