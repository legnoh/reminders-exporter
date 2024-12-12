package collector

import (
	"time"

	"github.com/legnoh/reminders-exporter/pkg/reminder"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "reminder"
)

type ReminderCollector struct {
	Filters []reminder.Filter
}

var (
	reminderAllTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "all_total",
		Help:      "all reminder count",
	}, []string{"status"})

	reminderListTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "list_total",
		Help:      "target list uncompleted reminder count",
	}, []string{"name"})

	reminderPriorityTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "priority_total",
		Help:      "target priority uncompleted reminder count",
	}, []string{"priority"})

	reminderCustomfilterTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "customfilter_total",
		Help:      "reminder count of user defined custom filter",
	}, []string{"name"})
)

func (c ReminderCollector) Describe(ch chan<- *prometheus.Desc) {
	reminderAllTotal.Describe(ch)
	reminderListTotal.Describe(ch)
	reminderPriorityTotal.Describe(ch)
	reminderCustomfilterTotal.Describe(ch)
}

func (rc ReminderCollector) Collect(ch chan<- prometheus.Metric) {

	remdb := reminder.FetchReminderData()

	reminderAllTotal.WithLabelValues("completed").Set(remdb.Count(reminder.Condition{Completed: true}))
	reminderAllTotal.WithLabelValues("incompleted").Set(remdb.Count(reminder.Condition{Completed: false}))
	reminderAllTotal.WithLabelValues("outdated").Set(remdb.Count(reminder.Condition{Deadline: time.Duration(1)}))
	reminderAllTotal.Collect(ch)

	for _, listName := range remdb.Lists {
		reminderListTotal.WithLabelValues(listName).Set(remdb.Count(reminder.Condition{List: []string{listName}}))
	}
	reminderListTotal.Collect(ch)

	for _, priorityName := range reminder.Priorities {
		reminderPriorityTotal.WithLabelValues(priorityName).Set(remdb.Count(reminder.Condition{Priority: []string{priorityName}}))
	}
	reminderPriorityTotal.Collect(ch)

	for _, filter := range rc.Filters {
		reminderCustomfilterTotal.WithLabelValues(filter.Name).Set(remdb.Count(filter.Condition))
	}
	reminderCustomfilterTotal.Collect(ch)
}
