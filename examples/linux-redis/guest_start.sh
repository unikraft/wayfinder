#!/bin/sh

set -ex

# Usage: guest_start.sh $IP_ADDR $IP_GW

echo "Init called as: $@"

export PATH="/usr/local/bin:$PATH"

mount -t proc proc /proc
ulimit -n 65535
echo 1024 > /proc/sys/net/core/somaxconn

# ip self, ip gateway
./busybox-x86_64 ip addr add $WAYFINDER_DOMAIN_IP_ADDR/24 dev eth0
./busybox-x86_64 ip addr add 127.0.0.1/24 dev lo
./busybox-x86_64 ip link set eth0 up
./busybox-x86_64 ip link set lo up
./busybox-x86_64 ip route add default via $WAYFINDER_DOMAIN_IP_GW_ADDR dev eth0

# Setting all runtime parameters: beware, some might fail

# set everything inside `proc`. Display error for what failed
if [ -n "$AIO_MAX_NR" ] && [ -w /proc/sys/fs/aio-max-nr ]; then
  echo $AIO_MAX_NR > /proc/sys/fs/aio-max-nr &
fi

if [ -n "$DIR_NOTIFY_ENABLE" ] && [ -w /proc/sys/fs/dir-notify-enable ]; then
  echo $DIR_NOTIFY_ENABLE > /proc/sys/fs/dir-notify-enable &
fi

if [ -n "$FILE_MAX" ] && [ -w /proc/sys/fs/file-max ]; then
  echo $FILE_MAX > /proc/sys/fs/file-max &
fi

if [ -n "$LEASE_BREAK_TIME" ] && [ -w /proc/sys/fs/lease-break-time ]; then
  echo $LEASE_BREAK_TIME > /proc/sys/fs/lease-break-time &
fi

if [ -n "$LEASES_ENABLE" ] && [ -w /proc/sys/fs/leases-enable ]; then
  echo $LEASES_ENABLE > /proc/sys/fs/leases-enable &
fi

if [ -n "$NR_OPEN" ] && [ -w /proc/sys/fs/nr_open ]; then
  echo $NR_OPEN > /proc/sys/fs/nr_open &
fi

if [ -n "$OVERFLOWGID" ] && [ -w /proc/sys/fs/overflowgid ]; then
  echo $OVERFLOWGID > /proc/sys/fs/overflowgid &
fi

if [ -n "$OVERFLOWUID" ] && [ -w /proc/sys/fs/overflowuid ]; then
  echo $OVERFLOWUID > /proc/sys/fs/overflowuid &
fi

if [ -n "$PIPE_MAX_SIZE" ] && [ -w /proc/sys/fs/pipe-max-size ]; then
  echo $PIPE_MAX_SIZE > /proc/sys/fs/pipe-max-size &
fi

if [ -n "$PROTECTED_HARDLINKS" ] && [ -w /proc/sys/fs/protected_hardlinks ]; then
  echo $PROTECTED_HARDLINKS > /proc/sys/fs/protected_hardlinks &
fi

if [ -n "$PROTECTED_SYMLINKS" ] && [ -w /proc/sys/fs/protected_symlinks ]; then
  echo $PROTECTED_SYMLINKS > /proc/sys/fs/protected_symlinks &
fi

if [ -n "$SUID_DUMPABLE" ] && [ -w /proc/sys/fs/suid_dumpable ]; then
  echo $SUID_DUMPABLE > /proc/sys/fs/suid_dumpable &
fi

if [ -n "$ACPI_VIDEO_FLAGS" ] && [ -w /proc/sys/kernel/acpi_video_flags ]; then
  echo $ACPI_VIDEO_FLAGS > /proc/sys/kernel/acpi_video_flags &
fi

if [ -n "$AUTO_MSGMNI" ] && [ -w /proc/sys/kernel/auto_msgmni ]; then
  echo $AUTO_MSGMNI > /proc/sys/kernel/auto_msgmni &
fi

if [ -n "$CAD_PID" ] && [ -w /proc/sys/kernel/cad_pid ]; then
  echo $CAD_PID > /proc/sys/kernel/cad_pid &
fi

if [ -n "$COMPAT_LOG" ] && [ -w /proc/sys/kernel/compat-log ]; then
  echo $COMPAT_LOG > /proc/sys/kernel/compat-log &
fi

if [ -n "$CORE_PIPE_LIMIT" ] && [ -w /proc/sys/kernel/core_pipe_limit ]; then
  echo $CORE_PIPE_LIMIT > /proc/sys/kernel/core_pipe_limit &
fi

if [ -n "$CORE_USES_PID" ] && [ -w /proc/sys/kernel/core_uses_pid ]; then
  echo $CORE_USES_PID > /proc/sys/kernel/core_uses_pid &
fi

if [ -n "$CTRL_ALT_DEL" ] && [ -w /proc/sys/kernel/ctrl-alt-del ]; then
  echo $CTRL_ALT_DEL > /proc/sys/kernel/ctrl-alt-del &
fi

if [ -n "$DMESG_RESTRICT" ] && [ -w /proc/sys/kernel/dmesg_restrict ]; then
  echo $DMESG_RESTRICT > /proc/sys/kernel/dmesg_restrict &
fi

if [ -n "$FTRACE_DUMP_ON_OOPS" ] && [ -w /proc/sys/kernel/ftrace_dump_on_oops ]; then
  echo $FTRACE_DUMP_ON_OOPS > /proc/sys/kernel/ftrace_dump_on_oops &
fi

if [ -n "$FTRACE_ENABLED" ] && [ -w /proc/sys/kernel/ftrace_enabled ]; then
  echo $FTRACE_ENABLED > /proc/sys/kernel/ftrace_enabled &
fi

if [ -n "$HUNG_TASK_CHECK_COUNT" ] && [ -w /proc/sys/kernel/hung_task_check_count ]; then
  echo $HUNG_TASK_CHECK_COUNT > /proc/sys/kernel/hung_task_check_count &
fi

if [ -n "$HUNG_TASK_PANIC" ] && [ -w /proc/sys/kernel/hung_task_panic ]; then
  echo $HUNG_TASK_PANIC > /proc/sys/kernel/hung_task_panic &
fi

if [ -n "$HUNG_TASK_TIMEOUT_SECS" ] && [ -w /proc/sys/kernel/hung_task_timeout_secs ]; then
  echo $HUNG_TASK_TIMEOUT_SECS > /proc/sys/kernel/hung_task_timeout_secs &
fi

if [ -n "$HUNG_TASK_WARNINGS" ] && [ -w /proc/sys/kernel/hung_task_warnings ]; then
  echo $HUNG_TASK_WARNINGS > /proc/sys/kernel/hung_task_warnings &
fi

if [ -n "$IO_DELAY_TYPE" ] && [ -w /proc/sys/kernel/io_delay_type ]; then
  echo $IO_DELAY_TYPE > /proc/sys/kernel/io_delay_type &
fi

if [ -n "$KEXEC_LOAD_DISABLED" ] && [ -w /proc/sys/kernel/kexec_load_disabled ]; then
  echo $KEXEC_LOAD_DISABLED > /proc/sys/kernel/kexec_load_disabled &
fi

if [ -n "$KPTR_RESTRICT" ] && [ -w /proc/sys/kernel/kptr_restrict ]; then
  echo $KPTR_RESTRICT > /proc/sys/kernel/kptr_restrict &
fi

if [ -n "$KSTACK_DEPTH_TO_PRINT" ] && [ -w /proc/sys/kernel/kstack_depth_to_print ]; then
  echo $KSTACK_DEPTH_TO_PRINT > /proc/sys/kernel/kstack_depth_to_print &
fi

if [ -n "$LATENCYTOP" ] && [ -w /proc/sys/kernel/latencytop ]; then
  echo $LATENCYTOP > /proc/sys/kernel/latencytop &
fi

if [ -n "$MAX_LOCK_DEPTH" ] && [ -w /proc/sys/kernel/max_lock_depth ]; then
  echo $MAX_LOCK_DEPTH > /proc/sys/kernel/max_lock_depth &
fi

if [ -n "$MSGMAX" ] && [ -w /proc/sys/kernel/msgmax ]; then
  echo $MSGMAX > /proc/sys/kernel/msgmax &
fi

if [ -n "$MSGMNB" ] && [ -w /proc/sys/kernel/msgmnb ]; then
  echo $MSGMNB > /proc/sys/kernel/msgmnb &
fi

if [ -n "$MSGMNI" ] && [ -w /proc/sys/kernel/msgmni ]; then
  echo $MSGMNI > /proc/sys/kernel/msgmni &
fi

if [ -n "$NMI_WATCHDOG" ] && [ -w /proc/sys/kernel/nmi_watchdog ]; then
  echo $NMI_WATCHDOG > /proc/sys/kernel/nmi_watchdog &
fi

if [ -n "$NS_LAST_PID" ] && [ -w /proc/sys/kernel/ns_last_pid ]; then
  echo $NS_LAST_PID > /proc/sys/kernel/ns_last_pid &
fi

if [ -n "$NUMA_BALANCING" ] && [ -w /proc/sys/kernel/numa_balancing ]; then
  echo $NUMA_BALANCING > /proc/sys/kernel/numa_balancing &
fi

if [ -n "$NUMA_BALANCING_SCAN_DELAY_MS" ] && [ -w /proc/sys/kernel/numa_balancing_scan_delay_ms ]; then
  echo $NUMA_BALANCING_SCAN_DELAY_MS > /proc/sys/kernel/numa_balancing_scan_delay_ms &
fi

if [ -n "$NUMA_BALANCING_SCAN_PERIOD_MAX_MS" ] && [ -w /proc/sys/kernel/numa_balancing_scan_period_max_ms ]; then
  echo $NUMA_BALANCING_SCAN_PERIOD_MAX_MS > /proc/sys/kernel/numa_balancing_scan_period_max_ms &
fi

if [ -n "$NUMA_BALANCING_SCAN_PERIOD_MIN_MS" ] && [ -w /proc/sys/kernel/numa_balancing_scan_period_min_ms ]; then
  echo $NUMA_BALANCING_SCAN_PERIOD_MIN_MS > /proc/sys/kernel/numa_balancing_scan_period_min_ms &
fi

