package libvirt

import (
  "fmt"
  "regexp"
  "strconv"
  "path"
  "path/filepath"

  libvirt "github.com/libvirt/libvirt-go"

  "github.com/unikraft/wayfinder/pkg/proc"
  "github.com/unikraft/wayfinder/internal/metrics"
)

func (d *Domain) CpuLookup() error {
  vcpus, err := d.domain.GetVcpus()
  if err != nil {
    return fmt.Errorf("could not get vCpus: %s", err)
  }

  cores := len(vcpus)
  cpuCores := metrics.CreateMeasurement(uint64(cores))

  d.AddMeasurement("cpu_cores", cpuCores)

  // cache old thread Ids for cleanup
  var oldThreadIds []int
  oldThreadIds = append(oldThreadIds, d.GetMetricIntArray("cpu_threadIds")...)
  oldThreadIds = append(oldThreadIds, d.GetMetricIntArray("cpu_otherThreadIds")...)

  // get core thread Ids
  vCpuThreads, err := d.domain.QemuMonitorCommand("info cpus", libvirt.DOMAIN_QEMU_MONITOR_COMMAND_HMP)
  if err != nil {
    return fmt.Errorf("could not get core thread ids: %s", err)
  }
  regThreadId := regexp.MustCompile("thread_id=([0-9]*)\\s")
  threadIdsRaw := regThreadId.FindAllStringSubmatch(vCpuThreads, -1)
  coreThreadIds := make([]int, len(threadIdsRaw))

  for i, thread := range threadIdsRaw {
    threadId, _ := strconv.Atoi(thread[1])
    coreThreadIds[i] = threadId
    oldThreadIds = removeFromArray(oldThreadIds, threadId)
  }

  d.AddMeasurement("cpu_threadIds", metrics.CreateMeasurement(coreThreadIds))

  // get thread ids
  pid, err := d.Pid()
  if err != nil {
    return fmt.Errorf("domain has no pid: %s", err)
  }

  tasksFolder := fmt.Sprint(d.p.Cfg.ProcFS, "/", pid, "/task/*")
  files, err := filepath.Glob(tasksFolder)
  if err != nil {
    return fmt.Errorf("could not get procfs tasks: %s", err)
  }

  otherThreadIds := make([]int, 0)
  i := 0
  for _, f := range files {
    taskId, _ := strconv.Atoi(path.Base(f))
    found := false
    for _, n := range coreThreadIds {
      if taskId == n {
        // taskId is for vCPU core. skip.
        found = true
        break
      }
    }
    if found {
      // taskId is for vCPU core. skip.
      continue
    }

    // taskId is not for a vCPU core
    otherThreadIds = append(otherThreadIds, taskId)
    oldThreadIds = removeFromArray(oldThreadIds, taskId)
    i++
  }

  d.AddMeasurement("cpu_otherThreadIds", metrics.CreateMeasurement(otherThreadIds))

  // remove cached but not existent thread Ids
  for _, id := range oldThreadIds {
    d.DeleteMetric(fmt.Sprint("cpu_times_", id))
    d.DeleteMetric(fmt.Sprint("cpu_runqueues_", id))
    d.DeleteMetric(fmt.Sprint("cpu_other_times_", id))
    d.DeleteMetric(fmt.Sprint("cpu_other_runqueues_", id))
  }

  return nil
}

func (d *Domain) CpuMeasure() error {
  d.cpuCollectMeasurements("cpu_threadIds", "cpu_")
  d.cpuCollectMeasurements("cpu_otherThreadIds", "cpu_other_")

  return nil
}

func (d *Domain) cpuCollectMeasurements(name string, prefix string) {
  threadIds := d.GetMetricIntArray(name)
  for _, threadId := range threadIds {
    schedstat := proc.GetProcPIDSchedStat(d.p.Cfg.ProcFS, threadId)
    d.AddMeasurement(fmt.Sprint(prefix, "times_", threadId), metrics.CreateMeasurement(schedstat.Cputime))
    d.AddMeasurement(fmt.Sprint(prefix, "runqueues_", threadId), metrics.CreateMeasurement(schedstat.Runqueue))
  }
}

func (d *Domain) CpuPrint() map[string]string {
  cores, _ := d.GetMetricUint64("cpu_cores", 0)

  return map[string]string{
    "cpu_cores": cores,

    // cpu util for vcores
    "cpu_total": d.CpuPrintThreadMetric("cpu_threadIds", "cpu_times"),
    "cpu_steal": d.CpuPrintThreadMetric("cpu_threadIds", "cpu_runqueues"),

    // cpu util for for other threads (i/o or emulation)
    "cpu_other_total": d.CpuPrintThreadMetric("cpu_otherThreadIDs", "cpu_other_times"),
    "cpu_other_steal": d.CpuPrintThreadMetric("cpu_otherThreadIDs", "cpu_other_runqueues"),
  }
}

func (d *Domain) CpuPrintThreadMetric(lookupMetric string, metric string) string {
  threadIds := d.GetMetricIntArray(lookupMetric)
  var measurementSum float64
  var measurementCount int

  for _, threadId := range threadIds {
    metricName := fmt.Sprint(metric, "_", threadId)
    measurementStr := d.GetMetricDiffUint64(metricName, true)
    if measurementStr == "" {
      continue
    }
    measurement, err := strconv.ParseUint(measurementStr, 10, 64)
    if err != nil {
      fmt.Printf("copuld not convert error!\n")
      continue
    }
    measurementSeconds := float64(measurement) / 1000000000 // since counters are nanoseconds
    measurementSum += measurementSeconds
    measurementCount++
  }

  var avg float64
  if measurementCount > 0 {
    avg = float64(measurementSum) / float64(measurementCount)
  }

  percent := avg * 100
  return fmt.Sprintf("%.0f", percent)
}


func removeFromArray(s []int, r int) []int {
  for i, v := range s {
    if v == r {
      return append(s[:i], s[i+1:]...)
    }
  }

  return s
}
