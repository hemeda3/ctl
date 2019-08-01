package cmd

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"github.com/wish/ctl/pkg/client/types"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

// REVIEW: Most of the processing here was guessed with reverse engineering
// by comparing with the output of kubectl
func printPodList(lst []types.PodDiscovery, labelColumns []string) {
	if len(lst) == 0 {
		fmt.Println("No pods found")
		return
	}
	// Insert default columns
	defaultColumns := viper.GetStringSlice("default_columns")
	var newLabelColumns []string
	if len(defaultColumns) == 0 {
		newLabelColumns = labelColumns
	} else if len(labelColumns) == 0 {
		newLabelColumns = defaultColumns
	} else {
		for _, s := range defaultColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
		for _, s := range labelColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprint(w, "CONTEXT\tNAMESPACE\tNAME\tREADY\tSTATUS\tRESTARTS\tAGE")
	for _, l := range newLabelColumns {
		fmt.Fprint(w, "\t", strings.ToUpper(l))
	}
	fmt.Fprintln(w)

	for _, v := range lst {
		fmt.Fprintf(w, "%s", v.Context)
		fmt.Fprintf(w, "\t%s", v.Namespace)
		fmt.Fprintf(w, "\t%s", v.Name)
		var ready int
		for _, s := range v.Status.ContainerStatuses {
			if s.Ready {
				ready++
			}
		}
		fmt.Fprintf(w, "\t%d/%d", ready, len(v.Spec.Containers))
		fmt.Fprintf(w, "\t%s", v.Status.Phase) // A bit off from kubectl output
		// Restarts
		var restarts int32
		for _, s := range v.Status.ContainerStatuses {
			restarts += s.RestartCount
		}
		fmt.Fprintf(w, "\t%d", restarts)
		fmt.Fprintf(w, "\t%v", time.Since(v.CreationTimestamp.Time).Round(time.Second))

		for _, l := range newLabelColumns {
			fmt.Fprint(w, "\t")
			if _, ok := v.Labels[l]; ok {
				fmt.Fprint(w, v.Labels[l])
			}
		}
		fmt.Fprintln(w)
	}
	w.Flush()
}

func describePod(pod types.PodDiscovery) {
	fmt.Printf("context: %s\n", pod.Context)
	b, _ := yaml.Marshal(pod.Pod)
	fmt.Println(string(b))
}

func describePodList(lst []types.PodDiscovery) {
	for _, pod := range lst {
		describePod(pod)
	}
}

func printConfigMapList(lst []types.ConfigMapDiscovery, labelColumns []string) {
	if len(lst) == 0 {
		fmt.Println("No config maps found")
		return
	}
	// Insert default columns
	defaultColumns := viper.GetStringSlice("default_columns")
	var newLabelColumns []string
	if len(defaultColumns) == 0 {
		newLabelColumns = labelColumns
	} else if len(labelColumns) == 0 {
		newLabelColumns = defaultColumns
	} else {
		for _, s := range defaultColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
		for _, s := range labelColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprint(w, "CONTEXT\tNAMESPACE\tNAME\tDATA\tAGE")
	for _, l := range newLabelColumns {
		fmt.Fprint(w, "\t", strings.ToUpper(l))
	}
	fmt.Fprintln(w)

	for _, v := range lst {
		fmt.Fprintf(w, "%s", v.Context)
		fmt.Fprintf(w, "\t%s", v.Namespace)
		fmt.Fprintf(w, "\t%s", v.Name)
		fmt.Fprintf(w, "\t%d", len(v.Data))
		// Age
		fmt.Fprintf(w, "\t%v", time.Since(v.CreationTimestamp.Time).Round(time.Second))

		for _, l := range newLabelColumns {
			fmt.Fprint(w, "\t")
			if _, ok := v.Labels[l]; ok {
				fmt.Fprint(w, v.Labels[l])
			}
		}
		fmt.Fprintln(w)
	}
	w.Flush()
}

func describeConfigMap(cm types.ConfigMapDiscovery) {
	fmt.Printf("context: %s\n", cm.Context)
	b, _ := yaml.Marshal(cm.ConfigMap)
	fmt.Println(string(b))
}

func describeConfigMapList(lst []types.ConfigMapDiscovery) {
	for _, cm := range lst {
		describeConfigMap(cm)
	}
}

func printDeploymentList(lst []types.DeploymentDiscovery, labelColumns []string) {
	if len(lst) == 0 {
		fmt.Println("No deployments found")
		return
	}
	// Insert default columns
	defaultColumns := viper.GetStringSlice("default_columns")
	var newLabelColumns []string
	if len(defaultColumns) == 0 {
		newLabelColumns = labelColumns
	} else if len(labelColumns) == 0 {
		newLabelColumns = defaultColumns
	} else {
		for _, s := range defaultColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
		for _, s := range labelColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprint(w, "CONTEXT\tNAMESPACE\tNAME\tDESIRED\tCURRENT\tUP-TO-DATE\tAVAILABLE\tAGE")
	for _, l := range newLabelColumns {
		fmt.Fprint(w, "\t", strings.ToUpper(l))
	}
	fmt.Fprintln(w)

	for _, v := range lst {
		fmt.Fprintf(w, "%s", v.Context)
		fmt.Fprintf(w, "\t%s", v.Namespace)
		fmt.Fprintf(w, "\t%s", v.Name)
		// Desired
		fmt.Fprintf(w, "\t%d", *v.Spec.Replicas)
		// Current
		fmt.Fprintf(w, "\t%d", v.Status.Replicas)
		// Up-to-date
		fmt.Fprintf(w, "\t%d", v.Status.UpdatedReplicas)
		// Available
		fmt.Fprintf(w, "\t%d", v.Status.AvailableReplicas)
		// Age
		fmt.Fprintf(w, "\t%v", time.Since(v.CreationTimestamp.Time).Round(time.Second))

		for _, l := range newLabelColumns {
			fmt.Fprint(w, "\t")
			if _, ok := v.Labels[l]; ok {
				fmt.Fprint(w, v.Labels[l])
			}
		}
		fmt.Fprintln(w)
	}
	w.Flush()
}

func describeDeployment(d types.DeploymentDiscovery) {
	fmt.Printf("context: %s\n", d.Context)
	b, _ := yaml.Marshal(d.Deployment)
	fmt.Println(string(b))
}

func describeDeploymentList(lst []types.DeploymentDiscovery) {
	for _, d := range lst {
		describeDeployment(d)
	}
}

func printReplicaSetList(lst []types.ReplicaSetDiscovery, labelColumns []string) {
	if len(lst) == 0 {
		fmt.Println("No replicasets found")
		return
	}
	// Insert default columns
	defaultColumns := viper.GetStringSlice("default_columns")
	var newLabelColumns []string
	if len(defaultColumns) == 0 {
		newLabelColumns = labelColumns
	} else if len(labelColumns) == 0 {
		newLabelColumns = defaultColumns
	} else {
		for _, s := range defaultColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
		for _, s := range labelColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprint(w, "CONTEXT\tNAMESPACE\tNAME\tDESIRED\tCURRENT\tREADY\tAGE")
	for _, l := range newLabelColumns {
		fmt.Fprint(w, "\t", strings.ToUpper(l))
	}
	fmt.Fprintln(w)

	for _, v := range lst {
		fmt.Fprintf(w, "%s", v.Context)
		fmt.Fprintf(w, "\t%s", v.Namespace)
		fmt.Fprintf(w, "\t%s", v.Name)
		// Desired
		fmt.Fprintf(w, "\t%d", *v.Spec.Replicas)
		// Current
		fmt.Fprintf(w, "\t%d", v.Status.Replicas)
		// Ready
		fmt.Fprintf(w, "\t%d", v.Status.ReadyReplicas)
		// Age
		fmt.Fprintf(w, "\t%v", time.Since(v.CreationTimestamp.Time).Round(time.Second))

		for _, l := range newLabelColumns {
			fmt.Fprint(w, "\t")
			if _, ok := v.Labels[l]; ok {
				fmt.Fprint(w, v.Labels[l])
			}
		}
		fmt.Fprintln(w)
	}
	w.Flush()
}

func describeReplicaSet(d types.ReplicaSetDiscovery) {
	fmt.Printf("context: %s\n", d.Context)
	b, _ := yaml.Marshal(d.ReplicaSet)
	fmt.Println(string(b))
}

func describeReplicaSetList(lst []types.ReplicaSetDiscovery) {
	for _, d := range lst {
		describeReplicaSet(d)
	}
}

func printJobList(lst []types.JobDiscovery, labelColumns []string) {
	if len(lst) == 0 {
		fmt.Println("No replicasets found")
		return
	}
	// Insert default columns
	defaultColumns := viper.GetStringSlice("default_columns")
	var newLabelColumns []string
	if len(defaultColumns) == 0 {
		newLabelColumns = labelColumns
	} else if len(labelColumns) == 0 {
		newLabelColumns = defaultColumns
	} else {
		for _, s := range defaultColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
		for _, s := range labelColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprint(w, "CONTEXT\tNAMESPACE\tNAME\tSTATE\tSTART\tEND")
	for _, l := range newLabelColumns {
		fmt.Fprint(w, "\t", strings.ToUpper(l))
	}
	fmt.Fprintln(w)

	for _, v := range lst {
		fmt.Fprintf(w, "%s", v.Context)
		fmt.Fprintf(w, "\t%s", v.Namespace)
		fmt.Fprintf(w, "\t%s", v.Name)
		// State
		if v.Status.Failed > 0 {
			fmt.Fprint(w, "\tFAILED")
		} else if v.Status.CompletionTime != nil {
			fmt.Fprint(w, "\tSUCCESSFUL")
		} else {
			fmt.Fprint(w, "\tIN PROGRESS")
		}
		// Start
		fmt.Fprintf(w, "\t%v", v.Status.StartTime)
		// END
		if v.Status.CompletionTime != nil {
			fmt.Fprintf(w, "\t%v", v.Status.CompletionTime)
		} else {
			fmt.Fprint(w, "\t<none>")
		}

		for _, l := range newLabelColumns {
			fmt.Fprint(w, "\t")
			if _, ok := v.Labels[l]; ok {
				fmt.Fprint(w, v.Labels[l])
			}
		}
		fmt.Fprintln(w)
	}
	w.Flush()
}

func describeJob(d types.JobDiscovery) {
	fmt.Printf("context: %s\n", d.Context)
	b, _ := yaml.Marshal(d.Job)
	fmt.Println(string(b))
}

func describeJobList(lst []types.JobDiscovery) {
	for _, d := range lst {
		describeJob(d)
	}
}

func printCronJobList(lst []types.CronJobDiscovery, labelColumns []string) {
	if len(lst) == 0 {
		fmt.Println("No cronjobs found")
		return
	}
	// Insert default columns
	defaultColumns := viper.GetStringSlice("default_columns")
	var newLabelColumns []string
	if len(defaultColumns) == 0 {
		newLabelColumns = labelColumns
	} else if len(labelColumns) == 0 {
		newLabelColumns = defaultColumns
	} else {
		for _, s := range defaultColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
		for _, s := range labelColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprint(w, "CONTEXT\tNAME\tSCHEDULE\tSUSPEND\tACTIVE\tLAST SCHEDULE\tNEXT RUN\tAGE")
	for _, l := range newLabelColumns {
		fmt.Fprint(w, "\t", strings.ToUpper(l))
	}
	fmt.Fprintln(w)

	for _, v := range lst {
		fmt.Fprintf(w, "%s", v.Context)
		fmt.Fprintf(w, "\t%s", v.Name)
		fmt.Fprintf(w, "\t%s", v.Spec.Schedule) // Schedule
		fmt.Fprintf(w, "\t%t", *v.Spec.Suspend) // Suspend
		fmt.Fprintf(w, "\t%d", len(v.Status.Active))
		// Last schedule
		// TODO fix rounding
		if v.Status.LastScheduleTime == nil {
			fmt.Fprintf(w, "\t<none>")
		} else {
			fmt.Fprintf(w, "\t%v", time.Since(v.Status.LastScheduleTime.Time).Round(time.Second))
		}
		// Next run
		s, _ := cron.ParseStandard(v.Spec.Schedule)
		fmt.Fprintf(w, "\t%v", time.Until(s.Next(time.Now())).Round(time.Second))
		// Age
		fmt.Fprintf(w, "\t%v", time.Since(v.CreationTimestamp.Time).Round(time.Second))
		// Labels
		for _, l := range newLabelColumns {
			fmt.Fprint(w, "\t")
			if _, ok := v.Labels[l]; ok {
				fmt.Fprint(w, v.Labels[l])
			}
		}
		fmt.Fprintln(w)
	}
	w.Flush()
}

func describeCronJob(c types.CronJobDiscovery) {
	fmt.Printf("context: %s\n", c.Context)
	b, _ := yaml.Marshal(c.CronJob)
	fmt.Println(string(b))
}

func describeCronJobList(lst []types.CronJobDiscovery) {
	for _, cj := range lst {
		describeCronJob(cj)
	}
}