if [ -n "$NUMA_BALANCING_SCAN_SIZE_MB" ] && [ -w /proc/sys/kernel/numa_balancing_scan_size_mb ]; then
  echo $NUMA_BALANCING_SCAN_SIZE_MB > /proc/sys/kernel/numa_balancing_scan_size_mb &
fi

if [ -n "$OVERFLOWGID" ] && [ -w /proc/sys/kernel/overflowgid ]; then
  echo $OVERFLOWGID > /proc/sys/kernel/overflowgid &
fi

if [ -n "$OVERFLOWUID" ] && [ -w /proc/sys/kernel/overflowuid ]; then
  echo $OVERFLOWUID > /proc/sys/kernel/overflowuid &
fi

if [ -n "$PANIC" ] && [ -w /proc/sys/kernel/panic ]; then
  echo $PANIC > /proc/sys/kernel/panic &
fi

if [ -n "$PANIC_ON_IO_NMI" ] && [ -w /proc/sys/kernel/panic_on_io_nmi ]; then
  echo $PANIC_ON_IO_NMI > /proc/sys/kernel/panic_on_io_nmi &
fi

if [ -n "$PANIC_ON_OOPS" ] && [ -w /proc/sys/kernel/panic_on_oops ]; then
  echo $PANIC_ON_OOPS > /proc/sys/kernel/panic_on_oops &
fi

if [ -n "$PANIC_ON_UNRECOVERED_NMI" ] && [ -w /proc/sys/kernel/panic_on_unrecovered_nmi ]; then
  echo $PANIC_ON_UNRECOVERED_NMI > /proc/sys/kernel/panic_on_unrecovered_nmi &
fi

if [ -n "$PANIC_ON_WARN" ] && [ -w /proc/sys/kernel/panic_on_warn ]; then
  echo $PANIC_ON_WARN > /proc/sys/kernel/panic_on_warn &
fi

if [ -n "$PERF_CPU_TIME_MAX_PERCENT" ] && [ -w /proc/sys/kernel/perf_cpu_time_max_percent ]; then
  echo $PERF_CPU_TIME_MAX_PERCENT > /proc/sys/kernel/perf_cpu_time_max_percent &
fi

if [ -n "$PERF_EVENT_MAX_SAMPLE_RATE" ] && [ -w /proc/sys/kernel/perf_event_max_sample_rate ]; then
  echo $PERF_EVENT_MAX_SAMPLE_RATE > /proc/sys/kernel/perf_event_max_sample_rate &
fi

if [ -n "$PERF_EVENT_MLOCK_KB" ] && [ -w /proc/sys/kernel/perf_event_mlock_kb ]; then
  echo $PERF_EVENT_MLOCK_KB > /proc/sys/kernel/perf_event_mlock_kb &
fi

if [ -n "$PERF_EVENT_PARANOID" ] && [ -w /proc/sys/kernel/perf_event_paranoid ]; then
  echo $PERF_EVENT_PARANOID > /proc/sys/kernel/perf_event_paranoid &
fi

if [ -n "$PID_MAX" ] && [ -w /proc/sys/kernel/pid_max ]; then
  echo $PID_MAX > /proc/sys/kernel/pid_max &
fi

if [ -n "$PRINT_FATAL_SIGNALS" ] && [ -w /proc/sys/kernel/print-fatal-signals ]; then
  echo $PRINT_FATAL_SIGNALS > /proc/sys/kernel/print-fatal-signals &
fi

if [ -n "$PRINTK" ] && [ -w /proc/sys/kernel/printk ]; then
  echo $PRINTK > /proc/sys/kernel/printk &
fi

if [ -n "$PRINTK_DELAY" ] && [ -w /proc/sys/kernel/printk_delay ]; then
  echo $PRINTK_DELAY > /proc/sys/kernel/printk_delay &
fi

if [ -n "$PRINTK_RATELIMIT" ] && [ -w /proc/sys/kernel/printk_ratelimit ]; then
  echo $PRINTK_RATELIMIT > /proc/sys/kernel/printk_ratelimit &
fi

if [ -n "$PRINTK_RATELIMIT_BURST" ] && [ -w /proc/sys/kernel/printk_ratelimit_burst ]; then
  echo $PRINTK_RATELIMIT_BURST > /proc/sys/kernel/printk_ratelimit_burst &
fi

if [ -n "$RANDOMIZE_VA_SPACE" ] && [ -w /proc/sys/kernel/randomize_va_space ]; then
  echo $RANDOMIZE_VA_SPACE > /proc/sys/kernel/randomize_va_space &
fi

if [ -n "$REAL_ROOT_DEV" ] && [ -w /proc/sys/kernel/real-root-dev ]; then
  echo $REAL_ROOT_DEV > /proc/sys/kernel/real-root-dev &
fi

if [ -n "$SCHED_AUTOGROUP_ENABLED" ] && [ -w /proc/sys/kernel/sched_autogroup_enabled ]; then
  echo $SCHED_AUTOGROUP_ENABLED > /proc/sys/kernel/sched_autogroup_enabled &
fi

if [ -n "$SCHED_CFS_BANDWIDTH_SLICE_US" ] && [ -w /proc/sys/kernel/sched_cfs_bandwidth_slice_us ]; then
  echo $SCHED_CFS_BANDWIDTH_SLICE_US > /proc/sys/kernel/sched_cfs_bandwidth_slice_us &
fi

if [ -n "$SCHED_CHILD_RUNS_FIRST" ] && [ -w /proc/sys/kernel/sched_child_runs_first ]; then
  echo $SCHED_CHILD_RUNS_FIRST > /proc/sys/kernel/sched_child_runs_first &
fi

if [ -n "$SCHED_LATENCY_NS" ] && [ -w /proc/sys/kernel/sched_latency_ns ]; then
  echo $SCHED_LATENCY_NS > /proc/sys/kernel/sched_latency_ns &
fi

if [ -n "$SCHED_MIGRATION_COST_NS" ] && [ -w /proc/sys/kernel/sched_migration_cost_ns ]; then
  echo $SCHED_MIGRATION_COST_NS > /proc/sys/kernel/sched_migration_cost_ns &
fi

if [ -n "$SCHED_MIN_GRANULARITY_NS" ] && [ -w /proc/sys/kernel/sched_min_granularity_ns ]; then
  echo $SCHED_MIN_GRANULARITY_NS > /proc/sys/kernel/sched_min_granularity_ns &
fi

if [ -n "$SCHED_NR_MIGRATE" ] && [ -w /proc/sys/kernel/sched_nr_migrate ]; then
  echo $SCHED_NR_MIGRATE > /proc/sys/kernel/sched_nr_migrate &
fi

if [ -n "$SCHED_RR_TIMESLICE_MS" ] && [ -w /proc/sys/kernel/sched_rr_timeslice_ms ]; then
  echo $SCHED_RR_TIMESLICE_MS > /proc/sys/kernel/sched_rr_timeslice_ms &
fi

if [ -n "$SCHED_RT_PERIOD_US" ] && [ -w /proc/sys/kernel/sched_rt_period_us ]; then
  echo $SCHED_RT_PERIOD_US > /proc/sys/kernel/sched_rt_period_us &
fi

if [ -n "$SCHED_RT_RUNTIME_US" ] && [ -w /proc/sys/kernel/sched_rt_runtime_us ]; then
  echo $SCHED_RT_RUNTIME_US > /proc/sys/kernel/sched_rt_runtime_us &
fi

if [ -n "$SCHED_SHARES_WINDOW_NS" ] && [ -w /proc/sys/kernel/sched_shares_window_ns ]; then
  echo $SCHED_SHARES_WINDOW_NS > /proc/sys/kernel/sched_shares_window_ns &
fi

if [ -n "$SCHED_TIME_AVG_MS" ] && [ -w /proc/sys/kernel/sched_time_avg_ms ]; then
  echo $SCHED_TIME_AVG_MS > /proc/sys/kernel/sched_time_avg_ms &
fi

if [ -n "$SCHED_TUNABLE_SCALING" ] && [ -w /proc/sys/kernel/sched_tunable_scaling ]; then
  echo $SCHED_TUNABLE_SCALING > /proc/sys/kernel/sched_tunable_scaling &
fi

if [ -n "$SCHED_WAKEUP_GRANULARITY_NS" ] && [ -w /proc/sys/kernel/sched_wakeup_granularity_ns ]; then
  echo $SCHED_WAKEUP_GRANULARITY_NS > /proc/sys/kernel/sched_wakeup_granularity_ns &
fi

if [ -n "$SEM" ] && [ -w /proc/sys/kernel/sem ]; then
  echo $SEM > /proc/sys/kernel/sem &
fi

if [ -n "$SHM_RMID_FORCED" ] && [ -w /proc/sys/kernel/shm_rmid_forced ]; then
  echo $SHM_RMID_FORCED > /proc/sys/kernel/shm_rmid_forced &
fi

if [ -n "$SHMALL" ] && [ -w /proc/sys/kernel/shmall ]; then
  echo $SHMALL > /proc/sys/kernel/shmall &
fi

if [ -n "$SHMMAX" ] && [ -w /proc/sys/kernel/shmmax ]; then
  echo $SHMMAX > /proc/sys/kernel/shmmax &
fi

if [ -n "$SHMMNI" ] && [ -w /proc/sys/kernel/shmmni ]; then
  echo $SHMMNI > /proc/sys/kernel/shmmni &
fi

if [ -n "$SOFTLOCKUP_ALL_CPU_BACKTRACE" ] && [ -w /proc/sys/kernel/softlockup_all_cpu_backtrace ]; then
  echo $SOFTLOCKUP_ALL_CPU_BACKTRACE > /proc/sys/kernel/softlockup_all_cpu_backtrace &
fi

if [ -n "$SOFTLOCKUP_PANIC" ] && [ -w /proc/sys/kernel/softlockup_panic ]; then
  echo $SOFTLOCKUP_PANIC > /proc/sys/kernel/softlockup_panic &
fi

if [ -n "$STACK_TRACER_ENABLED" ] && [ -w /proc/sys/kernel/stack_tracer_enabled ]; then
  echo $STACK_TRACER_ENABLED > /proc/sys/kernel/stack_tracer_enabled &
fi

if [ -n "$SYSCTL_WRITES_STRICT" ] && [ -w /proc/sys/kernel/sysctl_writes_strict ]; then
  echo $SYSCTL_WRITES_STRICT > /proc/sys/kernel/sysctl_writes_strict &
fi

