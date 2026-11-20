package worker

type Worker interface {
	GenerateQuery() string
}
