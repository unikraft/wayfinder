package stats

import (
  "fmt"
  "regexp"
  "strconv"

  libvirt "github.com/libvirt/libvirt-go"
)

var (
  // DomainCpu
)

func (s *DomainStats) LookupCpuCores() error {
  vcpus, err := s.domain.Domain().GetVcpus()
	if err != nil {
		return fmt.Errorf("could not get vCpus: %s", err)
	}

  cores := len(vcpus)
	cpuCores := s.domain.CreateMeasurement(uint64(cores))
	domain.AddMetricMeasurement("cpu_cores", cpuCores)

	// cache old thread Ids for cleanup
	var oldThreadIds []int
	oldThreadIds = append(oldThreadIds, domain.GetMetricIntArray("cpu_threadIds")...)
	oldThreadIds = append(oldThreadIds, domain.GetMetricIntArray("cpu_otherThreadIds")...)

	// get core thread Ids
	vCpuThreads, err := s.domain.Domain().QemuMonitorCommand("info cpus", libvirt.DOMAIN_QEMU_MONITOR_COMMAND_HMP)
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

	newMeasurementThreads := s.domain.CreateMeasurement(coreThreadIds)
	s.domain.AddMetricMeasurement("cpu_threadIds", newMeasurementThreads)

	// get thread ids
	tasksFolder := fmt.Sprint(s.domain.p.Cfg.ProcFS, "/", s.domain.Pid(), "/task/*")
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

	s.domain.AddMetricMeasurement("cpu_otherThreadIds", s.domain.CreateMeasurement(otherThreadIds))

	// remove cached but not existent thread Ids
	for _, id := range oldThreadIds {
		s.domain.DeleteMetric(fmt.Sprint("cpu_times_", id))
		s.domain.DeleteMetric(fmt.Sprint("cpu_runqueues_", id))
		s.domain.DeleteMetric(fmt.Sprint("cpu_other_times_", id))
		s.domain.DeleteMetric(fmt.Sprint("cpu_other_runqueues_", id))
	}

  return nil
}

func (s *DomainStats) MeasureCpuCores() {
  cpuCollectMeasurements("cpu_threadIds", "cpu_")
  cpuCollectMeasurements("cpu_otherThreadIds", "cpu_other_")
}

func (s *DomainStats) cpuCollectMeasurements(name string, prefix string) {
  threadIds := s.domain.GetMetricIntArray(name)
  for _, threadId := range threadIds {
    schedstat := proc.GetProcPIdSchedStat(threadId)
    s.domain.AddMetricMeasurement(fmt.Sprint(prefix, "times_", threadId), models.CreateMeasurement(schedstat.Cputime))
    s.domain.AddMetricMeasurement(fmt.Sprint(prefix, "runqueues_", threadId), models.CreateMeasurement(schedstat.Runqueue))
  }
}