if [ -n "$SYSRQ" ] && [ -w /proc/sys/kernel/sysrq ]; then
  echo $SYSRQ > /proc/sys/kernel/sysrq &
fi

if [ -n "$TAINTED" ] && [ -w /proc/sys/kernel/tainted ]; then
  echo $TAINTED > /proc/sys/kernel/tainted &
fi

if [ -n "$THREADS_MAX" ] && [ -w /proc/sys/kernel/threads-max ]; then
  echo $THREADS_MAX > /proc/sys/kernel/threads-max &
fi

if [ -n "$TIMER_MIGRATION" ] && [ -w /proc/sys/kernel/timer_migration ]; then
  echo $TIMER_MIGRATION > /proc/sys/kernel/timer_migration &
fi

if [ -n "$TRACEOFF_ON_WARNING" ] && [ -w /proc/sys/kernel/traceoff_on_warning ]; then
  echo $TRACEOFF_ON_WARNING > /proc/sys/kernel/traceoff_on_warning &
fi

if [ -n "$TRACEPOINT_PRINTK" ] && [ -w /proc/sys/kernel/tracepoint_printk ]; then
  echo $TRACEPOINT_PRINTK > /proc/sys/kernel/tracepoint_printk &
fi

if [ -n "$UNKNOWN_NMI_PANIC" ] && [ -w /proc/sys/kernel/unknown_nmi_panic ]; then
  echo $UNKNOWN_NMI_PANIC > /proc/sys/kernel/unknown_nmi_panic &
fi

if [ -n "$WATCHDOG" ] && [ -w /proc/sys/kernel/watchdog ]; then
  echo $WATCHDOG > /proc/sys/kernel/watchdog &
fi

if [ -n "$WATCHDOG_THRESH" ] && [ -w /proc/sys/kernel/watchdog_thresh ]; then
  echo $WATCHDOG_THRESH > /proc/sys/kernel/watchdog_thresh &
fi

if [ -n "$BPF_JIT_ENABLE" ] && [ -w /proc/sys/net/core/bpf_jit_enable ]; then
  echo $BPF_JIT_ENABLE > /proc/sys/net/core/bpf_jit_enable &
fi

if [ -n "$BUSY_POLL" ] && [ -w /proc/sys/net/core/busy_poll ]; then
  echo $BUSY_POLL > /proc/sys/net/core/busy_poll &
fi

if [ -n "$BUSY_READ" ] && [ -w /proc/sys/net/core/busy_read ]; then
  echo $BUSY_READ > /proc/sys/net/core/busy_read &
fi

if [ -n "$DEFAULT_QDISC" ] && [ -w /proc/sys/net/core/default_qdisc ]; then
  echo $DEFAULT_QDISC > /proc/sys/net/core/default_qdisc &
fi

if [ -n "$DEV_WEIGHT" ] && [ -w /proc/sys/net/core/dev_weight ]; then
  echo $DEV_WEIGHT > /proc/sys/net/core/dev_weight &
fi

if [ -n "$FLOW_LIMIT_CPU_BITMAP" ] && [ -w /proc/sys/net/core/flow_limit_cpu_bitmap ]; then
  echo $FLOW_LIMIT_CPU_BITMAP > /proc/sys/net/core/flow_limit_cpu_bitmap &
fi

if [ -n "$FLOW_LIMIT_TABLE_LEN" ] && [ -w /proc/sys/net/core/flow_limit_table_len ]; then
  echo $FLOW_LIMIT_TABLE_LEN > /proc/sys/net/core/flow_limit_table_len &
fi

if [ -n "$MESSAGE_BURST" ] && [ -w /proc/sys/net/core/message_burst ]; then
  echo $MESSAGE_BURST > /proc/sys/net/core/message_burst &
fi

if [ -n "$MESSAGE_COST" ] && [ -w /proc/sys/net/core/message_cost ]; then
  echo $MESSAGE_COST > /proc/sys/net/core/message_cost &
fi

if [ -n "$NETDEV_BUDGET" ] && [ -w /proc/sys/net/core/netdev_budget ]; then
  echo $NETDEV_BUDGET > /proc/sys/net/core/netdev_budget &
fi

if [ -n "$NETDEV_MAX_BACKLOG" ] && [ -w /proc/sys/net/core/netdev_max_backlog ]; then
  echo $NETDEV_MAX_BACKLOG > /proc/sys/net/core/netdev_max_backlog &
fi

if [ -n "$NETDEV_TSTAMP_PREQUEUE" ] && [ -w /proc/sys/net/core/netdev_tstamp_prequeue ]; then
  echo $NETDEV_TSTAMP_PREQUEUE > /proc/sys/net/core/netdev_tstamp_prequeue &
fi

if [ -n "$OPTMEM_MAX" ] && [ -w /proc/sys/net/core/optmem_max ]; then
  echo $OPTMEM_MAX > /proc/sys/net/core/optmem_max &
fi

if [ -n "$RMEM_DEFAULT" ] && [ -w /proc/sys/net/core/rmem_default ]; then
  echo $RMEM_DEFAULT > /proc/sys/net/core/rmem_default &
fi

if [ -n "$RMEM_MAX" ] && [ -w /proc/sys/net/core/rmem_max ]; then
  echo $RMEM_MAX > /proc/sys/net/core/rmem_max &
fi

if [ -n "$RPS_SOCK_FLOW_ENTRIES" ] && [ -w /proc/sys/net/core/rps_sock_flow_entries ]; then
  echo $RPS_SOCK_FLOW_ENTRIES > /proc/sys/net/core/rps_sock_flow_entries &
fi

if [ -n "$SOMAXCONN" ] && [ -w /proc/sys/net/core/somaxconn ]; then
  echo $SOMAXCONN > /proc/sys/net/core/somaxconn &
fi

if [ -n "$TSTAMP_ALLOW_DATA" ] && [ -w /proc/sys/net/core/tstamp_allow_data ]; then
  echo $TSTAMP_ALLOW_DATA > /proc/sys/net/core/tstamp_allow_data &
fi

if [ -n "$WMEM_DEFAULT" ] && [ -w /proc/sys/net/core/wmem_default ]; then
  echo $WMEM_DEFAULT > /proc/sys/net/core/wmem_default &
fi

if [ -n "$WMEM_MAX" ] && [ -w /proc/sys/net/core/wmem_max ]; then
  echo $WMEM_MAX > /proc/sys/net/core/wmem_max &
fi

if [ -n "$RPS_FLOW_CNT" ] && [ -w /sys/class/net/eth0/queues/rx-0/rps_flow_cnt ]; then
  echo $RPS_FLOW_CNT > /sys/class/net/eth0/queues/rx-0/rps_flow_cnt &
fi

if [ -n "$RPS_CPUS" ] && [ -w /sys/class/net/eth0/queues/rx-0/rps_cpus ]; then
  echo $RPS_CPUS > /sys/class/net/eth0/queues/rx-0/rps_cpus &
fi

if [ -n "$FWMARK_REFLECT" ] && [ -w /proc/sys/net/ipv4/fwmark_reflect ]; then
  echo $FWMARK_REFLECT > /proc/sys/net/ipv4/fwmark_reflect &
fi

if [ -n "$ICMP_ECHO_IGNORE_ALL" ] && [ -w /proc/sys/net/ipv4/icmp_echo_ignore_all ]; then
  echo $ICMP_ECHO_IGNORE_ALL > /proc/sys/net/ipv4/icmp_echo_ignore_all &
fi

if [ -n "$ICMP_ECHO_IGNORE_BROADCASTS" ] && [ -w /proc/sys/net/ipv4/icmp_echo_ignore_broadcasts ]; then
  echo $ICMP_ECHO_IGNORE_BROADCASTS > /proc/sys/net/ipv4/icmp_echo_ignore_broadcasts &
fi

if [ -n "$ICMP_ERRORS_USE_INBOUND_IFADDR" ] && [ -w /proc/sys/net/ipv4/icmp_errors_use_inbound_ifaddr ]; then
  echo $ICMP_ERRORS_USE_INBOUND_IFADDR > /proc/sys/net/ipv4/icmp_errors_use_inbound_ifaddr &
fi

if [ -n "$ICMP_IGNORE_BOGUS_ERROR_RESPONSES" ] && [ -w /proc/sys/net/ipv4/icmp_ignore_bogus_error_responses ]; then
  echo $ICMP_IGNORE_BOGUS_ERROR_RESPONSES > /proc/sys/net/ipv4/icmp_ignore_bogus_error_responses &
fi

if [ -n "$ICMP_MSGS_BURST" ] && [ -w /proc/sys/net/ipv4/icmp_msgs_burst ]; then
  echo $ICMP_MSGS_BURST > /proc/sys/net/ipv4/icmp_msgs_burst &
fi

if [ -n "$ICMP_MSGS_PER_SEC" ] && [ -w /proc/sys/net/ipv4/icmp_msgs_per_sec ]; then
  echo $ICMP_MSGS_PER_SEC > /proc/sys/net/ipv4/icmp_msgs_per_sec &
fi

if [ -n "$ICMP_RATELIMIT" ] && [ -w /proc/sys/net/ipv4/icmp_ratelimit ]; then
  echo $ICMP_RATELIMIT > /proc/sys/net/ipv4/icmp_ratelimit &
fi

if [ -n "$ICMP_RATEMASK" ] && [ -w /proc/sys/net/ipv4/icmp_ratemask ]; then
  echo $ICMP_RATEMASK > /proc/sys/net/ipv4/icmp_ratemask &
fi

if [ -n "$IGMP_MAX_MEMBERSHIPS" ] && [ -w /proc/sys/net/ipv4/igmp_max_memberships ]; then
  echo $IGMP_MAX_MEMBERSHIPS > /proc/sys/net/ipv4/igmp_max_memberships &
fi

if [ -n "$IGMP_MAX_MSF" ] && [ -w /proc/sys/net/ipv4/igmp_max_msf ]; then
  echo $IGMP_MAX_MSF > /proc/sys/net/ipv4/igmp_max_msf &
fi

if [ -n "$IGMP_QRV" ] && [ -w /proc/sys/net/ipv4/igmp_qrv ]; then
  echo $IGMP_QRV > /proc/sys/net/ipv4/igmp_qrv &
fi

if [ -n "$INET_PEER_MAXTTL" ] && [ -w /proc/sys/net/ipv4/inet_peer_maxttl ]; then
  echo $INET_PEER_MAXTTL > /proc/sys/net/ipv4/inet_peer_maxttl &
