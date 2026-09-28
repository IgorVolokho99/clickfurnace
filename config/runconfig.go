package config

type ClickHouseConfiguration struct {
	host     string
	port     int
	dbname   string
	username string
}

type TableConfiguration struct {
	name    string
	columns []ColumnConfiguration
}

type ColumnConfiguration struct {
	name       string
	columntype string
}

type WorkloadConfiguration struct {
	rate      int
	batchsize int
	workers   int
	duration  int
}

type Config struct {
}
