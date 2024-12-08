package main

import (
	"fmt"
	"time"

	"github.com/legnoh/reminder-exporter/internal"
)

// var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

func main() {
	// flag.Parse()

	// reg := prometheus.NewRegistry()

	// reg.MustRegister(
	// 	collectors.NewGoCollector(),
	// 	collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	// )

	// http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	// log.Fatal(http.ListenAndServe(*addr, nil))
	reminders := internal.GetAllReminder()
	condition := internal.Reminder{
		DueDate: time.Now(),
	}
	outdated := internal.SearchReminder(reminders, condition)
	fmt.Printf(outdated[0].Title)
}