fi

if [ -n "$INET_PEER_MINTTL" ] && [ -w /proc/sys/net/ipv4/inet_peer_minttl ]; then
  echo $INET_PEER_MINTTL > /proc/sys/net/ipv4/inet_peer_minttl &
fi

if [ -n "$INET_PEER_THRESHOLD" ] && [ -w /proc/sys/net/ipv4/inet_peer_threshold ]; then
  echo $INET_PEER_THRESHOLD > /proc/sys/net/ipv4/inet_peer_threshold &
fi

if [ -n "$IP_DEFAULT_TTL" ] && [ -w /proc/sys/net/ipv4/ip_default_ttl ]; then
  echo $IP_DEFAULT_TTL > /proc/sys/net/ipv4/ip_default_ttl &
fi

if [ -n "$IP_DYNADDR" ] && [ -w /proc/sys/net/ipv4/ip_dynaddr ]; then
  echo $IP_DYNADDR > /proc/sys/net/ipv4/ip_dynaddr &
fi

if [ -n "$IP_EARLY_DEMUX" ] && [ -w /proc/sys/net/ipv4/ip_early_demux ]; then
  echo $IP_EARLY_DEMUX > /proc/sys/net/ipv4/ip_early_demux &
fi

if [ -n "$IP_FORWARD" ] && [ -w /proc/sys/net/ipv4/ip_forward ]; then
  echo $IP_FORWARD > /proc/sys/net/ipv4/ip_forward &
fi

if [ -n "$IP_FORWARD_USE_PMTU" ] && [ -w /proc/sys/net/ipv4/ip_forward_use_pmtu ]; then
  echo $IP_FORWARD_USE_PMTU > /proc/sys/net/ipv4/ip_forward_use_pmtu &
fi

if [ -n "$IP_NO_PMTU_DISC" ] && [ -w /proc/sys/net/ipv4/ip_no_pmtu_disc ]; then
  echo $IP_NO_PMTU_DISC > /proc/sys/net/ipv4/ip_no_pmtu_disc &
fi

if [ -n "$IP_NONLOCAL_BIND" ] && [ -w /proc/sys/net/ipv4/ip_nonlocal_bind ]; then
  echo $IP_NONLOCAL_BIND > /proc/sys/net/ipv4/ip_nonlocal_bind &
fi

if [ -n "$IPFRAG_HIGH_THRESH" ] && [ -w /proc/sys/net/ipv4/ipfrag_high_thresh ]; then
  echo $IPFRAG_HIGH_THRESH > /proc/sys/net/ipv4/ipfrag_high_thresh &
fi

if [ -n "$IPFRAG_LOW_THRESH" ] && [ -w /proc/sys/net/ipv4/ipfrag_low_thresh ]; then
  echo $IPFRAG_LOW_THRESH > /proc/sys/net/ipv4/ipfrag_low_thresh &
fi

if [ -n "$IPFRAG_MAX_DIST" ] && [ -w /proc/sys/net/ipv4/ipfrag_max_dist ]; then
  echo $IPFRAG_MAX_DIST > /proc/sys/net/ipv4/ipfrag_max_dist &
fi

if [ -n "$IPFRAG_SECRET_INTERVAL" ] && [ -w /proc/sys/net/ipv4/ipfrag_secret_interval ]; then
  echo $IPFRAG_SECRET_INTERVAL > /proc/sys/net/ipv4/ipfrag_secret_interval &
fi

if [ -n "$IPFRAG_TIME" ] && [ -w /proc/sys/net/ipv4/ipfrag_time ]; then
  echo $IPFRAG_TIME > /proc/sys/net/ipv4/ipfrag_time &
fi

if [ -n "$TCP_ABORT_ON_OVERFLOW" ] && [ -w /proc/sys/net/ipv4/tcp_abort_on_overflow ]; then
  echo $TCP_ABORT_ON_OVERFLOW > /proc/sys/net/ipv4/tcp_abort_on_overflow &
fi

if [ -n "$TCP_ADV_WIN_SCALE" ] && [ -w /proc/sys/net/ipv4/tcp_adv_win_scale ]; then
  echo $TCP_ADV_WIN_SCALE > /proc/sys/net/ipv4/tcp_adv_win_scale &
fi

if [ -n "$TCP_APP_WIN" ] && [ -w /proc/sys/net/ipv4/tcp_app_win ]; then
  echo $TCP_APP_WIN > /proc/sys/net/ipv4/tcp_app_win &
fi

if [ -n "$TCP_AUTOCORKING" ] && [ -w /proc/sys/net/ipv4/tcp_autocorking ]; then
  echo $TCP_AUTOCORKING > /proc/sys/net/ipv4/tcp_autocorking &
fi

if [ -n "$TCP_BASE_MSS" ] && [ -w /proc/sys/net/ipv4/tcp_base_mss ]; then
  echo $TCP_BASE_MSS > /proc/sys/net/ipv4/tcp_base_mss &
fi

if [ -n "$TCP_CHALLENGE_ACK_LIMIT" ] && [ -w /proc/sys/net/ipv4/tcp_challenge_ack_limit ]; then
  echo $TCP_CHALLENGE_ACK_LIMIT > /proc/sys/net/ipv4/tcp_challenge_ack_limit &
fi

if [ -n "$TCP_CONGESTION_CONTROL" ] && [ -w /proc/sys/net/ipv4/tcp_congestion_control ]; then
  echo $TCP_CONGESTION_CONTROL > /proc/sys/net/ipv4/tcp_congestion_control &
fi

if [ -n "$TCP_DSACK" ] && [ -w /proc/sys/net/ipv4/tcp_dsack ]; then
  echo $TCP_DSACK > /proc/sys/net/ipv4/tcp_dsack &
fi

if [ -n "$TCP_EARLY_RETRANS" ] && [ -w /proc/sys/net/ipv4/tcp_early_retrans ]; then
  echo $TCP_EARLY_RETRANS > /proc/sys/net/ipv4/tcp_early_retrans &
fi

if [ -n "$TCP_ECN" ] && [ -w /proc/sys/net/ipv4/tcp_ecn ]; then
  echo $TCP_ECN > /proc/sys/net/ipv4/tcp_ecn &
fi

if [ -n "$TCP_FACK" ] && [ -w /proc/sys/net/ipv4/tcp_fack ]; then
  echo $TCP_FACK > /proc/sys/net/ipv4/tcp_fack &
fi

if [ -n "$TCP_FASTOPEN" ] && [ -w /proc/sys/net/ipv4/tcp_fastopen ]; then
  echo $TCP_FASTOPEN > /proc/sys/net/ipv4/tcp_fastopen &
fi

if [ -n "$TCP_FIN_TIMEOUT" ] && [ -w /proc/sys/net/ipv4/tcp_fin_timeout ]; then
  echo $TCP_FIN_TIMEOUT > /proc/sys/net/ipv4/tcp_fin_timeout &
fi

if [ -n "$TCP_FRTO" ] && [ -w /proc/sys/net/ipv4/tcp_frto ]; then
  echo $TCP_FRTO > /proc/sys/net/ipv4/tcp_frto &
fi

if [ -n "$TCP_FWMARK_ACCEPT" ] && [ -w /proc/sys/net/ipv4/tcp_fwmark_accept ]; then
  echo $TCP_FWMARK_ACCEPT > /proc/sys/net/ipv4/tcp_fwmark_accept &
fi

if [ -n "$TCP_INVALID_RATELIMIT" ] && [ -w /proc/sys/net/ipv4/tcp_invalid_ratelimit ]; then
  echo $TCP_INVALID_RATELIMIT > /proc/sys/net/ipv4/tcp_invalid_ratelimit &
fi

if [ -n "$TCP_KEEPALIVE_INTVL" ] && [ -w /proc/sys/net/ipv4/tcp_keepalive_intvl ]; then
  echo $TCP_KEEPALIVE_INTVL > /proc/sys/net/ipv4/tcp_keepalive_intvl &
fi

if [ -n "$TCP_KEEPALIVE_PROBES" ] && [ -w /proc/sys/net/ipv4/tcp_keepalive_probes ]; then
  echo $TCP_KEEPALIVE_PROBES > /proc/sys/net/ipv4/tcp_keepalive_probes &
fi

if [ -n "$TCP_KEEPALIVE_TIME" ] && [ -w /proc/sys/net/ipv4/tcp_keepalive_time ]; then
  echo $TCP_KEEPALIVE_TIME > /proc/sys/net/ipv4/tcp_keepalive_time &
fi

if [ -n "$TCP_LIMIT_OUTPUT_BYTES" ] && [ -w /proc/sys/net/ipv4/tcp_limit_output_bytes ]; then
  echo $TCP_LIMIT_OUTPUT_BYTES > /proc/sys/net/ipv4/tcp_limit_output_bytes &
fi

if [ -n "$TCP_LOW_LATENCY" ] && [ -w /proc/sys/net/ipv4/tcp_low_latency ]; then
  echo $TCP_LOW_LATENCY > /proc/sys/net/ipv4/tcp_low_latency &
fi

if [ -n "$TCP_MAX_ORPHANS" ] && [ -w /proc/sys/net/ipv4/tcp_max_orphans ]; then
  echo $TCP_MAX_ORPHANS > /proc/sys/net/ipv4/tcp_max_orphans &
fi

if [ -n "$TCP_MAX_REORDERING" ] && [ -w /proc/sys/net/ipv4/tcp_max_reordering ]; then
  echo $TCP_MAX_REORDERING > /proc/sys/net/ipv4/tcp_max_reordering &
fi

if [ -n "$TCP_MAX_SYN_BACKLOG" ] && [ -w /proc/sys/net/ipv4/tcp_max_syn_backlog ]; then
  echo $TCP_MAX_SYN_BACKLOG > /proc/sys/net/ipv4/tcp_max_syn_backlog &
fi

if [ -n "$TCP_MAX_TW_BUCKETS" ] && [ -w /proc/sys/net/ipv4/tcp_max_tw_buckets ]; then
  echo $TCP_MAX_TW_BUCKETS > /proc/sys/net/ipv4/tcp_max_tw_buckets &
fi

if [ -n "$TCP_MEM" ] && [ -w /proc/sys/net/ipv4/tcp_mem ]; then
  echo $TCP_MEM > /proc/sys/net/ipv4/tcp_mem &
