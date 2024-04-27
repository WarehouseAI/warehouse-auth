package logger

import "log"

// EmptyLogger for testing pg without logging goose migrations
type EmptyLogger struct {
}


func (e *EmptyLogger) Fatal(v ...interface{}) {
	log.Fatal(v...)
}
func (e *EmptyLogger) Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}
func (e *EmptyLogger) Print(v ...interface{}) {

}
func (e *EmptyLogger) Println(v ...interface{}) {

}
func (e *EmptyLogger) Printf(format string, v ...interface{}) {
}

