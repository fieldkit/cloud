package ingester

type Config struct {
	ProductionLogging bool   `envconfig:"production_logging"`
	Addr              string `split_words:"true" default:"127.0.0.1:8080" required:"true"`
	PostgresURL       string `split_words:"true" default:"postgres://localhost/fieldkit?sslmode=disable" required:"true"`
	AwsProfile        string `envconfig:"aws_profile" default:"fieldkit" required:"true"`
	AwsId             string `split_words:"true" default:""`
	AwsSecret         string `split_words:"true" default:""`
	Archiver          string `split_words:"true" default:"default" required:"true"`
	SessionKey        string `split_words:"true"`
	StatsdAddress     string `split_words:"true" default:""`
	StreamsBucketName string `split_words:"true" default:""`
	Help              bool
}