fi

if [ -n "$TCP_MIN_TSO_SEGS" ] && [ -w /proc/sys/net/ipv4/tcp_min_tso_segs ]; then
  echo $TCP_MIN_TSO_SEGS > /proc/sys/net/ipv4/tcp_min_tso_segs &
fi

if [ -n "$TCP_MODERATE_RCVBUF" ] && [ -w /proc/sys/net/ipv4/tcp_moderate_rcvbuf ]; then
  echo $TCP_MODERATE_RCVBUF > /proc/sys/net/ipv4/tcp_moderate_rcvbuf &
fi

if [ -n "$TCP_MTU_PROBING" ] && [ -w /proc/sys/net/ipv4/tcp_mtu_probing ]; then
  echo $TCP_MTU_PROBING > /proc/sys/net/ipv4/tcp_mtu_probing &
fi

if [ -n "$TCP_NO_METRICS_SAVE" ] && [ -w /proc/sys/net/ipv4/tcp_no_metrics_save ]; then
  echo $TCP_NO_METRICS_SAVE > /proc/sys/net/ipv4/tcp_no_metrics_save &
fi

if [ -n "$TCP_NOTSENT_LOWAT" ] && [ -w /proc/sys/net/ipv4/tcp_notsent_lowat ]; then
  echo $TCP_NOTSENT_LOWAT > /proc/sys/net/ipv4/tcp_notsent_lowat &
fi

if [ -n "$TCP_ORPHAN_RETRIES" ] && [ -w /proc/sys/net/ipv4/tcp_orphan_retries ]; then
  echo $TCP_ORPHAN_RETRIES > /proc/sys/net/ipv4/tcp_orphan_retries &
fi

if [ -n "$TCP_REORDERING" ] && [ -w /proc/sys/net/ipv4/tcp_reordering ]; then
  echo $TCP_REORDERING > /proc/sys/net/ipv4/tcp_reordering &
fi

if [ -n "$TCP_RETRANS_COLLAPSE" ] && [ -w /proc/sys/net/ipv4/tcp_retrans_collapse ]; then
  echo $TCP_RETRANS_COLLAPSE > /proc/sys/net/ipv4/tcp_retrans_collapse &
fi

if [ -n "$TCP_RETRIES1" ] && [ -w /proc/sys/net/ipv4/tcp_retries1 ]; then
  echo $TCP_RETRIES1 > /proc/sys/net/ipv4/tcp_retries1 &
fi

if [ -n "$TCP_RETRIES2" ] && [ -w /proc/sys/net/ipv4/tcp_retries2 ]; then
  echo $TCP_RETRIES2 > /proc/sys/net/ipv4/tcp_retries2 &
fi

if [ -n "$TCP_RFC1337" ] && [ -w /proc/sys/net/ipv4/tcp_rfc1337 ]; then
  echo $TCP_RFC1337 > /proc/sys/net/ipv4/tcp_rfc1337 &
fi

if [ -n "$TCP_RMEM" ] && [ -w /proc/sys/net/ipv4/tcp_rmem ]; then
  echo $TCP_RMEM > /proc/sys/net/ipv4/tcp_rmem &
fi

if [ -n "$TCP_SACK" ] && [ -w /proc/sys/net/ipv4/tcp_sack ]; then
  echo $TCP_SACK > /proc/sys/net/ipv4/tcp_sack &
fi

if [ -n "$TCP_SLOW_START_AFTER_IDLE" ] && [ -w /proc/sys/net/ipv4/tcp_slow_start_after_idle ]; then
  echo $TCP_SLOW_START_AFTER_IDLE > /proc/sys/net/ipv4/tcp_slow_start_after_idle &
fi

if [ -n "$TCP_STDURG" ] && [ -w /proc/sys/net/ipv4/tcp_stdurg ]; then
  echo $TCP_STDURG > /proc/sys/net/ipv4/tcp_stdurg &
fi

if [ -n "$TCP_SYN_RETRIES" ] && [ -w /proc/sys/net/ipv4/tcp_syn_retries ]; then
  echo $TCP_SYN_RETRIES > /proc/sys/net/ipv4/tcp_syn_retries &
fi

if [ -n "$TCP_SYNACK_RETRIES" ] && [ -w /proc/sys/net/ipv4/tcp_synack_retries ]; then
  echo $TCP_SYNACK_RETRIES > /proc/sys/net/ipv4/tcp_synack_retries &
fi

if [ -n "$TCP_SYNCOOKIES" ] && [ -w /proc/sys/net/ipv4/tcp_syncookies ]; then
  echo $TCP_SYNCOOKIES > /proc/sys/net/ipv4/tcp_syncookies &
fi

if [ -n "$TCP_THIN_DUPACK" ] && [ -w /proc/sys/net/ipv4/tcp_thin_dupack ]; then
  echo $TCP_THIN_DUPACK > /proc/sys/net/ipv4/tcp_thin_dupack &
fi

if [ -n "$TCP_THIN_LINEAR_TIMEOUTS" ] && [ -w /proc/sys/net/ipv4/tcp_thin_linear_timeouts ]; then
  echo $TCP_THIN_LINEAR_TIMEOUTS > /proc/sys/net/ipv4/tcp_thin_linear_timeouts &
fi

if [ -n "$TCP_TIMESTAMPS" ] && [ -w /proc/sys/net/ipv4/tcp_timestamps ]; then
  echo $TCP_TIMESTAMPS > /proc/sys/net/ipv4/tcp_timestamps &
fi

if [ -n "$TCP_TSO_WIN_DIVISOR" ] && [ -w /proc/sys/net/ipv4/tcp_tso_win_divisor ]; then
  echo $TCP_TSO_WIN_DIVISOR > /proc/sys/net/ipv4/tcp_tso_win_divisor &
fi

if [ -n "$TCP_TW_RECYCLE" ] && [ -w /proc/sys/net/ipv4/tcp_tw_recycle ]; then
  echo $TCP_TW_RECYCLE > /proc/sys/net/ipv4/tcp_tw_recycle &
fi

if [ -n "$TCP_TW_REUSE" ] && [ -w /proc/sys/net/ipv4/tcp_tw_reuse ]; then
  echo $TCP_TW_REUSE > /proc/sys/net/ipv4/tcp_tw_reuse &
fi

if [ -n "$TCP_WINDOW_SCALING" ] && [ -w /proc/sys/net/ipv4/tcp_window_scaling ]; then
  echo $TCP_WINDOW_SCALING > /proc/sys/net/ipv4/tcp_window_scaling &
fi

if [ -n "$TCP_WMEM" ] && [ -w /proc/sys/net/ipv4/tcp_wmem ]; then
  echo $TCP_WMEM > /proc/sys/net/ipv4/tcp_wmem &
fi

if [ -n "$TCP_WORKAROUND_SIGNED_WINDOWS" ] && [ -w /proc/sys/net/ipv4/tcp_workaround_signed_windows ]; then
  echo $TCP_WORKAROUND_SIGNED_WINDOWS > /proc/sys/net/ipv4/tcp_workaround_signed_windows &
fi

if [ -n "$UDP_MEM" ] && [ -w /proc/sys/net/ipv4/udp_mem ]; then
  echo $UDP_MEM > /proc/sys/net/ipv4/udp_mem &
fi

if [ -n "$UDP_RMEM_MIN" ] && [ -w /proc/sys/net/ipv4/udp_rmem_min ]; then
  echo $UDP_RMEM_MIN > /proc/sys/net/ipv4/udp_rmem_min &
fi

if [ -n "$UDP_WMEM_MIN" ] && [ -w /proc/sys/net/ipv4/udp_wmem_min ]; then
  echo $UDP_WMEM_MIN > /proc/sys/net/ipv4/udp_wmem_min &
fi

if [ -n "$XFRM4_GC_THRESH" ] && [ -w /proc/sys/net/ipv4/xfrm4_gc_thresh ]; then
  echo $XFRM4_GC_THRESH > /proc/sys/net/ipv4/xfrm4_gc_thresh &
fi

if [ -n "$ERROR_BURST" ] && [ -w /proc/sys/net/ipv4/route/error_burst ]; then
  echo $ERROR_BURST > /proc/sys/net/ipv4/route/error_burst &
fi

if [ -n "$ERROR_COST" ] && [ -w /proc/sys/net/ipv4/route/error_cost ]; then
  echo $ERROR_COST > /proc/sys/net/ipv4/route/error_cost &
fi

if [ -n "$GC_ELASTICITY" ] && [ -w /proc/sys/net/ipv4/route/gc_elasticity ]; then
  echo $GC_ELASTICITY > /proc/sys/net/ipv4/route/gc_elasticity &
fi

if [ -n "$GC_INTERVAL" ] && [ -w /proc/sys/net/ipv4/route/gc_interval ]; then
  echo $GC_INTERVAL > /proc/sys/net/ipv4/route/gc_interval &
fi

if [ -n "$GC_MIN_INTERVAL" ] && [ -w /proc/sys/net/ipv4/route/gc_min_interval ]; then
  echo $GC_MIN_INTERVAL > /proc/sys/net/ipv4/route/gc_min_interval &
fi

if [ -n "$GC_MIN_INTERVAL_MS" ] && [ -w /proc/sys/net/ipv4/route/gc_min_interval_ms ]; then
  echo $GC_MIN_INTERVAL_MS > /proc/sys/net/ipv4/route/gc_min_interval_ms &
fi

if [ -n "$GC_THRESH" ] && [ -w /proc/sys/net/ipv4/route/gc_thresh ]; then
  echo $GC_THRESH > /proc/sys/net/ipv4/route/gc_thresh &
fi

if [ -n "$GC_TIMEOUT" ] && [ -w /proc/sys/net/ipv4/route/gc_timeout ]; then
  echo $GC_TIMEOUT > /proc/sys/net/ipv4/route/gc_timeout &
fi

if [ -n "$MAX_SIZE" ] && [ -w /proc/sys/net/ipv4/route/max_size ]; then
  echo $MAX_SIZE > /proc/sys/net/ipv4/route/max_size &
fi

if [ -n "$MIN_ADV_MSS" ] && [ -w /proc/sys/net/ipv4/route/min_adv_mss ]; then
  echo $MIN_ADV_MSS > /proc/sys/net/ipv4/route/min_adv_mss &
fi

