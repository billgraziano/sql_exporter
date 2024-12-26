FORK Details
============
* Always display the time stamp.  By default, they seem to collect every 10 seconds.  I only want to poll the servers every minute.  I changed so that it always displays the timestamp. This doesn't seem to affect `/sql_exporter_metrics` (which is good).  Ideally this would be a global configuration flag.  Or even a per query configuration flag.

```go
// metric.go:87
value := row[v].(sql.NullFloat64)
if value.Valid {
    metric := NewMetric(&mf, value.Float64, labelValues...)
    // if mf.config.TimestampValue == "" {
    // 	ch <- metric
    // } else {
    // 	ts := row[mf.config.TimestampValue].(sql.NullTime)
    // 	if ts.Valid {
    // 		ch <- NewMetricWithTimestamp(ts.Time, metric)
    // 	}
    // }
    // BG: Always send a timestamp
    ch <- NewMetricWithTimestamp(time.Now(), metric)
}
```

* Log file is renamed with YYYY_MM in the file name.  The log file is still appended if it exists.

```go
// log.go:19
// initLogFile opens the log file for writing if a log file is specified.
func initLogFile(logFile string) (*os.File, error) {
	if logFile == "" {
		return nil, nil
	}
	logFile = logFileWithTimeStamp(logFile)
	logFileHandler, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %w", err)
	}
	return logFileHandler, nil
}

// logFileWithTimeStamp inserts a YYYY_MM timestamp in a file name
func logFileWithTimeStamp(logFile string) string {
	ext := filepath.Ext(logFile)
	name := strings.TrimSuffix(logFile, ext)
	tsname := fmt.Sprintf("%s_%d_%d%s", name, time.Now().Year(), time.Now().Month(), ext)
	return tsname
}
```
