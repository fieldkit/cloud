package tests

import (
	"log"
	"os"
)

type MigratorLog struct{}

func (l *MigratorLog) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l *MigratorLog) Println(args ...interface{}) {
	log.Println(args...)
}

func (l *MigratorLog) Verbose() bool {
	return false
}

func (l *MigratorLog) fatal(args ...interface{}) {
	l.Println(args...)
	os.Exit(1)
}

func (l *MigratorLog) fatalErr(err error) {
	l.fatal("error:", err)
}