if [ -n "$MIN_PMTU" ] && [ -w /proc/sys/net/ipv4/route/min_pmtu ]; then
  echo $MIN_PMTU > /proc/sys/net/ipv4/route/min_pmtu &
fi

if [ -n "$MTU_EXPIRES" ] && [ -w /proc/sys/net/ipv4/route/mtu_expires ]; then
  echo $MTU_EXPIRES > /proc/sys/net/ipv4/route/mtu_expires &
fi

if [ -n "$REDIRECT_LOAD" ] && [ -w /proc/sys/net/ipv4/route/redirect_load ]; then
  echo $REDIRECT_LOAD > /proc/sys/net/ipv4/route/redirect_load &
fi

if [ -n "$REDIRECT_NUMBER" ] && [ -w /proc/sys/net/ipv4/route/redirect_number ]; then
  echo $REDIRECT_NUMBER > /proc/sys/net/ipv4/route/redirect_number &
fi

if [ -n "$REDIRECT_SILENCE" ] && [ -w /proc/sys/net/ipv4/route/redirect_silence ]; then
  echo $REDIRECT_SILENCE > /proc/sys/net/ipv4/route/redirect_silence &
fi

if [ -n "$ACCEPT_LOCAL" ] && [ -w /proc/sys/net/ipv4/conf/all/accept_local ]; then
  echo $ACCEPT_LOCAL > /proc/sys/net/ipv4/conf/all/accept_local &
fi

if [ -n "$ACCEPT_REDIRECTS" ] && [ -w /proc/sys/net/ipv4/conf/all/accept_redirects ]; then
  echo $ACCEPT_REDIRECTS > /proc/sys/net/ipv4/conf/all/accept_redirects &
fi

if [ -n "$ACCEPT_SOURCE_ROUTE" ] && [ -w /proc/sys/net/ipv4/conf/all/accept_source_route ]; then
  echo $ACCEPT_SOURCE_ROUTE > /proc/sys/net/ipv4/conf/all/accept_source_route &
fi

if [ -n "$ARP_ACCEPT" ] && [ -w /proc/sys/net/ipv4/conf/all/arp_accept ]; then
  echo $ARP_ACCEPT > /proc/sys/net/ipv4/conf/all/arp_accept &
fi

if [ -n "$ARP_ANNOUNCE" ] && [ -w /proc/sys/net/ipv4/conf/all/arp_announce ]; then
  echo $ARP_ANNOUNCE > /proc/sys/net/ipv4/conf/all/arp_announce &
fi

if [ -n "$ARP_FILTER" ] && [ -w /proc/sys/net/ipv4/conf/all/arp_filter ]; then
  echo $ARP_FILTER > /proc/sys/net/ipv4/conf/all/arp_filter &
fi

if [ -n "$ARP_IGNORE" ] && [ -w /proc/sys/net/ipv4/conf/all/arp_ignore ]; then
  echo $ARP_IGNORE > /proc/sys/net/ipv4/conf/all/arp_ignore &
fi

if [ -n "$ARP_NOTIFY" ] && [ -w /proc/sys/net/ipv4/conf/all/arp_notify ]; then
  echo $ARP_NOTIFY > /proc/sys/net/ipv4/conf/all/arp_notify &
fi

if [ -n "$BOOTP_RELAY" ] && [ -w /proc/sys/net/ipv4/conf/all/bootp_relay ]; then
  echo $BOOTP_RELAY > /proc/sys/net/ipv4/conf/all/bootp_relay &
fi

if [ -n "$DISABLE_POLICY" ] && [ -w /proc/sys/net/ipv4/conf/all/disable_policy ]; then
  echo $DISABLE_POLICY > /proc/sys/net/ipv4/conf/all/disable_policy &
fi

if [ -n "$DISABLE_XFRM" ] && [ -w /proc/sys/net/ipv4/conf/all/disable_xfrm ]; then
  echo $DISABLE_XFRM > /proc/sys/net/ipv4/conf/all/disable_xfrm &
fi

if [ -n "$FORCE_IGMP_VERSION" ] && [ -w /proc/sys/net/ipv4/conf/all/force_igmp_version ]; then
  echo $FORCE_IGMP_VERSION > /proc/sys/net/ipv4/conf/all/force_igmp_version &
fi

if [ -n "$FORWARDING" ] && [ -w /proc/sys/net/ipv4/conf/all/forwarding ]; then
  echo $FORWARDING > /proc/sys/net/ipv4/conf/all/forwarding &
fi

if [ -n "$IGMPV2_UNSOLICITED_REPORT_INTERVAL" ] && [ -w /proc/sys/net/ipv4/conf/all/igmpv2_unsolicited_report_interval ]; then
  echo $IGMPV2_UNSOLICITED_REPORT_INTERVAL > /proc/sys/net/ipv4/conf/all/igmpv2_unsolicited_report_interval &
fi

if [ -n "$IGMPV3_UNSOLICITED_REPORT_INTERVAL" ] && [ -w /proc/sys/net/ipv4/conf/all/igmpv3_unsolicited_report_interval ]; then
  echo $IGMPV3_UNSOLICITED_REPORT_INTERVAL > /proc/sys/net/ipv4/conf/all/igmpv3_unsolicited_report_interval &
fi

if [ -n "$LOG_MARTIANS" ] && [ -w /proc/sys/net/ipv4/conf/all/log_martians ]; then
  echo $LOG_MARTIANS > /proc/sys/net/ipv4/conf/all/log_martians &
fi

if [ -n "$MEDIUM_ID" ] && [ -w /proc/sys/net/ipv4/conf/all/medium_id ]; then
  echo $MEDIUM_ID > /proc/sys/net/ipv4/conf/all/medium_id &
fi

if [ -n "$PROMOTE_SECONDARIES" ] && [ -w /proc/sys/net/ipv4/conf/all/promote_secondaries ]; then
  echo $PROMOTE_SECONDARIES > /proc/sys/net/ipv4/conf/all/promote_secondaries &
fi

if [ -n "$PROXY_ARP" ] && [ -w /proc/sys/net/ipv4/conf/all/proxy_arp ]; then
  echo $PROXY_ARP > /proc/sys/net/ipv4/conf/all/proxy_arp &
fi

if [ -n "$PROXY_ARP_PVLAN" ] && [ -w /proc/sys/net/ipv4/conf/all/proxy_arp_pvlan ]; then
  echo $PROXY_ARP_PVLAN > /proc/sys/net/ipv4/conf/all/proxy_arp_pvlan &
fi

if [ -n "$ROUTE_LOCALNET" ] && [ -w /proc/sys/net/ipv4/conf/all/route_localnet ]; then
  echo $ROUTE_LOCALNET > /proc/sys/net/ipv4/conf/all/route_localnet &
fi

if [ -n "$RP_FILTER" ] && [ -w /proc/sys/net/ipv4/conf/all/rp_filter ]; then
  echo $RP_FILTER > /proc/sys/net/ipv4/conf/all/rp_filter &
fi

if [ -n "$SECURE_REDIRECTS" ] && [ -w /proc/sys/net/ipv4/conf/all/secure_redirects ]; then
  echo $SECURE_REDIRECTS > /proc/sys/net/ipv4/conf/all/secure_redirects &
fi

if [ -n "$SEND_REDIRECTS" ] && [ -w /proc/sys/net/ipv4/conf/all/send_redirects ]; then
  echo $SEND_REDIRECTS > /proc/sys/net/ipv4/conf/all/send_redirects &
fi

if [ -n "$SHARED_MEDIA" ] && [ -w /proc/sys/net/ipv4/conf/all/shared_media ]; then
  echo $SHARED_MEDIA > /proc/sys/net/ipv4/conf/all/shared_media &
fi

if [ -n "$SRC_VALID_MARK" ] && [ -w /proc/sys/net/ipv4/conf/all/src_valid_mark ]; then
  echo $SRC_VALID_MARK > /proc/sys/net/ipv4/conf/all/src_valid_mark &
fi

if [ -n "$TAG" ] && [ -w /proc/sys/net/ipv4/conf/all/tag ]; then
  echo $TAG > /proc/sys/net/ipv4/conf/all/tag &
fi

if [ -n "$ANYCAST_DELAY" ] && [ -w /proc/sys/net/ipv4/neigh/eth0/anycast_delay ]; then
  echo $ANYCAST_DELAY > /proc/sys/net/ipv4/neigh/eth0/anycast_delay &
fi

if [ -n "$APP_SOLICIT" ] && [ -w /proc/sys/net/ipv4/neigh/eth0/app_solicit ]; then
  echo $APP_SOLICIT > /proc/sys/net/ipv4/neigh/eth0/app_solicit &
fi

if [ -n "$BASE_REACHABLE_TIME" ] && [ -w /proc/sys/net/ipv4/neigh/eth0/base_reachable_time ]; then
  echo $BASE_REACHABLE_TIME > /proc/sys/net/ipv4/neigh/eth0/base_reachable_time &
fi

if [ -n "$BASE_REACHABLE_TIME_MS" ] && [ -w /proc/sys/net/ipv4/neigh/eth0/base_reachable_time_ms ]; then
  echo $BASE_REACHABLE_TIME_MS > /proc/sys/net/ipv4/neigh/eth0/base_reachable_time_ms &
fi

if [ -n "$DELAY_FIRST_PROBE_TIME" ] && [ -w /proc/sys/net/ipv4/neigh/eth0/delay_first_probe_time ]; then
  echo $DELAY_FIRST_PROBE_TIME > /proc/sys/net/ipv4/neigh/eth0/delay_first_probe_time &
fi

if [ -n "$GC_STALE_TIME" ] && [ -w /proc/sys/net/ipv4/neigh/eth0/gc_stale_time ]; then
  echo $GC_STALE_TIME > /proc/sys/net/ipv4/neigh/eth0/gc_stale_time &
fi

if [ -n "$LOCKTIME" ] && [ -w /proc/sys/net/ipv4/neigh/eth0/locktime ]; then
  echo $LOCKTIME > /proc/sys/net/ipv4/neigh/eth0/locktime &
fi

if [ -n "$MCAST_SOLICIT" ] && [ -w /proc/sys/net/ipv4/neigh/eth0/mcast_solicit ]; then
  echo $MCAST_SOLICIT > /proc/sys/net/ipv4/neigh/eth0/mcast_solicit &
fi

if [ -n "$PROXY_DELAY" ] && [ -w /proc/sys/net/ipv4/neigh/eth0/proxy_delay ]; then
  echo $PROXY_DELAY > /proc/sys/net/ipv4/neigh/eth0/proxy_delay &
