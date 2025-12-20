package main

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/billgraziano/sql_exporter"
	"github.com/fsnotify/fsnotify"
)

// watches configuration files and reloads collectors
func watchConfig(e sql_exporter.Exporter, configFile string) error {
	folder := filepath.Dir(configFile)
	slog.Info("watchconfig", "folder", folder)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("fsnotify.newwatcher: %w", err)
	}
	ch := make(chan int)
	go coalesce(watcher.Events, ch)

	go func() {
		for {
			select {
			case count, ok := <-ch:
				if !ok {
					return
				}
				slog.Info("configuration file changed", "events", count)

				// Reload the configuration
				err = sql_exporter.Reload(e, &configFile)
				if err != nil {
					slog.Error("sql_exporter.reload", "error", err)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				slog.Error(err.Error())
			}
		}
	}()

	err = watcher.Add(folder)
	if err != nil {
		return fmt.Errorf("watcher.add: %w", err)
	}
	slog.Info("watching", "folder", folder)
	return nil
}

// coalesce watches fsnotify events and returns when no new event has happened
// for two seconds or after five seconds
func coalesce(in <-chan fsnotify.Event, out chan<- int) {

	timer := time.NewTicker(1 * time.Second)
	var events int // count of events

	active := false
	first := time.Time{}
	last := time.Time{}

	for {
		select {
		case e := <-in:
			events++
			slog.Debug("watch-in", "name", e.Name, "op", e.Op.String(), "count", events)
			last = time.Now()
			if !active {
				first = time.Now()
			}
			active = true

		case <-timer.C:
			if active {
				if time.Since(first) > time.Duration(5*time.Second) || time.Since(last) > time.Duration(2*time.Second) {
					slog.Debug("watch-out", "active", active, "first", first, "last", last)
					out <- events
					active = false
					events = 0
				}
			}
		}
	}
}