fi

if [ -n "$PROXY_QLEN" ] && [ -w /proc/sys/net/ipv4/neigh/eth0/proxy_qlen ]; then
  echo $PROXY_QLEN > /proc/sys/net/ipv4/neigh/eth0/proxy_qlen &
fi

if [ -n "$RETRANS_TIME" ] && [ -w /proc/sys/net/ipv4/neigh/eth0/retrans_time ]; then
  echo $RETRANS_TIME > /proc/sys/net/ipv4/neigh/eth0/retrans_time &
fi

if [ -n "$RETRANS_TIME_MS" ] && [ -w /proc/sys/net/ipv4/neigh/eth0/retrans_time_ms ]; then
  echo $RETRANS_TIME_MS > /proc/sys/net/ipv4/neigh/eth0/retrans_time_ms &
fi

if [ -n "$UCAST_SOLICIT" ] && [ -w /proc/sys/net/ipv4/neigh/eth0/ucast_solicit ]; then
  echo $UCAST_SOLICIT > /proc/sys/net/ipv4/neigh/eth0/ucast_solicit &
fi

if [ -n "$UNRES_QLEN" ] && [ -w /proc/sys/net/ipv4/neigh/eth0/unres_qlen ]; then
  echo $UNRES_QLEN > /proc/sys/net/ipv4/neigh/eth0/unres_qlen &
fi

if [ -n "$UNRES_QLEN_BYTES" ] && [ -w /proc/sys/net/ipv4/neigh/eth0/unres_qlen_bytes ]; then
  echo $UNRES_QLEN_BYTES > /proc/sys/net/ipv4/neigh/eth0/unres_qlen_bytes &
fi

if [ -n "$MAX_DGRAM_QLEN" ] && [ -w /proc/sys/net/unix/max_dgram_qlen ]; then
  echo $MAX_DGRAM_QLEN > /proc/sys/net/unix/max_dgram_qlen &
fi

if [ -n "$ADMIN_RESERVE_KBYTES" ] && [ -w /proc/sys/vm/admin_reserve_kbytes ]; then
  echo $ADMIN_RESERVE_KBYTES > /proc/sys/vm/admin_reserve_kbytes &
fi

if [ -n "$BLOCK_DUMP" ] && [ -w /proc/sys/vm/block_dump ]; then
  echo $BLOCK_DUMP > /proc/sys/vm/block_dump &
fi

if [ -n "$DIRTY_BACKGROUND_BYTES" ] && [ -w /proc/sys/vm/dirty_background_bytes ]; then
  echo $DIRTY_BACKGROUND_BYTES > /proc/sys/vm/dirty_background_bytes &
fi

if [ -n "$DIRTY_BACKGROUND_RATIO" ] && [ -w /proc/sys/vm/dirty_background_ratio ]; then
  echo $DIRTY_BACKGROUND_RATIO > /proc/sys/vm/dirty_background_ratio &
fi

if [ -n "$DIRTY_BYTES" ] && [ -w /proc/sys/vm/dirty_bytes ]; then
  echo $DIRTY_BYTES > /proc/sys/vm/dirty_bytes &
fi

if [ -n "$DIRTY_EXPIRE_CENTISECS" ] && [ -w /proc/sys/vm/dirty_expire_centisecs ]; then
  echo $DIRTY_EXPIRE_CENTISECS > /proc/sys/vm/dirty_expire_centisecs &
fi

if [ -n "$DIRTY_RATIO" ] && [ -w /proc/sys/vm/dirty_ratio ]; then
  echo $DIRTY_RATIO > /proc/sys/vm/dirty_ratio &
fi

if [ -n "$DIRTY_WRITEBACK_CENTISECS" ] && [ -w /proc/sys/vm/dirty_writeback_centisecs ]; then
  echo $DIRTY_WRITEBACK_CENTISECS > /proc/sys/vm/dirty_writeback_centisecs &
fi

if [ -n "$DIRTYTIME_EXPIRE_SECONDS" ] && [ -w /proc/sys/vm/dirtytime_expire_seconds ]; then
  echo $DIRTYTIME_EXPIRE_SECONDS > /proc/sys/vm/dirtytime_expire_seconds &
fi

if [ -n "$DROP_CACHES" ] && [ -w /proc/sys/vm/drop_caches ]; then
  echo $DROP_CACHES > /proc/sys/vm/drop_caches &
fi

if [ -n "$EXTFRAG_THRESHOLD" ] && [ -w /proc/sys/vm/extfrag_threshold ]; then
  echo $EXTFRAG_THRESHOLD > /proc/sys/vm/extfrag_threshold &
fi

if [ -n "$HUGEPAGES_TREAT_AS_MOVABLE" ] && [ -w /proc/sys/vm/hugepages_treat_as_movable ]; then
  echo $HUGEPAGES_TREAT_AS_MOVABLE > /proc/sys/vm/hugepages_treat_as_movable &
fi

if [ -n "$HUGETLB_SHM_GROUP" ] && [ -w /proc/sys/vm/hugetlb_shm_group ]; then
  echo $HUGETLB_SHM_GROUP > /proc/sys/vm/hugetlb_shm_group &
fi

if [ -n "$LAPTOP_MODE" ] && [ -w /proc/sys/vm/laptop_mode ]; then
  echo $LAPTOP_MODE > /proc/sys/vm/laptop_mode &
fi

if [ -n "$LEGACY_VA_LAYOUT" ] && [ -w /proc/sys/vm/legacy_va_layout ]; then
  echo $LEGACY_VA_LAYOUT > /proc/sys/vm/legacy_va_layout &
fi

if [ -n "$LOWMEM_RESERVE_RATIO" ] && [ -w /proc/sys/vm/lowmem_reserve_ratio ]; then
  echo $LOWMEM_RESERVE_RATIO > /proc/sys/vm/lowmem_reserve_ratio &
fi

if [ -n "$MAX_MAP_COUNT" ] && [ -w /proc/sys/vm/max_map_count ]; then
  echo $MAX_MAP_COUNT > /proc/sys/vm/max_map_count &
fi

if [ -n "$MEMORY_FAILURE_EARLY_KILL" ] && [ -w /proc/sys/vm/memory_failure_early_kill ]; then
  echo $MEMORY_FAILURE_EARLY_KILL > /proc/sys/vm/memory_failure_early_kill &
fi

if [ -n "$MEMORY_FAILURE_RECOVERY" ] && [ -w /proc/sys/vm/memory_failure_recovery ]; then
  echo $MEMORY_FAILURE_RECOVERY > /proc/sys/vm/memory_failure_recovery &
fi

if [ -n "$NR_HUGEPAGES" ] && [ -w /proc/sys/vm/nr_hugepages ]; then
  echo $NR_HUGEPAGES > /proc/sys/vm/nr_hugepages &
fi

if [ -n "$NR_HUGEPAGES_MEMPOLICY" ] && [ -w /proc/sys/vm/nr_hugepages_mempolicy ]; then
  echo $NR_HUGEPAGES_MEMPOLICY > /proc/sys/vm/nr_hugepages_mempolicy &
fi

if [ -n "$NR_OVERCOMMIT_HUGEPAGES" ] && [ -w /proc/sys/vm/nr_overcommit_hugepages ]; then
  echo $NR_OVERCOMMIT_HUGEPAGES > /proc/sys/vm/nr_overcommit_hugepages &
fi

if [ -n "$NUMA_ZONELIST_ORDER" ] && [ -w /proc/sys/vm/numa_zonelist_order ]; then
  echo $NUMA_ZONELIST_ORDER > /proc/sys/vm/numa_zonelist_order &
fi

if [ -n "$OOM_DUMP_TASKS" ] && [ -w /proc/sys/vm/oom_dump_tasks ]; then
  echo $OOM_DUMP_TASKS > /proc/sys/vm/oom_dump_tasks &
fi

if [ -n "$OOM_KILL_ALLOCATING_TASK" ] && [ -w /proc/sys/vm/oom_kill_allocating_task ]; then
  echo $OOM_KILL_ALLOCATING_TASK > /proc/sys/vm/oom_kill_allocating_task &
fi

if [ -n "$OVERCOMMIT_KBYTES" ] && [ -w /proc/sys/vm/overcommit_kbytes ]; then
  echo $OVERCOMMIT_KBYTES > /proc/sys/vm/overcommit_kbytes &
fi

if [ -n "$OVERCOMMIT_MEMORY" ] && [ -w /proc/sys/vm/overcommit_memory ]; then
  echo $OVERCOMMIT_MEMORY > /proc/sys/vm/overcommit_memory &
fi

if [ -n "$OVERCOMMIT_RATIO" ] && [ -w /proc/sys/vm/overcommit_ratio ]; then
  echo $OVERCOMMIT_RATIO > /proc/sys/vm/overcommit_ratio &
fi

if [ -n "$PAGE_CLUSTER" ] && [ -w /proc/sys/vm/page-cluster ]; then
  echo $PAGE_CLUSTER > /proc/sys/vm/page-cluster &
fi

if [ -n "$PANIC_ON_OOM" ] && [ -w /proc/sys/vm/panic_on_oom ]; then
  echo $PANIC_ON_OOM > /proc/sys/vm/panic_on_oom &
fi

if [ -n "$PERCPU_PAGELIST_FRACTION" ] && [ -w /proc/sys/vm/percpu_pagelist_fraction ]; then
  echo $PERCPU_PAGELIST_FRACTION > /proc/sys/vm/percpu_pagelist_fraction &
fi

if [ -n "$STAT_INTERVAL" ] && [ -w /proc/sys/vm/stat_interval ]; then
  echo $STAT_INTERVAL > /proc/sys/vm/stat_interval &
fi

if [ -n "$SWAPPINESS" ] && [ -w /proc/sys/vm/swappiness ]; then
  echo $SWAPPINESS > /proc/sys/vm/swappiness &
fi

if [ -n "$USER_RESERVE_KBYTES" ] && [ -w /proc/sys/vm/user_reserve_kbytes ]; then
  echo $USER_RESERVE_KBYTES > /proc/sys/vm/user_reserve_kbytes &
fi

if [ -n "$VFS_CACHE_PRESSURE" ] && [ -w /proc/sys/vm/vfs_cache_pressure ]; then
  echo $VFS_CACHE_PRESSURE > /proc/sys/vm/vfs_cache_pressure &
fi

if [ -n "$ZONE_RECLAIM_MODE" ] && [ -w /proc/sys/vm/zone_reclaim_mode ]; then
  echo $ZONE_RECLAIM_MODE > /proc/sys/vm/zone_reclaim_mode &
fi

if [ -n "$INODE_READAHEAD_BLKS" ] && [ -w /sys/fs/ext2/vda/inode_readahead_blks ]; then
  echo $INODE_READAHEAD_BLKS > /sys/fs/ext2/vda/inode_readahead_blks &
fi

if [ -n "$MB_MAX_TO_SCAN" ] && [ -w /sys/fs/ext2/vda/mb_max_to_scan ]; then
  echo $MB_MAX_TO_SCAN > /sys/fs/ext2/vda/mb_max_to_scan &
fi

if [ -n "$MSG_RATELIMIT_BURST" ] && [ -w /sys/fs/ext2/vda/msg_ratelimit_burst ]; then
  echo $MSG_RATELIMIT_BURST > /sys/fs/ext2/vda/msg_ratelimit_burst &
fi

if [ -n "$MB_STREAM_REQ" ] && [ -w /sys/fs/ext2/vda/mb_stream_req ]; then
  echo $MB_STREAM_REQ > /sys/fs/ext2/vda/mb_stream_req &
fi

if [ -n "$MB_MIN_TO_SCAN" ] && [ -w /sys/fs/ext2/vda/mb_min_to_scan ]; then
  echo $MB_MIN_TO_SCAN > /sys/fs/ext2/vda/mb_min_to_scan &
fi

if [ -n "$MB_STATS" ] && [ -w /sys/fs/ext2/vda/mb_stats ]; then
  echo $MB_STATS > /sys/fs/ext2/vda/mb_stats &
fi

if [ -n "$ERR_RATELIMIT_BURST" ] && [ -w /sys/fs/ext2/vda/err_ratelimit_burst ]; then
  echo $ERR_RATELIMIT_BURST > /sys/fs/ext2/vda/err_ratelimit_burst &
fi

if [ -n "$MB_GROUP_PREALLOC" ] && [ -w /sys/fs/ext2/vda/mb_group_prealloc ]; then
  echo $MB_GROUP_PREALLOC > /sys/fs/ext2/vda/mb_group_prealloc &
fi

if [ -n "$INODE_GOAL" ] && [ -w /sys/fs/ext2/vda/inode_goal ]; then
  echo $INODE_GOAL > /sys/fs/ext2/vda/inode_goal &
fi

if [ -n "$RESERVED_CLUSTERS" ] && [ -w /sys/fs/ext2/vda/reserved_clusters ]; then
  echo $RESERVED_CLUSTERS > /sys/fs/ext2/vda/reserved_clusters &
fi

if [ -n "$EXTENT_MAX_ZEROOUT_KB" ] && [ -w /sys/fs/ext2/vda/extent_max_zeroout_kb ]; then
  echo $EXTENT_MAX_ZEROOUT_KB > /sys/fs/ext2/vda/extent_max_zeroout_kb &
fi

if [ -n "$ERR_RATELIMIT_INTERVAL_MS" ] && [ -w /sys/fs/ext2/vda/err_ratelimit_interval_ms ]; then
  echo $ERR_RATELIMIT_INTERVAL_MS > /sys/fs/ext2/vda/err_ratelimit_interval_ms &
fi

if [ -n "$WARNING_RATELIMIT_BURST" ] && [ -w /sys/fs/ext2/vda/warning_ratelimit_burst ]; then
  echo $WARNING_RATELIMIT_BURST > /sys/fs/ext2/vda/warning_ratelimit_burst &
fi

if [ -n "$WARNING_RATELIMIT_INTERVAL_MS" ] && [ -w /sys/fs/ext2/vda/warning_ratelimit_interval_ms ]; then
  echo $WARNING_RATELIMIT_INTERVAL_MS > /sys/fs/ext2/vda/warning_ratelimit_interval_ms &
fi

if [ -n "$MB_ORDER2_REQ" ] && [ -w /sys/fs/ext2/vda/mb_order2_req ]; then
  echo $MB_ORDER2_REQ > /sys/fs/ext2/vda/mb_order2_req &
fi

if [ -n "$MSG_RATELIMIT_INTERVAL_MS" ] && [ -w /sys/fs/ext2/vda/msg_ratelimit_interval_ms ]; then
  echo $MSG_RATELIMIT_INTERVAL_MS > /sys/fs/ext2/vda/msg_ratelimit_interval_ms &
fi

if [ -n "$NR_OVERCOMMIT_HUGEPAGES" ] && [ -w /sys/kernel/mm/hugepages/hugepages-2048kB/nr_overcommit_hugepages ]; then
  echo $NR_OVERCOMMIT_HUGEPAGES > /sys/kernel/mm/hugepages/hugepages-2048kB/nr_overcommit_hugepages &
fi

if [ -n "$NR_HUGEPAGES" ] && [ -w /sys/kernel/mm/hugepages/hugepages-2048kB/nr_hugepages ]; then
  echo $NR_HUGEPAGES > /sys/kernel/mm/hugepages/hugepages-2048kB/nr_hugepages &
fi

if [ -n "$NR_HUGEPAGES_MEMPOLICY" ] && [ -w /sys/kernel/mm/hugepages/hugepages-2048kB/nr_hugepages_mempolicy ]; then
  echo $NR_HUGEPAGES_MEMPOLICY > /sys/kernel/mm/hugepages/hugepages-2048kB/nr_hugepages_mempolicy &
fi

if [ -n "$ENABLED" ] && [ -w /sys/kernel/mm/transparent_hugepage/enabled ]; then
  echo $ENABLED > /sys/kernel/mm/transparent_hugepage/enabled &
fi

if [ -n "$USE_ZERO_PAGE" ] && [ -w /sys/kernel/mm/transparent_hugepage/use_zero_page ]; then
  echo $USE_ZERO_PAGE > /sys/kernel/mm/transparent_hugepage/use_zero_page &
fi

if [ -n "$DEFRAG" ] && [ -w /sys/kernel/mm/transparent_hugepage/defrag ]; then
  echo $DEFRAG > /sys/kernel/mm/transparent_hugepage/defrag &
fi

if [ -n "$MTU" ] && [ -w /sys/class/net/eth0/mtu ]; then
  echo $MTU > /sys/class/net/eth0/mtu &
fi

if [ -n "$NETDEV_GROUP" ] && [ -w /sys/class/net/eth0/netdev_group ]; then
  echo $NETDEV_GROUP > /sys/class/net/eth0/netdev_group &
fi

if [ -n "$FLAGS" ] && [ -w /sys/class/net/eth0/flags ]; then
  echo $FLAGS > /sys/class/net/eth0/flags &
fi

if [ -n "$GRO_FLUSH_TIMEOUT" ] && [ -w /sys/class/net/eth0/gro_flush_timeout ]; then
  echo $GRO_FLUSH_TIMEOUT > /sys/class/net/eth0/gro_flush_timeout &
fi

if [ -n "$TX_QUEUE_LEN" ] && [ -w /sys/class/net/eth0/tx_queue_len ]; then
  echo $TX_QUEUE_LEN > /sys/class/net/eth0/tx_queue_len &
fi

sleep 1

REDIS_MAX_MEMORY=$(echo $WAYFINDER_DOMAIN_MAX_MEM | sed 's/iB$//')

sed -i \
    -e "s/%REDIS_MAX_MEMORY%/$REDIS_MAX_MEMORY/g" \
    -e "s/%REDIS_TCP_BACKLOG%/$REDIS_TCP_BACKLOG/g" \
    -e "s/%REDIS_TIMEOUT%/$REDIS_TIMEOUT/g" \
    -e "s/%REDIS_KEEPALIVE%/$REDIS_KEEPALIVE/g" \
    -e "s/%REDIS_DATABASES%/$REDIS_DATABASES/g" \
    -e "s/%REDIS_SAVE_TIME%/$REDIS_SAVE_TIME/g" \
    -e "s/%REDIS_SAVE_RECORDS%/$REDIS_SAVE_RECORDS/g" \
    -e "s/%REDIS_RBCOMRESSION%/$REDIS_RBCOMRESSION/g" \
    -e "s/%REDIS_RDBCHECKSUM%/$REDIS_RDBCHECKSUM/g" \
    -e "s/%REDIS_MAX_MEMORY_POLICY%/$REDIS_MAX_MEMORY_POLICY/g" \
    -e "s/%REDIS_MAX_MEMORY_SAMPLE%/$REDIS_MAX_MEMORY_SAMPLE/g" \
    -e "s/%REDIS_APPENDONLY%/$REDIS_APPENDONLY/g" \
    -e "s/%REDIS_AOF%/$REDIS_AOF/g" \
    -e "s/%REDIS_AOF_REWRITE%/$REDIS_AOF_REWRITE/g" \
    -e "s/%REDIS_AOF_PERCENT%/$REDIS_AOF_PERCENT/g" \
    -e "s/%REDIS_AOF_SIZE%/$REDIS_AOF_SIZE/g" \
    -e "s/%REDIS_SLOWLOG_LEN%/$REDIS_SLOWLOG_LEN/g" \
    -e "s/%REDIS_HASH_ENTRIES%/$REDIS_HASH_ENTRIES/g" \
    -e "s/%REDIS_HASH_VALUES%/$REDIS_HASH_VALUES/g" \
    -e "s/%REDIS_ZSET_ENTRIES%/$REDIS_ZSET_ENTRIES/g" \
    -e "s/%REDIS_ZSET_VALUES%/$REDIS_ZSET_VALUES/g" \
    -e "s/%REDIS_LISTPACK_SIZE%/$REDIS_LISTPACK_SIZE/g" \
    -e "s/%REDIS_COMPRESS%/$REDIS_COMPRESS/g" \
    -e "s/%REDIS_INTSET%/$REDIS_INTSET/g" \
    -e "s/%REDIS_REHASHING%/$REDIS_REHASHING/g" \
    -e "s/%REDIS_DEFRAG%/$REDIS_DEFRAG/g" \
    -e "s/%REDIS_JEMALOC_THREAD%/$REDIS_JEMALOC_THREAD/g" \
    /redis.conf

# dirty, but not sure how to do this differently
redis-server /redis.conf
