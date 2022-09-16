#!/bin/bash

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

NODE_OPTIONS=()

if [ -n "$ENABLE_FIPS" ] && [ $ENABLE_FIPS == "y" ]; then
  NODE_OPTIONS+=("--enable-fips")
fi

if [ -n "$NO_DEPRECATION" ] && [ $NO_DEPRECATION == "y" ]; then
  NODE_OPTIONS+=("--no-deprecation")
fi

if [ -n "$NO_WARNINGS" ] && [ $NO_WARNINGS == "y" ]; then
  NODE_OPTIONS+=("--no-warnings")
fi

if [ -n "$TRACE_SYNC_IO" ] && [ $TRACE_SYNC_IO == "y" ]; then
  NODE_OPTIONS+=("--trace-sync-io")
fi

if [ -n "$TRACE_DEPRECATION" ] && [ $TRACE_DEPRECATION == "y" ]; then
  NODE_OPTIONS+=("--trace-deprecation")
fi

if [ -n "$TRACE_WARNINGS" ] && [ $TRACE_WARNINGS == "y" ]; then
  NODE_OPTIONS+=("--trace-warnings")
fi

if [ -n "$TRACK_HEAP_OBJECTS" ] && [ $TRACK_HEAP_OBJECTS == "y" ]; then
  NODE_OPTIONS+=("--track-heap-objects")
fi

if [ -n "$ZERO_FILL_BUFFERS" ] && [ $ZERO_FILL_BUFFERS == "y" ]; then
  NODE_OPTIONS+=("--zero-fill-buffers")
fi

if [ -n "$EXPERIMENTAL_EXTRAS" ] && [ $EXPERIMENTAL_EXTRAS == "y" ]; then
  NODE_OPTIONS+=("--experimental_extras")
else
  NODE_OPTIONS+=("--no-experimental_extras")
fi

if [ -n "$USE_STRICT" ] && [ $USE_STRICT == "y" ]; then
  NODE_OPTIONS+=("--use_strict")
else
  NODE_OPTIONS+=("--no-use_strict")
fi

if [ -n "$ES_STAGING" ] && [ $ES_STAGING == "y" ]; then
  NODE_OPTIONS+=("--es_staging")
else
  NODE_OPTIONS+=("--no-es_staging")
fi

if [ -n "$HARMONY" ] && [ $HARMONY == "y" ]; then
  NODE_OPTIONS+=("--harmony")
else
  NODE_OPTIONS+=("--no-harmony")
fi

if [ -n "$HARMONY_SHIPPING" ] && [ $HARMONY_SHIPPING == "y" ]; then
  NODE_OPTIONS+=("--harmony_shipping")
else
  NODE_OPTIONS+=("--no-harmony_shipping")
fi

if [ -n "$HARMONY_ARRAY_PROTOTYPE_VALUES" ] && [ $HARMONY_ARRAY_PROTOTYPE_VALUES == "y" ]; then
  NODE_OPTIONS+=("--harmony_array_prototype_values")
else
  NODE_OPTIONS+=("--no-harmony_array_prototype_values")
fi

if [ -n "$HARMONY_FUNCTION_SENT" ] && [ $HARMONY_FUNCTION_SENT == "y" ]; then
  NODE_OPTIONS+=("--harmony_function_sent")
else
  NODE_OPTIONS+=("--no-harmony_function_sent")
fi

if [ -n "$HARMONY_SHAREDARRAYBUFFER" ] && [ $HARMONY_SHAREDARRAYBUFFER == "y" ]; then
  NODE_OPTIONS+=("--harmony_sharedarraybuffer")
else
  NODE_OPTIONS+=("--no-harmony_sharedarraybuffer")
fi

if [ -n "$HARMONY_DO_EXPRESSIONS" ] && [ $HARMONY_DO_EXPRESSIONS == "y" ]; then
  NODE_OPTIONS+=("--harmony_do_expressions")
else
  NODE_OPTIONS+=("--no-harmony_do_expressions")
fi

if [ -n "$HARMONY_REGEXP_NAMED_CAPTURES" ] && [ $HARMONY_REGEXP_NAMED_CAPTURES == "y" ]; then
  NODE_OPTIONS+=("--harmony_regexp_named_captures")
else
  NODE_OPTIONS+=("--no-harmony_regexp_named_captures")
fi

if [ -n "$HARMONY_REGEXP_PROPERTY" ] && [ $HARMONY_REGEXP_PROPERTY == "y" ]; then
  NODE_OPTIONS+=("--harmony_regexp_property")
else
  NODE_OPTIONS+=("--no-harmony_regexp_property")
fi

if [ -n "$HARMONY_FUNCTION_TOSTRING" ] && [ $HARMONY_FUNCTION_TOSTRING == "y" ]; then
  NODE_OPTIONS+=("--harmony_function_tostring")
else
  NODE_OPTIONS+=("--no-harmony_function_tostring")
fi

if [ -n "$HARMONY_CLASS_FIELDS" ] && [ $HARMONY_CLASS_FIELDS == "y" ]; then
  NODE_OPTIONS+=("--harmony_class_fields")
else
  NODE_OPTIONS+=("--no-harmony_class_fields")
fi

if [ -n "$HARMONY_ASYNC_ITERATION" ] && [ $HARMONY_ASYNC_ITERATION == "y" ]; then
  NODE_OPTIONS+=("--harmony_async_iteration")
else
  NODE_OPTIONS+=("--no-harmony_async_iteration")
fi

if [ -n "$HARMONY_DYNAMIC_IMPORT" ] && [ $HARMONY_DYNAMIC_IMPORT == "y" ]; then
  NODE_OPTIONS+=("--harmony_dynamic_import")
else
  NODE_OPTIONS+=("--no-harmony_dynamic_import")
fi

if [ -n "$HARMONY_PROMISE_FINALLY" ] && [ $HARMONY_PROMISE_FINALLY == "y" ]; then
  NODE_OPTIONS+=("--harmony_promise_finally")
else
  NODE_OPTIONS+=("--no-harmony_promise_finally")
fi

if [ -n "$HARMONY_REGEXP_LOOKBEHIND" ] && [ $HARMONY_REGEXP_LOOKBEHIND == "y" ]; then
  NODE_OPTIONS+=("--harmony_regexp_lookbehind")
else
  NODE_OPTIONS+=("--no-harmony_regexp_lookbehind")
fi

if [ -n "$HARMONY_RESTRICTIVE_GENERATORS" ] && [ $HARMONY_RESTRICTIVE_GENERATORS == "y" ]; then
  NODE_OPTIONS+=("--harmony_restrictive_generators")
else
  NODE_OPTIONS+=("--no-harmony_restrictive_generators")
fi

if [ -n "$HARMONY_OBJECT_REST_SPREAD" ] && [ $HARMONY_OBJECT_REST_SPREAD == "y" ]; then
  NODE_OPTIONS+=("--harmony_object_rest_spread")
else
  NODE_OPTIONS+=("--no-harmony_object_rest_spread")
fi

if [ -n "$HARMONY_TEMPLATE_ESCAPES" ] && [ $HARMONY_TEMPLATE_ESCAPES == "y" ]; then
  NODE_OPTIONS+=("--harmony_template_escapes")
else
  NODE_OPTIONS+=("--no-harmony_template_escapes")
fi

if [ -n "$FUTURE" ] && [ $FUTURE == "y" ]; then
  NODE_OPTIONS+=("--future")
else
  NODE_OPTIONS+=("--no-future")
fi

if [ -n "$ALLOCATION_SITE_PRETENURING" ] && [ $ALLOCATION_SITE_PRETENURING == "y" ]; then
  NODE_OPTIONS+=("--allocation_site_pretenuring")
else
  NODE_OPTIONS+=("--no-allocation_site_pretenuring")
fi

if [ -n "$PAGE_PROMOTION" ] && [ $PAGE_PROMOTION == "y" ]; then
  NODE_OPTIONS+=("--page_promotion")
else
  NODE_OPTIONS+=("--no-page_promotion")
fi

if [ -n "$TRACE_PRETENURING" ] && [ $TRACE_PRETENURING == "y" ]; then
  NODE_OPTIONS+=("--trace_pretenuring")
else
  NODE_OPTIONS+=("--no-trace_pretenuring")
fi

if [ -n "$TRACE_PRETENURING_STATISTICS" ] && [ $TRACE_PRETENURING_STATISTICS == "y" ]; then
  NODE_OPTIONS+=("--trace_pretenuring_statistics")
else
  NODE_OPTIONS+=("--no-trace_pretenuring_statistics")
fi

if [ -n "$TRACK_FIELDS" ] && [ $TRACK_FIELDS == "y" ]; then
  NODE_OPTIONS+=("--track_fields")
else
  NODE_OPTIONS+=("--no-track_fields")
fi

if [ -n "$TRACK_DOUBLE_FIELDS" ] && [ $TRACK_DOUBLE_FIELDS == "y" ]; then
  NODE_OPTIONS+=("--track_double_fields")
else
  NODE_OPTIONS+=("--no-track_double_fields")
fi

if [ -n "$TRACK_HEAP_OBJECT_FIELDS" ] && [ $TRACK_HEAP_OBJECT_FIELDS == "y" ]; then
  NODE_OPTIONS+=("--track_heap_object_fields")
else
  NODE_OPTIONS+=("--no-track_heap_object_fields")
fi

if [ -n "$TRACK_COMPUTED_FIELDS" ] && [ $TRACK_COMPUTED_FIELDS == "y" ]; then
  NODE_OPTIONS+=("--track_computed_fields")
else
  NODE_OPTIONS+=("--no-track_computed_fields")
fi

if [ -n "$TRACK_FIELD_TYPES" ] && [ $TRACK_FIELD_TYPES == "y" ]; then
  NODE_OPTIONS+=("--track_field_types")
else
  NODE_OPTIONS+=("--no-track_field_types")
fi

if [ -n "$TYPE_PROFILE" ] && [ $TYPE_PROFILE == "y" ]; then
  NODE_OPTIONS+=("--type_profile")
else
  NODE_OPTIONS+=("--no-type_profile")
fi

if [ -n "$OPTIMIZE_FOR_SIZE" ] && [ $OPTIMIZE_FOR_SIZE == "y" ]; then
  NODE_OPTIONS+=("--optimize_for_size")
else
  NODE_OPTIONS+=("--no-optimize_for_size")
fi

if [ -n "$UNBOX_DOUBLE_ARRAYS" ] && [ $UNBOX_DOUBLE_ARRAYS == "y" ]; then
  NODE_OPTIONS+=("--unbox_double_arrays")
else
  NODE_OPTIONS+=("--no-unbox_double_arrays")
fi

if [ -n "$STRING_SLICES" ] && [ $STRING_SLICES == "y" ]; then
  NODE_OPTIONS+=("--string_slices")
else
  NODE_OPTIONS+=("--no-string_slices")
fi

if [ -n "$IGNITION_REO" ] && [ $IGNITION_REO == "y" ]; then
  NODE_OPTIONS+=("--ignition_reo")
else
  NODE_OPTIONS+=("--no-ignition_reo")
fi

if [ -n "$IGNITION_FILTER_EXPRESSION_POSITIONS" ] && [ $IGNITION_FILTER_EXPRESSION_POSITIONS == "y" ]; then
  NODE_OPTIONS+=("--ignition_filter_expression_positions")
else
  NODE_OPTIONS+=("--no-ignition_filter_expression_positions")
fi

if [ -n "$PRINT_BYTECODE" ] && [ $PRINT_BYTECODE == "y" ]; then
  NODE_OPTIONS+=("--print_bytecode")
else
  NODE_OPTIONS+=("--no-print_bytecode")
fi

if [ -n "$TRACE_IGNITION_CODEGEN" ] && [ $TRACE_IGNITION_CODEGEN == "y" ]; then
  NODE_OPTIONS+=("--trace_ignition_codegen")
else
  NODE_OPTIONS+=("--no-trace_ignition_codegen")
fi

if [ -n "$TRACE_IGNITION_DISPATCHES" ] && [ $TRACE_IGNITION_DISPATCHES == "y" ]; then
  NODE_OPTIONS+=("--trace_ignition_dispatches")
else
  NODE_OPTIONS+=("--no-trace_ignition_dispatches")
fi

if [ -n "$FAST_MATH" ] && [ $FAST_MATH == "y" ]; then
  NODE_OPTIONS+=("--fast_math")
else
  NODE_OPTIONS+=("--no-fast_math")
fi

if [ -n "$TRACE_ENVIRONMENT_LIVENESS" ] && [ $TRACE_ENVIRONMENT_LIVENESS == "y" ]; then
  NODE_OPTIONS+=("--trace_environment_liveness")
else
  NODE_OPTIONS+=("--no-trace_environment_liveness")
fi

if [ -n "$TRACE_STORE_ELIMINATION" ] && [ $TRACE_STORE_ELIMINATION == "y" ]; then
  NODE_OPTIONS+=("--trace_store_elimination")
else
  NODE_OPTIONS+=("--no-trace_store_elimination")
fi

if [ -n "$TRACE_ALLOC" ] && [ $TRACE_ALLOC == "y" ]; then
  NODE_OPTIONS+=("--trace_alloc")
else
  NODE_OPTIONS+=("--no-trace_alloc")
fi

if [ -n "$TRACE_ALL_USES" ] && [ $TRACE_ALL_USES == "y" ]; then
  NODE_OPTIONS+=("--trace_all_uses")
else
  NODE_OPTIONS+=("--no-trace_all_uses")
fi

if [ -n "$TRACE_REPRESENTATION" ] && [ $TRACE_REPRESENTATION == "y" ]; then
  NODE_OPTIONS+=("--trace_representation")
else
  NODE_OPTIONS+=("--no-trace_representation")
fi

if [ -n "$TRACE_TRACK_ALLOCATION_SITES" ] && [ $TRACE_TRACK_ALLOCATION_SITES == "y" ]; then
  NODE_OPTIONS+=("--trace_track_allocation_sites")
else
  NODE_OPTIONS+=("--no-trace_track_allocation_sites")
fi

if [ -n "$TRACE_MIGRATION" ] && [ $TRACE_MIGRATION == "y" ]; then
  NODE_OPTIONS+=("--trace_migration")
else
  NODE_OPTIONS+=("--no-trace_migration")
fi

if [ -n "$TRACE_GENERALIZATION" ] && [ $TRACE_GENERALIZATION == "y" ]; then
  NODE_OPTIONS+=("--trace_generalization")
else
  NODE_OPTIONS+=("--no-trace_generalization")
fi

if [ -n "$PRINT_DEOPT_STRESS" ] && [ $PRINT_DEOPT_STRESS == "y" ]; then
  NODE_OPTIONS+=("--print_deopt_stress")
else
  NODE_OPTIONS+=("--no-print_deopt_stress")
fi

if [ -n "$POLYMORPHIC_INLINING" ] && [ $POLYMORPHIC_INLINING == "y" ]; then
  NODE_OPTIONS+=("--polymorphic_inlining")
else
  NODE_OPTIONS+=("--no-polymorphic_inlining")
fi

if [ -n "$USE_OSR" ] && [ $USE_OSR == "y" ]; then
  NODE_OPTIONS+=("--use_osr")
else
  NODE_OPTIONS+=("--no-use_osr")
fi

if [ -n "$ANALYZE_ENVIRONMENT_LIVENESS" ] && [ $ANALYZE_ENVIRONMENT_LIVENESS == "y" ]; then
  NODE_OPTIONS+=("--analyze_environment_liveness")
else
  NODE_OPTIONS+=("--no-analyze_environment_liveness")
fi

if [ -n "$TRACE_OSR" ] && [ $TRACE_OSR == "y" ]; then
  NODE_OPTIONS+=("--trace_osr")
else
  NODE_OPTIONS+=("--no-trace_osr")
fi

if [ -n "$INLINE_ACCESSORS" ] && [ $INLINE_ACCESSORS == "y" ]; then
  NODE_OPTIONS+=("--inline_accessors")
else
  NODE_OPTIONS+=("--no-inline_accessors")
fi

if [ -n "$INLINE_INTO_TRY" ] && [ $INLINE_INTO_TRY == "y" ]; then
  NODE_OPTIONS+=("--inline_into_try")
else
  NODE_OPTIONS+=("--no-inline_into_try")
fi

if [ -n "$CONCURRENT_RECOMPILATION" ] && [ $CONCURRENT_RECOMPILATION == "y" ]; then
  NODE_OPTIONS+=("--concurrent_recompilation")
else
  NODE_OPTIONS+=("--no-concurrent_recompilation")
fi

if [ -n "$TRACE_CONCURRENT_RECOMPILATION" ] && [ $TRACE_CONCURRENT_RECOMPILATION == "y" ]; then
  NODE_OPTIONS+=("--trace_concurrent_recompilation")
else
  NODE_OPTIONS+=("--no-trace_concurrent_recompilation")
fi

if [ -n "$BLOCK_CONCURRENT_RECOMPILATION" ] && [ $BLOCK_CONCURRENT_RECOMPILATION == "y" ]; then
  NODE_OPTIONS+=("--block_concurrent_recompilation")
else
  NODE_OPTIONS+=("--no-block_concurrent_recompilation")
fi

if [ -n "$TURBO_SP_FRAME_ACCESS" ] && [ $TURBO_SP_FRAME_ACCESS == "y" ]; then
  NODE_OPTIONS+=("--turbo_sp_frame_access")
else
  NODE_OPTIONS+=("--no-turbo_sp_frame_access")
fi

if [ -n "$TURBO_PREPROCESS_RANGES" ] && [ $TURBO_PREPROCESS_RANGES == "y" ]; then
  NODE_OPTIONS+=("--turbo_preprocess_ranges")
else
  NODE_OPTIONS+=("--no-turbo_preprocess_ranges")
fi

if [ -n "$TRACE_TURBO" ] && [ $TRACE_TURBO == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo")
else
  NODE_OPTIONS+=("--no-trace_turbo")
fi

if [ -n "$TRACE_TURBO_GRAPH" ] && [ $TRACE_TURBO_GRAPH == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_graph")
else
  NODE_OPTIONS+=("--no-trace_turbo_graph")
fi

if [ -n "$TRACE_TURBO_TYPES" ] && [ $TRACE_TURBO_TYPES == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_types")
else
  NODE_OPTIONS+=("--no-trace_turbo_types")
fi

if [ -n "$TRACE_TURBO_SCHEDULER" ] && [ $TRACE_TURBO_SCHEDULER == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_scheduler")
else
  NODE_OPTIONS+=("--no-trace_turbo_scheduler")
fi

if [ -n "$TRACE_TURBO_REDUCTION" ] && [ $TRACE_TURBO_REDUCTION == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_reduction")
else
  NODE_OPTIONS+=("--no-trace_turbo_reduction")
fi

if [ -n "$TRACE_TURBO_TRIMMING" ] && [ $TRACE_TURBO_TRIMMING == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_trimming")
else
  NODE_OPTIONS+=("--no-trace_turbo_trimming")
fi

if [ -n "$TRACE_TURBO_JT" ] && [ $TRACE_TURBO_JT == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_jt")
else
  NODE_OPTIONS+=("--no-trace_turbo_jt")
fi

if [ -n "$TRACE_TURBO_CEQ" ] && [ $TRACE_TURBO_CEQ == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_ceq")
else
  NODE_OPTIONS+=("--no-trace_turbo_ceq")
fi

if [ -n "$TRACE_TURBO_LOOP" ] && [ $TRACE_TURBO_LOOP == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_loop")
else
  NODE_OPTIONS+=("--no-trace_turbo_loop")
fi

if [ -n "$TURBO_VERIFY" ] && [ $TURBO_VERIFY == "y" ]; then
  NODE_OPTIONS+=("--turbo_verify")
else
  NODE_OPTIONS+=("--no-turbo_verify")
fi

if [ -n "$TRACE_VERIFY_CSA" ] && [ $TRACE_VERIFY_CSA == "y" ]; then
  NODE_OPTIONS+=("--trace_verify_csa")
else
  NODE_OPTIONS+=("--no-trace_verify_csa")
fi

if [ -n "$TURBO_STATS" ] && [ $TURBO_STATS == "y" ]; then
  NODE_OPTIONS+=("--turbo_stats")
else
  NODE_OPTIONS+=("--no-turbo_stats")
fi

if [ -n "$TURBO_STATS_NVP" ] && [ $TURBO_STATS_NVP == "y" ]; then
  NODE_OPTIONS+=("--turbo_stats_nvp")
else
  NODE_OPTIONS+=("--no-turbo_stats_nvp")
fi

if [ -n "$TURBO_SPLITTING" ] && [ $TURBO_SPLITTING == "y" ]; then
  NODE_OPTIONS+=("--turbo_splitting")
else
  NODE_OPTIONS+=("--no-turbo_splitting")
fi

if [ -n "$FUNCTION_CONTEXT_SPECIALIZATION" ] && [ $FUNCTION_CONTEXT_SPECIALIZATION == "y" ]; then
  NODE_OPTIONS+=("--function_context_specialization")
else
  NODE_OPTIONS+=("--no-function_context_specialization")
fi

if [ -n "$TURBO_INLINING" ] && [ $TURBO_INLINING == "y" ]; then
  NODE_OPTIONS+=("--turbo_inlining")
else
  NODE_OPTIONS+=("--no-turbo_inlining")
fi

if [ -n "$TRACE_TURBO_INLINING" ] && [ $TRACE_TURBO_INLINING == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_inlining")
else
  NODE_OPTIONS+=("--no-trace_turbo_inlining")
fi

if [ -n "$TURBO_LOAD_ELIMINATION" ] && [ $TURBO_LOAD_ELIMINATION == "y" ]; then
  NODE_OPTIONS+=("--turbo_load_elimination")
else
  NODE_OPTIONS+=("--no-turbo_load_elimination")
fi

if [ -n "$TRACE_TURBO_LOAD_ELIMINATION" ] && [ $TRACE_TURBO_LOAD_ELIMINATION == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_load_elimination")
else
  NODE_OPTIONS+=("--no-trace_turbo_load_elimination")
fi

if [ -n "$TURBO_PROFILING" ] && [ $TURBO_PROFILING == "y" ]; then
  NODE_OPTIONS+=("--turbo_profiling")
else
  NODE_OPTIONS+=("--no-turbo_profiling")
fi

if [ -n "$TURBO_VERIFY_ALLOCATION" ] && [ $TURBO_VERIFY_ALLOCATION == "y" ]; then
  NODE_OPTIONS+=("--turbo_verify_allocation")
else
  NODE_OPTIONS+=("--no-turbo_verify_allocation")
fi

if [ -n "$TURBO_MOVE_OPTIMIZATION" ] && [ $TURBO_MOVE_OPTIMIZATION == "y" ]; then
  NODE_OPTIONS+=("--turbo_move_optimization")
else
  NODE_OPTIONS+=("--no-turbo_move_optimization")
fi

if [ -n "$TURBO_JT" ] && [ $TURBO_JT == "y" ]; then
  NODE_OPTIONS+=("--turbo_jt")
else
  NODE_OPTIONS+=("--no-turbo_jt")
fi

if [ -n "$TURBO_LOOP_PEELING" ] && [ $TURBO_LOOP_PEELING == "y" ]; then
  NODE_OPTIONS+=("--turbo_loop_peeling")
else
  NODE_OPTIONS+=("--no-turbo_loop_peeling")
fi

if [ -n "$TURBO_LOOP_VARIABLE" ] && [ $TURBO_LOOP_VARIABLE == "y" ]; then
  NODE_OPTIONS+=("--turbo_loop_variable")
else
  NODE_OPTIONS+=("--no-turbo_loop_variable")
fi

if [ -n "$TURBO_CF_OPTIMIZATION" ] && [ $TURBO_CF_OPTIMIZATION == "y" ]; then
  NODE_OPTIONS+=("--turbo_cf_optimization")
else
  NODE_OPTIONS+=("--no-turbo_cf_optimization")
fi

if [ -n "$TURBO_FRAME_ELISION" ] && [ $TURBO_FRAME_ELISION == "y" ]; then
  NODE_OPTIONS+=("--turbo_frame_elision")
else
  NODE_OPTIONS+=("--no-turbo_frame_elision")
fi

if [ -n "$TURBO_ESCAPE" ] && [ $TURBO_ESCAPE == "y" ]; then
  NODE_OPTIONS+=("--turbo_escape")
else
  NODE_OPTIONS+=("--no-turbo_escape")
fi

if [ -n "$TURBO_INSTRUCTION_SCHEDULING" ] && [ $TURBO_INSTRUCTION_SCHEDULING == "y" ]; then
  NODE_OPTIONS+=("--turbo_instruction_scheduling")
else
  NODE_OPTIONS+=("--no-turbo_instruction_scheduling")
fi

if [ -n "$TURBO_STRESS_INSTRUCTION_SCHEDULING" ] && [ $TURBO_STRESS_INSTRUCTION_SCHEDULING == "y" ]; then
  NODE_OPTIONS+=("--turbo_stress_instruction_scheduling")
else
  NODE_OPTIONS+=("--no-turbo_stress_instruction_scheduling")
fi

if [ -n "$TURBO_STORE_ELIMINATION" ] && [ $TURBO_STORE_ELIMINATION == "y" ]; then
  NODE_OPTIONS+=("--turbo_store_elimination")
else
  NODE_OPTIONS+=("--no-turbo_store_elimination")
fi

if [ -n "$MINIMAL" ] && [ $MINIMAL == "y" ]; then
  NODE_OPTIONS+=("--minimal")
else
  NODE_OPTIONS+=("--no-minimal")
fi

if [ -n "$EXPOSE_WASM" ] && [ $EXPOSE_WASM == "y" ]; then
  NODE_OPTIONS+=("--expose_wasm")
else
  NODE_OPTIONS+=("--no-expose_wasm")
fi

if [ -n "$ASSUME_ASMJS_ORIGIN" ] && [ $ASSUME_ASMJS_ORIGIN == "y" ]; then
  NODE_OPTIONS+=("--assume_asmjs_origin")
else
  NODE_OPTIONS+=("--no-assume_asmjs_origin")
fi

if [ -n "$WASM_DISABLE_STRUCTURED_CLONING" ] && [ $WASM_DISABLE_STRUCTURED_CLONING == "y" ]; then
  NODE_OPTIONS+=("--wasm_disable_structured_cloning")
else
  NODE_OPTIONS+=("--no-wasm_disable_structured_cloning")
fi

if [ -n "$TRACE_WASM_DECODER" ] && [ $TRACE_WASM_DECODER == "y" ]; then
  NODE_OPTIONS+=("--trace_wasm_decoder")
else
  NODE_OPTIONS+=("--no-trace_wasm_decoder")
fi

if [ -n "$TRACE_WASM_DECODE_TIME" ] && [ $TRACE_WASM_DECODE_TIME == "y" ]; then
  NODE_OPTIONS+=("--trace_wasm_decode_time")
else
  NODE_OPTIONS+=("--no-trace_wasm_decode_time")
fi

if [ -n "$TRACE_WASM_COMPILER" ] && [ $TRACE_WASM_COMPILER == "y" ]; then
  NODE_OPTIONS+=("--trace_wasm_compiler")
else
  NODE_OPTIONS+=("--no-trace_wasm_compiler")
fi

if [ -n "$TRACE_WASM_INTERPRETER" ] && [ $TRACE_WASM_INTERPRETER == "y" ]; then
  NODE_OPTIONS+=("--trace_wasm_interpreter")
else
  NODE_OPTIONS+=("--no-trace_wasm_interpreter")
fi

if [ -n "$WASM_BREAK_ON_DECODER_ERROR" ] && [ $WASM_BREAK_ON_DECODER_ERROR == "y" ]; then
  NODE_OPTIONS+=("--wasm_break_on_decoder_error")
else
  NODE_OPTIONS+=("--no-wasm_break_on_decoder_error")
fi

if [ -n "$VALIDATE_ASM" ] && [ $VALIDATE_ASM == "y" ]; then
  NODE_OPTIONS+=("--validate_asm")
else
  NODE_OPTIONS+=("--no-validate_asm")
fi

if [ -n "$SUPPRESS_ASM_MESSAGES" ] && [ $SUPPRESS_ASM_MESSAGES == "y" ]; then
  NODE_OPTIONS+=("--suppress_asm_messages")
else
  NODE_OPTIONS+=("--no-suppress_asm_messages")
fi

if [ -n "$TRACE_ASM_TIME" ] && [ $TRACE_ASM_TIME == "y" ]; then
  NODE_OPTIONS+=("--trace_asm_time")
else
  NODE_OPTIONS+=("--no-trace_asm_time")
fi

if [ -n "$DUMP_WASM_MODULE" ] && [ $DUMP_WASM_MODULE == "y" ]; then
  NODE_OPTIONS+=("--dump_wasm_module")
else
  NODE_OPTIONS+=("--no-dump_wasm_module")
fi

if [ -n "$WASM_OPT" ] && [ $WASM_OPT == "y" ]; then
  NODE_OPTIONS+=("--wasm_opt")
else
  NODE_OPTIONS+=("--no-wasm_opt")
fi

if [ -n "$WASM_NO_BOUNDS_CHECKS" ] && [ $WASM_NO_BOUNDS_CHECKS == "y" ]; then
  NODE_OPTIONS+=("--wasm_no_bounds_checks")
else
  NODE_OPTIONS+=("--no-wasm_no_bounds_checks")
fi

if [ -n "$WASM_NO_STACK_CHECKS" ] && [ $WASM_NO_STACK_CHECKS == "y" ]; then
  NODE_OPTIONS+=("--wasm_no_stack_checks")
else
  NODE_OPTIONS+=("--no-wasm_no_stack_checks")
fi

if [ -n "$WASM_TRAP_HANDLER" ] && [ $WASM_TRAP_HANDLER == "y" ]; then
  NODE_OPTIONS+=("--wasm_trap_handler")
else
  NODE_OPTIONS+=("--no-wasm_trap_handler")
fi

if [ -n "$WASM_GUARD_PAGES" ] && [ $WASM_GUARD_PAGES == "y" ]; then
  NODE_OPTIONS+=("--wasm_guard_pages")
else
  NODE_OPTIONS+=("--no-wasm_guard_pages")
fi

if [ -n "$WASM_CODE_FUZZER_GEN_TEST" ] && [ $WASM_CODE_FUZZER_GEN_TEST == "y" ]; then
  NODE_OPTIONS+=("--wasm_code_fuzzer_gen_test")
else
  NODE_OPTIONS+=("--no-wasm_code_fuzzer_gen_test")
fi

if [ -n "$PRINT_WASM_CODE" ] && [ $PRINT_WASM_CODE == "y" ]; then
  NODE_OPTIONS+=("--print_wasm_code")
else
  NODE_OPTIONS+=("--no-print_wasm_code")
fi

if [ -n "$TRACE_OPT_VERBOSE" ] && [ $TRACE_OPT_VERBOSE == "y" ]; then
  NODE_OPTIONS+=("--trace_opt_verbose")
else
  NODE_OPTIONS+=("--no-trace_opt_verbose")
fi

if [ -n "$EXPERIMENTAL_NEW_SPACE_GROWTH_HEURISTIC" ] && [ $EXPERIMENTAL_NEW_SPACE_GROWTH_HEURISTIC == "y" ]; then
  NODE_OPTIONS+=("--experimental_new_space_growth_heuristic")
else
  NODE_OPTIONS+=("--no-experimental_new_space_growth_heuristic")
fi

if [ -n "$GC_GLOBAL" ] && [ $GC_GLOBAL == "y" ]; then
  NODE_OPTIONS+=("--gc_global")
else
  NODE_OPTIONS+=("--no-gc_global")
fi

if [ -n "$TRACE_GC" ] && [ $TRACE_GC == "y" ]; then
  NODE_OPTIONS+=("--trace_gc")
else
  NODE_OPTIONS+=("--no-trace_gc")
fi

if [ -n "$TRACE_GC_NVP" ] && [ $TRACE_GC_NVP == "y" ]; then
  NODE_OPTIONS+=("--trace_gc_nvp")
else
  NODE_OPTIONS+=("--no-trace_gc_nvp")
fi

if [ -n "$TRACE_GC_IGNORE_SCAVENGER" ] && [ $TRACE_GC_IGNORE_SCAVENGER == "y" ]; then
  NODE_OPTIONS+=("--trace_gc_ignore_scavenger")
else
  NODE_OPTIONS+=("--no-trace_gc_ignore_scavenger")
fi

if [ -n "$TRACE_IDLE_NOTIFICATION" ] && [ $TRACE_IDLE_NOTIFICATION == "y" ]; then
  NODE_OPTIONS+=("--trace_idle_notification")
else
  NODE_OPTIONS+=("--no-trace_idle_notification")
fi

if [ -n "$TRACE_IDLE_NOTIFICATION_VERBOSE" ] && [ $TRACE_IDLE_NOTIFICATION_VERBOSE == "y" ]; then
  NODE_OPTIONS+=("--trace_idle_notification_verbose")
else
  NODE_OPTIONS+=("--no-trace_idle_notification_verbose")
fi

if [ -n "$TRACE_GC_VERBOSE" ] && [ $TRACE_GC_VERBOSE == "y" ]; then
  NODE_OPTIONS+=("--trace_gc_verbose")
else
  NODE_OPTIONS+=("--no-trace_gc_verbose")
fi

if [ -n "$TRACE_FRAGMENTATION" ] && [ $TRACE_FRAGMENTATION == "y" ]; then
  NODE_OPTIONS+=("--trace_fragmentation")
else
  NODE_OPTIONS+=("--no-trace_fragmentation")
fi

if [ -n "$TRACE_FRAGMENTATION_VERBOSE" ] && [ $TRACE_FRAGMENTATION_VERBOSE == "y" ]; then
  NODE_OPTIONS+=("--trace_fragmentation_verbose")
else
  NODE_OPTIONS+=("--no-trace_fragmentation_verbose")
fi

if [ -n "$TRACE_EVACUATION" ] && [ $TRACE_EVACUATION == "y" ]; then
  NODE_OPTIONS+=("--trace_evacuation")
else
  NODE_OPTIONS+=("--no-trace_evacuation")
fi

if [ -n "$TRACE_MUTATOR_UTILIZATION" ] && [ $TRACE_MUTATOR_UTILIZATION == "y" ]; then
  NODE_OPTIONS+=("--trace_mutator_utilization")
else
  NODE_OPTIONS+=("--no-trace_mutator_utilization")
fi

if [ -n "$INCREMENTAL_MARKING" ] && [ $INCREMENTAL_MARKING == "y" ]; then
  NODE_OPTIONS+=("--incremental_marking")
else
  NODE_OPTIONS+=("--no-incremental_marking")
fi

if [ -n "$INCREMENTAL_MARKING_WRAPPERS" ] && [ $INCREMENTAL_MARKING_WRAPPERS == "y" ]; then
  NODE_OPTIONS+=("--incremental_marking_wrappers")
else
  NODE_OPTIONS+=("--no-incremental_marking_wrappers")
fi

if [ -n "$MINOR_MC" ] && [ $MINOR_MC == "y" ]; then
  NODE_OPTIONS+=("--minor_mc")
else
  NODE_OPTIONS+=("--no-minor_mc")
fi

if [ -n "$BLACK_ALLOCATION" ] && [ $BLACK_ALLOCATION == "y" ]; then
  NODE_OPTIONS+=("--black_allocation")
else
  NODE_OPTIONS+=("--no-black_allocation")
fi

if [ -n "$CONCURRENT_SWEEPING" ] && [ $CONCURRENT_SWEEPING == "y" ]; then
  NODE_OPTIONS+=("--concurrent_sweeping")
else
  NODE_OPTIONS+=("--no-concurrent_sweeping")
fi

if [ -n "$PARALLEL_COMPACTION" ] && [ $PARALLEL_COMPACTION == "y" ]; then
  NODE_OPTIONS+=("--parallel_compaction")
else
  NODE_OPTIONS+=("--no-parallel_compaction")
fi

if [ -n "$PARALLEL_POINTER_UPDATE" ] && [ $PARALLEL_POINTER_UPDATE == "y" ]; then
  NODE_OPTIONS+=("--parallel_pointer_update")
else
  NODE_OPTIONS+=("--no-parallel_pointer_update")
fi

if [ -n "$TRACE_INCREMENTAL_MARKING" ] && [ $TRACE_INCREMENTAL_MARKING == "y" ]; then
  NODE_OPTIONS+=("--trace_incremental_marking")
else
  NODE_OPTIONS+=("--no-trace_incremental_marking")
fi

if [ -n "$TRACK_GC_OBJECT_STATS" ] && [ $TRACK_GC_OBJECT_STATS == "y" ]; then
  NODE_OPTIONS+=("--track_gc_object_stats")
else
  NODE_OPTIONS+=("--no-track_gc_object_stats")
fi

if [ -n "$TRACE_GC_OBJECT_STATS" ] && [ $TRACE_GC_OBJECT_STATS == "y" ]; then
  NODE_OPTIONS+=("--trace_gc_object_stats")
else
  NODE_OPTIONS+=("--no-trace_gc_object_stats")
fi

if [ -n "$TRACK_DETACHED_CONTEXTS" ] && [ $TRACK_DETACHED_CONTEXTS == "y" ]; then
  NODE_OPTIONS+=("--track_detached_contexts")
else
  NODE_OPTIONS+=("--no-track_detached_contexts")
fi

if [ -n "$TRACE_DETACHED_CONTEXTS" ] && [ $TRACE_DETACHED_CONTEXTS == "y" ]; then
  NODE_OPTIONS+=("--trace_detached_contexts")
else
  NODE_OPTIONS+=("--no-trace_detached_contexts")
fi

if [ -n "$MOVE_OBJECT_START" ] && [ $MOVE_OBJECT_START == "y" ]; then
  NODE_OPTIONS+=("--move_object_start")
else
  NODE_OPTIONS+=("--no-move_object_start")
fi

if [ -n "$MEMORY_REDUCER" ] && [ $MEMORY_REDUCER == "y" ]; then
  NODE_OPTIONS+=("--memory_reducer")
else
  NODE_OPTIONS+=("--no-memory_reducer")
fi

if [ -n "$ALWAYS_COMPACT" ] && [ $ALWAYS_COMPACT == "y" ]; then
  NODE_OPTIONS+=("--always_compact")
else
  NODE_OPTIONS+=("--no-always_compact")
fi

if [ -n "$NEVER_COMPACT" ] && [ $NEVER_COMPACT == "y" ]; then
  NODE_OPTIONS+=("--never_compact")
else
  NODE_OPTIONS+=("--no-never_compact")
fi

if [ -n "$COMPACT_CODE_SPACE" ] && [ $COMPACT_CODE_SPACE == "y" ]; then
  NODE_OPTIONS+=("--compact_code_space")
else
  NODE_OPTIONS+=("--no-compact_code_space")
fi

if [ -n "$CLEANUP_CODE_CACHES_AT_GC" ] && [ $CLEANUP_CODE_CACHES_AT_GC == "y" ]; then
  NODE_OPTIONS+=("--cleanup_code_caches_at_gc")
else
  NODE_OPTIONS+=("--no-cleanup_code_caches_at_gc")
fi

if [ -n "$USE_MARKING_PROGRESS_BAR" ] && [ $USE_MARKING_PROGRESS_BAR == "y" ]; then
  NODE_OPTIONS+=("--use_marking_progress_bar")
else
  NODE_OPTIONS+=("--no-use_marking_progress_bar")
fi

if [ -n "$FORCE_MARKING_DEQUE_OVERFLOWS" ] && [ $FORCE_MARKING_DEQUE_OVERFLOWS == "y" ]; then
  NODE_OPTIONS+=("--force_marking_deque_overflows")
else
  NODE_OPTIONS+=("--no-force_marking_deque_overflows")
fi

if [ -n "$STRESS_COMPACTION" ] && [ $STRESS_COMPACTION == "y" ]; then
  NODE_OPTIONS+=("--stress_compaction")
else
  NODE_OPTIONS+=("--no-stress_compaction")
fi

if [ -n "$MANUAL_EVACUATION_CANDIDATES_SELECTION" ] && [ $MANUAL_EVACUATION_CANDIDATES_SELECTION == "y" ]; then
  NODE_OPTIONS+=("--manual_evacuation_candidates_selection")
else
  NODE_OPTIONS+=("--no-manual_evacuation_candidates_selection")
fi

if [ -n "$FAST_PROMOTION_NEW_SPACE" ] && [ $FAST_PROMOTION_NEW_SPACE == "y" ]; then
  NODE_OPTIONS+=("--fast_promotion_new_space")
else
  NODE_OPTIONS+=("--no-fast_promotion_new_space")
fi

if [ -n "$DEBUG_CODE" ] && [ $DEBUG_CODE == "y" ]; then
  NODE_OPTIONS+=("--debug_code")
else
  NODE_OPTIONS+=("--no-debug_code")
fi

if [ -n "$CODE_COMMENTS" ] && [ $CODE_COMMENTS == "y" ]; then
  NODE_OPTIONS+=("--code_comments")
else
  NODE_OPTIONS+=("--no-code_comments")
fi

if [ -n "$ENABLE_SSE3" ] && [ $ENABLE_SSE3 == "y" ]; then
  NODE_OPTIONS+=("--enable_sse3")
else
  NODE_OPTIONS+=("--no-enable_sse3")
fi

if [ -n "$ENABLE_SSSE" ] && [ $ENABLE_SSSE == "y" ]; then
  NODE_OPTIONS+=("--enable_ssse3")
else
  NODE_OPTIONS+=("--no-enable_ssse3")
fi

if [ -n "$ENABLE_SSE4_1" ] && [ $ENABLE_SSE4_1 == "y" ]; then
  NODE_OPTIONS+=("--enable_sse4_1")
else
  NODE_OPTIONS+=("--no-enable_sse4_1")
fi

if [ -n "$ENABLE_SAHF" ] && [ $ENABLE_SAHF == "y" ]; then
  NODE_OPTIONS+=("--enable_sahf")
else
  NODE_OPTIONS+=("--no-enable_sahf")
fi

if [ -n "$ENABLE_AVX" ] && [ $ENABLE_AVX == "y" ]; then
  NODE_OPTIONS+=("--enable_avx")
else
  NODE_OPTIONS+=("--no-enable_avx")
fi

if [ -n "$ENABLE_FMA" ] && [ $ENABLE_FMA == "y" ]; then
  NODE_OPTIONS+=("--enable_fma3")
else
  NODE_OPTIONS+=("--no-enable_fma3")
fi

if [ -n "$ENABLE_BMI1" ] && [ $ENABLE_BMI1 == "y" ]; then
  NODE_OPTIONS+=("--enable_bmi1")
else
  NODE_OPTIONS+=("--no-enable_bmi1")
fi

if [ -n "$ENABLE_BMI2" ] && [ $ENABLE_BMI2 == "y" ]; then
  NODE_OPTIONS+=("--enable_bmi2")
else
  NODE_OPTIONS+=("--no-enable_bmi2")
fi

if [ -n "$ENABLE_LZCNT" ] && [ $ENABLE_LZCNT == "y" ]; then
  NODE_OPTIONS+=("--enable_lzcnt")
else
  NODE_OPTIONS+=("--no-enable_lzcnt")
fi

if [ -n "$ENABLE_POPCNT" ] && [ $ENABLE_POPCNT == "y" ]; then
  NODE_OPTIONS+=("--enable_popcnt")
else
  NODE_OPTIONS+=("--no-enable_popcnt")
fi

if [ -n "$ENABLE_VLDR_IMM" ] && [ $ENABLE_VLDR_IMM == "y" ]; then
  NODE_OPTIONS+=("--enable_vldr_imm")
else
  NODE_OPTIONS+=("--no-enable_vldr_imm")
fi

if [ -n "$FORCE_LONG_BRANCHES" ] && [ $FORCE_LONG_BRANCHES == "y" ]; then
  NODE_OPTIONS+=("--force_long_branches")
else
  NODE_OPTIONS+=("--no-force_long_branches")
fi

if [ -n "$ENABLE_ARMV7" ] && [ $ENABLE_ARMV7 == "y" ]; then
  NODE_OPTIONS+=("--enable_armv7")
else
  NODE_OPTIONS+=("--no-enable_armv7")
fi

if [ -n "$ENABLE_VFP" ] && [ $ENABLE_VFP == "y" ]; then
  NODE_OPTIONS+=("--enable_vfp3")
else
  NODE_OPTIONS+=("--no-enable_vfp3")
fi

if [ -n "$ENABLE_" ] && [ $ENABLE_ == "y" ]; then
  NODE_OPTIONS+=("--enable_32dregs")
else
  NODE_OPTIONS+=("--no-enable_32dregs")
fi

if [ -n "$ENABLE_NEON" ] && [ $ENABLE_NEON == "y" ]; then
  NODE_OPTIONS+=("--enable_neon")
else
  NODE_OPTIONS+=("--no-enable_neon")
fi

if [ -n "$ENABLE_SUDIV" ] && [ $ENABLE_SUDIV == "y" ]; then
  NODE_OPTIONS+=("--enable_sudiv")
else
  NODE_OPTIONS+=("--no-enable_sudiv")
fi

if [ -n "$ENABLE_ARMV8" ] && [ $ENABLE_ARMV8 == "y" ]; then
  NODE_OPTIONS+=("--enable_armv8")
else
  NODE_OPTIONS+=("--no-enable_armv8")
fi

if [ -n "$ENABLE_REGEXP_UNALIGNED_ACCESSES" ] && [ $ENABLE_REGEXP_UNALIGNED_ACCESSES == "y" ]; then
  NODE_OPTIONS+=("--enable_regexp_unaligned_accesses")
else
  NODE_OPTIONS+=("--no-enable_regexp_unaligned_accesses")
fi

if [ -n "$SCRIPT_STREAMING" ] && [ $SCRIPT_STREAMING == "y" ]; then
  NODE_OPTIONS+=("--script_streaming")
else
  NODE_OPTIONS+=("--no-script_streaming")
fi

if [ -n "$DISABLE_OLD_API_ACCESSORS" ] && [ $DISABLE_OLD_API_ACCESSORS == "y" ]; then
  NODE_OPTIONS+=("--disable_old_api_accessors")
else
  NODE_OPTIONS+=("--no-disable_old_api_accessors")
fi

if [ -n "$EXPOSE_FREE_BUFFER" ] && [ $EXPOSE_FREE_BUFFER == "y" ]; then
  NODE_OPTIONS+=("--expose_free_buffer")
else
  NODE_OPTIONS+=("--no-expose_free_buffer")
fi

if [ -n "$EXPOSE_GC" ] && [ $EXPOSE_GC == "y" ]; then
  NODE_OPTIONS+=("--expose_gc")
else
  NODE_OPTIONS+=("--no-expose_gc")
fi

if [ -n "$EXPOSE_EXTERNALIZE_STRING" ] && [ $EXPOSE_EXTERNALIZE_STRING == "y" ]; then
  NODE_OPTIONS+=("--expose_externalize_string")
else
  NODE_OPTIONS+=("--no-expose_externalize_string")
fi

if [ -n "$EXPOSE_TRIGGER_FAILURE" ] && [ $EXPOSE_TRIGGER_FAILURE == "y" ]; then
  NODE_OPTIONS+=("--expose_trigger_failure")
else
  NODE_OPTIONS+=("--no-expose_trigger_failure")
fi

if [ -n "$BUILTINS_IN_STACK_TRACES" ] && [ $BUILTINS_IN_STACK_TRACES == "y" ]; then
  NODE_OPTIONS+=("--builtins_in_stack_traces")
else
  NODE_OPTIONS+=("--no-builtins_in_stack_traces")
fi

if [ -n "$ALLOW_UNSAFE_FUNCTION_CONSTRUCTOR" ] && [ $ALLOW_UNSAFE_FUNCTION_CONSTRUCTOR == "y" ]; then
  NODE_OPTIONS+=("--allow_unsafe_function_constructor")
else
  NODE_OPTIONS+=("--no-allow_unsafe_function_constructor")
fi

if [ -n "$INLINE_NEW" ] && [ $INLINE_NEW == "y" ]; then
  NODE_OPTIONS+=("--inline_new")
else
  NODE_OPTIONS+=("--no-inline_new")
fi

if [ -n "$TRACE_CODEGEN" ] && [ $TRACE_CODEGEN == "y" ]; then
  NODE_OPTIONS+=("--trace_codegen")
else
  NODE_OPTIONS+=("--no-trace_codegen")
fi

if [ -n "$LAZY" ] && [ $LAZY == "y" ]; then
  NODE_OPTIONS+=("--lazy")
else
  NODE_OPTIONS+=("--no-lazy")
fi

if [ -n "$TRACE_OPT" ] && [ $TRACE_OPT == "y" ]; then
  NODE_OPTIONS+=("--trace_opt")
else
  NODE_OPTIONS+=("--no-trace_opt")
fi

if [ -n "$TRACE_OPT_STATS" ] && [ $TRACE_OPT_STATS == "y" ]; then
  NODE_OPTIONS+=("--trace_opt_stats")
else
  NODE_OPTIONS+=("--no-trace_opt_stats")
fi

if [ -n "$TRACE_FILE_NAMES" ] && [ $TRACE_FILE_NAMES == "y" ]; then
  NODE_OPTIONS+=("--trace_file_names")
else
  NODE_OPTIONS+=("--no-trace_file_names")
fi

if [ -n "$OPT" ] && [ $OPT == "y" ]; then
  NODE_OPTIONS+=("--opt")
else
  NODE_OPTIONS+=("--no-opt")
fi

if [ -n "$ALWAYS_OPT" ] && [ $ALWAYS_OPT == "y" ]; then
  NODE_OPTIONS+=("--always_opt")
else
  NODE_OPTIONS+=("--no-always_opt")
fi

if [ -n "$ALWAYS_OSR" ] && [ $ALWAYS_OSR == "y" ]; then
  NODE_OPTIONS+=("--always_osr")
else
  NODE_OPTIONS+=("--no-always_osr")
fi

if [ -n "$PREPARE_ALWAYS_OPT" ] && [ $PREPARE_ALWAYS_OPT == "y" ]; then
  NODE_OPTIONS+=("--prepare_always_opt")
else
  NODE_OPTIONS+=("--no-prepare_always_opt")
fi

if [ -n "$TRACE_DEOPT" ] && [ $TRACE_DEOPT == "y" ]; then
  NODE_OPTIONS+=("--trace_deopt")
else
  NODE_OPTIONS+=("--no-trace_deopt")
fi

if [ -n "$SERIALIZE_TOPLEVEL" ] && [ $SERIALIZE_TOPLEVEL == "y" ]; then
  NODE_OPTIONS+=("--serialize_toplevel")
else
  NODE_OPTIONS+=("--no-serialize_toplevel")
fi

if [ -n "$SERIALIZE_EAGER" ] && [ $SERIALIZE_EAGER == "y" ]; then
  NODE_OPTIONS+=("--serialize_eager")
else
  NODE_OPTIONS+=("--no-serialize_eager")
fi

if [ -n "$TRACE_SERIALIZER" ] && [ $TRACE_SERIALIZER == "y" ]; then
  NODE_OPTIONS+=("--trace_serializer")
else
  NODE_OPTIONS+=("--no-trace_serializer")
fi

if [ -n "$COMPILATION_CACHE" ] && [ $COMPILATION_CACHE == "y" ]; then
  NODE_OPTIONS+=("--compilation_cache")
else
  NODE_OPTIONS+=("--no-compilation_cache")
fi

if [ -n "$CACHE_PROTOTYPE_TRANSITIONS" ] && [ $CACHE_PROTOTYPE_TRANSITIONS == "y" ]; then
  NODE_OPTIONS+=("--cache_prototype_transitions")
else
  NODE_OPTIONS+=("--no-cache_prototype_transitions")
fi

if [ -n "$COMPILER_DISPATCHER" ] && [ $COMPILER_DISPATCHER == "y" ]; then
  NODE_OPTIONS+=("--compiler_dispatcher")
else
  NODE_OPTIONS+=("--no-compiler_dispatcher")
fi

if [ -n "$TRACE_COMPILER_DISPATCHER" ] && [ $TRACE_COMPILER_DISPATCHER == "y" ]; then
  NODE_OPTIONS+=("--trace_compiler_dispatcher")
else
  NODE_OPTIONS+=("--no-trace_compiler_dispatcher")
fi

if [ -n "$TRACE_COMPILER_DISPATCHER_JOBS" ] && [ $TRACE_COMPILER_DISPATCHER_JOBS == "y" ]; then
  NODE_OPTIONS+=("--trace_compiler_dispatcher_jobs")
else
  NODE_OPTIONS+=("--no-trace_compiler_dispatcher_jobs")
fi

if [ -n "$TRACE_JS_ARRAY_ABUSE" ] && [ $TRACE_JS_ARRAY_ABUSE == "y" ]; then
  NODE_OPTIONS+=("--trace_js_array_abuse")
else
  NODE_OPTIONS+=("--no-trace_js_array_abuse")
fi

if [ -n "$TRACE_EXTERNAL_ARRAY_ABUSE" ] && [ $TRACE_EXTERNAL_ARRAY_ABUSE == "y" ]; then
  NODE_OPTIONS+=("--trace_external_array_abuse")
else
  NODE_OPTIONS+=("--no-trace_external_array_abuse")
fi

if [ -n "$TRACE_ARRAY_ABUSE" ] && [ $TRACE_ARRAY_ABUSE == "y" ]; then
  NODE_OPTIONS+=("--trace_array_abuse")
else
  NODE_OPTIONS+=("--no-trace_array_abuse")
fi

if [ -n "$ENABLE_LIVEEDIT" ] && [ $ENABLE_LIVEEDIT == "y" ]; then
  NODE_OPTIONS+=("--enable_liveedit")
else
  NODE_OPTIONS+=("--no-enable_liveedit")
fi

if [ -n "$TRACE_SIDE_EFFECT_FREE_DEBUG_EVALUATE" ] && [ $TRACE_SIDE_EFFECT_FREE_DEBUG_EVALUATE == "y" ]; then
  NODE_OPTIONS+=("--trace_side_effect_free_debug_evaluate")
else
  NODE_OPTIONS+=("--no-trace_side_effect_free_debug_evaluate")
fi

if [ -n "$HARD_ABORT" ] && [ $HARD_ABORT == "y" ]; then
  NODE_OPTIONS+=("--hard_abort")
else
  NODE_OPTIONS+=("--no-hard_abort")
fi

if [ -n "$CLEAR_EXCEPTIONS_ON_JS_ENTRY" ] && [ $CLEAR_EXCEPTIONS_ON_JS_ENTRY == "y" ]; then
  NODE_OPTIONS+=("--clear_exceptions_on_js_entry")
else
  NODE_OPTIONS+=("--no-clear_exceptions_on_js_entry")
fi

if [ -n "$HEAP_PROFILER_TRACE_OBJECTS" ] && [ $HEAP_PROFILER_TRACE_OBJECTS == "y" ]; then
  NODE_OPTIONS+=("--heap_profiler_trace_objects")
else
  NODE_OPTIONS+=("--no-heap_profiler_trace_objects")
fi

if [ -n "$SAMPLING_HEAP_PROFILER_SUPPRESS_RANDOMNESS" ] && [ $SAMPLING_HEAP_PROFILER_SUPPRESS_RANDOMNESS == "y" ]; then
  NODE_OPTIONS+=("--sampling_heap_profiler_suppress_randomness")
else
  NODE_OPTIONS+=("--no-sampling_heap_profiler_suppress_randomness")
fi

if [ -n "$USE_IDLE_NOTIFICATION" ] && [ $USE_IDLE_NOTIFICATION == "y" ]; then
  NODE_OPTIONS+=("--use_idle_notification")
else
  NODE_OPTIONS+=("--no-use_idle_notification")
fi

if [ -n "$USE_IC" ] && [ $USE_IC == "y" ]; then
  NODE_OPTIONS+=("--use_ic")
else
  NODE_OPTIONS+=("--no-use_ic")
fi

if [ -n "$TRACE_IC" ] && [ $TRACE_IC == "y" ]; then
  NODE_OPTIONS+=("--trace_ic")
else
  NODE_OPTIONS+=("--no-trace_ic")
fi

if [ -n "$NATIVE_CODE_COUNTERS" ] && [ $NATIVE_CODE_COUNTERS == "y" ]; then
  NODE_OPTIONS+=("--native_code_counters")
else
  NODE_OPTIONS+=("--no-native_code_counters")
fi

if [ -n "$THIN_STRINGS" ] && [ $THIN_STRINGS == "y" ]; then
  NODE_OPTIONS+=("--thin_strings")
else
  NODE_OPTIONS+=("--no-thin_strings")
fi

if [ -n "$TRACE_WEAK_ARRAYS" ] && [ $TRACE_WEAK_ARRAYS == "y" ]; then
  NODE_OPTIONS+=("--trace_weak_arrays")
else
  NODE_OPTIONS+=("--no-trace_weak_arrays")
fi

if [ -n "$TRACE_PROTOTYPE_USERS" ] && [ $TRACE_PROTOTYPE_USERS == "y" ]; then
  NODE_OPTIONS+=("--trace_prototype_users")
else
  NODE_OPTIONS+=("--no-trace_prototype_users")
fi

if [ -n "$USE_VERBOSE_PRINTER" ] && [ $USE_VERBOSE_PRINTER == "y" ]; then
  NODE_OPTIONS+=("--use_verbose_printer")
else
  NODE_OPTIONS+=("--no-use_verbose_printer")
fi

if [ -n "$TRACE_FOR_IN_ENUMERATE" ] && [ $TRACE_FOR_IN_ENUMERATE == "y" ]; then
  NODE_OPTIONS+=("--trace_for_in_enumerate")
else
  NODE_OPTIONS+=("--no-trace_for_in_enumerate")
fi

if [ -n "$ALLOW_NATIVES_SYNTAX" ] && [ $ALLOW_NATIVES_SYNTAX == "y" ]; then
  NODE_OPTIONS+=("--allow_natives_syntax")
else
  NODE_OPTIONS+=("--no-allow_natives_syntax")
fi

if [ -n "$TRACE_PARSE" ] && [ $TRACE_PARSE == "y" ]; then
  NODE_OPTIONS+=("--trace_parse")
else
  NODE_OPTIONS+=("--no-trace_parse")
fi

if [ -n "$TRACE_PREPARSE" ] && [ $TRACE_PREPARSE == "y" ]; then
  NODE_OPTIONS+=("--trace_preparse")
else
  NODE_OPTIONS+=("--no-trace_preparse")
fi

if [ -n "$LAZY_INNER_FUNCTIONS" ] && [ $LAZY_INNER_FUNCTIONS == "y" ]; then
  NODE_OPTIONS+=("--lazy_inner_functions")
else
  NODE_OPTIONS+=("--no-lazy_inner_functions")
fi

if [ -n "$AGGRESSIVE_LAZY_INNER_FUNCTIONS" ] && [ $AGGRESSIVE_LAZY_INNER_FUNCTIONS == "y" ]; then
  NODE_OPTIONS+=("--aggressive_lazy_inner_functions")
else
  NODE_OPTIONS+=("--no-aggressive_lazy_inner_functions")
fi

if [ -n "$PREPARSER_SCOPE_ANALYSIS" ] && [ $PREPARSER_SCOPE_ANALYSIS == "y" ]; then
  NODE_OPTIONS+=("--preparser_scope_analysis")
else
  NODE_OPTIONS+=("--no-preparser_scope_analysis")
fi

if [ -n "$TRACE_SIM" ] && [ $TRACE_SIM == "y" ]; then
  NODE_OPTIONS+=("--trace_sim")
else
  NODE_OPTIONS+=("--no-trace_sim")
fi

if [ -n "$DEBUG_SIM" ] && [ $DEBUG_SIM == "y" ]; then
  NODE_OPTIONS+=("--debug_sim")
else
  NODE_OPTIONS+=("--no-debug_sim")
fi

if [ -n "$CHECK_ICACHE" ] && [ $CHECK_ICACHE == "y" ]; then
  NODE_OPTIONS+=("--check_icache")
else
  NODE_OPTIONS+=("--no-check_icache")
fi

if [ -n "$LOG_REGS_MODIFIED" ] && [ $LOG_REGS_MODIFIED == "y" ]; then
  NODE_OPTIONS+=("--log_regs_modified")
else
  NODE_OPTIONS+=("--no-log_regs_modified")
fi

if [ -n "$LOG_COLOUR" ] && [ $LOG_COLOUR == "y" ]; then
  NODE_OPTIONS+=("--log_colour")
else
  NODE_OPTIONS+=("--no-log_colour")
fi

if [ -n "$IGNORE_ASM_UNIMPLEMENTED_BREAK" ] && [ $IGNORE_ASM_UNIMPLEMENTED_BREAK == "y" ]; then
  NODE_OPTIONS+=("--ignore_asm_unimplemented_break")
else
  NODE_OPTIONS+=("--no-ignore_asm_unimplemented_break")
fi

if [ -n "$TRACE_SIM_MESSAGES" ] && [ $TRACE_SIM_MESSAGES == "y" ]; then
  NODE_OPTIONS+=("--trace_sim_messages")
else
  NODE_OPTIONS+=("--no-trace_sim_messages")
fi

if [ -n "$STACK_TRACE_ON_ILLEGAL" ] && [ $STACK_TRACE_ON_ILLEGAL == "y" ]; then
  NODE_OPTIONS+=("--stack_trace_on_illegal")
else
  NODE_OPTIONS+=("--no-stack_trace_on_illegal")
fi

if [ -n "$ABORT_ON_UNCAUGHT_EXCEPTION" ] && [ $ABORT_ON_UNCAUGHT_EXCEPTION == "y" ]; then
  NODE_OPTIONS+=("--abort_on_uncaught_exception")
else
  NODE_OPTIONS+=("--no-abort_on_uncaught_exception")
fi

if [ -n "$RANDOMIZE_HASHES" ] && [ $RANDOMIZE_HASHES == "y" ]; then
  NODE_OPTIONS+=("--randomize_hashes")
else
  NODE_OPTIONS+=("--no-randomize_hashes")
fi

if [ -n "$TRACE_RAIL" ] && [ $TRACE_RAIL == "y" ]; then
  NODE_OPTIONS+=("--trace_rail")
else
  NODE_OPTIONS+=("--no-trace_rail")
fi

if [ -n "$PRINT_ALL_EXCEPTIONS" ] && [ $PRINT_ALL_EXCEPTIONS == "y" ]; then
  NODE_OPTIONS+=("--print_all_exceptions")
else
  NODE_OPTIONS+=("--no-print_all_exceptions")
fi

if [ -n "$RUNTIME_CALL_STATS" ] && [ $RUNTIME_CALL_STATS == "y" ]; then
  NODE_OPTIONS+=("--runtime_call_stats")
else
  NODE_OPTIONS+=("--no-runtime_call_stats")
fi

if [ -n "$PROFILE_DESERIALIZATION" ] && [ $PROFILE_DESERIALIZATION == "y" ]; then
  NODE_OPTIONS+=("--profile_deserialization")
else
  NODE_OPTIONS+=("--no-profile_deserialization")
fi

if [ -n "$SERIALIZATION_STATISTICS" ] && [ $SERIALIZATION_STATISTICS == "y" ]; then
  NODE_OPTIONS+=("--serialization_statistics")
else
  NODE_OPTIONS+=("--no-serialization_statistics")
fi

if [ -n "$REGEXP_OPTIMIZATION" ] && [ $REGEXP_OPTIMIZATION == "y" ]; then
  NODE_OPTIONS+=("--regexp_optimization")
else
  NODE_OPTIONS+=("--no-regexp_optimization")
fi

if [ -n "$TESTING_BOOL_FLAG" ] && [ $TESTING_BOOL_FLAG == "y" ]; then
  NODE_OPTIONS+=("--testing_bool_flag")
else
  NODE_OPTIONS+=("--no-testing_bool_flag")
fi

if [ -n "$TESTING_MAYBE_BOOL_FLAG" ] && [ $TESTING_MAYBE_BOOL_FLAG == "y" ]; then
  NODE_OPTIONS+=("--testing_maybe_bool_flag")
else
  NODE_OPTIONS+=("--no-testing_maybe_bool_flag")
fi

if [ -n "$DUMP_COUNTERS" ] && [ $DUMP_COUNTERS == "y" ]; then
  NODE_OPTIONS+=("--dump_counters")
else
  NODE_OPTIONS+=("--no-dump_counters")
fi

if [ -n "$DUMP_COUNTERS_NVP" ] && [ $DUMP_COUNTERS_NVP == "y" ]; then
  NODE_OPTIONS+=("--dump_counters_nvp")
else
  NODE_OPTIONS+=("--no-dump_counters_nvp")
fi

if [ -n "$LOG" ] && [ $LOG == "y" ]; then
  NODE_OPTIONS+=("--log")
else
  NODE_OPTIONS+=("--no-log")
fi

if [ -n "$LOG_ALL" ] && [ $LOG_ALL == "y" ]; then
  NODE_OPTIONS+=("--log_all")
else
  NODE_OPTIONS+=("--no-log_all")
fi

if [ -n "$LOG_API" ] && [ $LOG_API == "y" ]; then
  NODE_OPTIONS+=("--log_api")
else
  NODE_OPTIONS+=("--no-log_api")
fi

if [ -n "$LOG_CODE" ] && [ $LOG_CODE == "y" ]; then
  NODE_OPTIONS+=("--log_code")
else
  NODE_OPTIONS+=("--no-log_code")
fi

if [ -n "$LOG_GC" ] && [ $LOG_GC == "y" ]; then
  NODE_OPTIONS+=("--log_gc")
else
  NODE_OPTIONS+=("--no-log_gc")
fi

if [ -n "$LOG_HANDLES" ] && [ $LOG_HANDLES == "y" ]; then
  NODE_OPTIONS+=("--log_handles")
else
  NODE_OPTIONS+=("--no-log_handles")
fi

if [ -n "$LOG_SUSPECT" ] && [ $LOG_SUSPECT == "y" ]; then
  NODE_OPTIONS+=("--log_suspect")
else
  NODE_OPTIONS+=("--no-log_suspect")
fi

if [ -n "$PROF" ] && [ $PROF == "y" ]; then
  NODE_OPTIONS+=("--prof")
else
  NODE_OPTIONS+=("--no-prof")
fi

if [ -n "$PROF_CPP" ] && [ $PROF_CPP == "y" ]; then
  NODE_OPTIONS+=("--prof_cpp")
else
  NODE_OPTIONS+=("--no-prof_cpp")
fi

if [ -n "$PROF_BROWSER_MODE" ] && [ $PROF_BROWSER_MODE == "y" ]; then
  NODE_OPTIONS+=("--prof_browser_mode")
else
  NODE_OPTIONS+=("--no-prof_browser_mode")
fi

if [ -n "$LOGFILE_PER_ISOLATE" ] && [ $LOGFILE_PER_ISOLATE == "y" ]; then
  NODE_OPTIONS+=("--logfile_per_isolate")
else
  NODE_OPTIONS+=("--no-logfile_per_isolate")
fi

if [ -n "$LL_PROF" ] && [ $LL_PROF == "y" ]; then
  NODE_OPTIONS+=("--ll_prof")
else
  NODE_OPTIONS+=("--no-ll_prof")
fi

if [ -n "$PERF_BASIC_PROF" ] && [ $PERF_BASIC_PROF == "y" ]; then
  NODE_OPTIONS+=("--perf_basic_prof")
else
  NODE_OPTIONS+=("--no-perf_basic_prof")
fi

if [ -n "$PERF_BASIC_PROF_ONLY_FUNCTIONS" ] && [ $PERF_BASIC_PROF_ONLY_FUNCTIONS == "y" ]; then
  NODE_OPTIONS+=("--perf_basic_prof_only_functions")
else
  NODE_OPTIONS+=("--no-perf_basic_prof_only_functions")
fi

if [ -n "$PERF_PROF" ] && [ $PERF_PROF == "y" ]; then
  NODE_OPTIONS+=("--perf_prof")
else
  NODE_OPTIONS+=("--no-perf_prof")
fi

if [ -n "$PERF_PROF_UNWINDING_INFO" ] && [ $PERF_PROF_UNWINDING_INFO == "y" ]; then
  NODE_OPTIONS+=("--perf_prof_unwinding_info")
else
  NODE_OPTIONS+=("--no-perf_prof_unwinding_info")
fi

if [ -n "$LOG_INTERNAL_TIMER_EVENTS" ] && [ $LOG_INTERNAL_TIMER_EVENTS == "y" ]; then
  NODE_OPTIONS+=("--log_internal_timer_events")
else
  NODE_OPTIONS+=("--no-log_internal_timer_events")
fi

if [ -n "$LOG_TIMER_EVENTS" ] && [ $LOG_TIMER_EVENTS == "y" ]; then
  NODE_OPTIONS+=("--log_timer_events")
else
  NODE_OPTIONS+=("--no-log_timer_events")
fi

if [ -n "$LOG_INSTRUCTION_STATS" ] && [ $LOG_INSTRUCTION_STATS == "y" ]; then
  NODE_OPTIONS+=("--log_instruction_stats")
else
  NODE_OPTIONS+=("--no-log_instruction_stats")
fi

if [ -n "$REDIRECT_CODE_TRACES" ] && [ $REDIRECT_CODE_TRACES == "y" ]; then
  NODE_OPTIONS+=("--redirect_code_traces")
else
  NODE_OPTIONS+=("--no-redirect_code_traces")
fi

if [ -n "$PRINT_OPT_SOURCE" ] && [ $PRINT_OPT_SOURCE == "y" ]; then
  NODE_OPTIONS+=("--print_opt_source")
else
  NODE_OPTIONS+=("--no-print_opt_source")
fi

if [ -n "$TRACE_ELEMENTS_TRANSITIONS" ] && [ $TRACE_ELEMENTS_TRANSITIONS == "y" ]; then
  NODE_OPTIONS+=("--trace_elements_transitions")
else
  NODE_OPTIONS+=("--no-trace_elements_transitions")
fi

if [ -n "$TRACE_CREATION_ALLOCATION_SITES" ] && [ $TRACE_CREATION_ALLOCATION_SITES == "y" ]; then
  NODE_OPTIONS+=("--trace_creation_allocation_sites")
else
  NODE_OPTIONS+=("--no-trace_creation_allocation_sites")
fi

if [ -n "$PRINT_CODE_STUBS" ] && [ $PRINT_CODE_STUBS == "y" ]; then
  NODE_OPTIONS+=("--print_code_stubs")
else
  NODE_OPTIONS+=("--no-print_code_stubs")
fi

if [ -n "$TEST_SECONDARY_STUB_CACHE" ] && [ $TEST_SECONDARY_STUB_CACHE == "y" ]; then
  NODE_OPTIONS+=("--test_secondary_stub_cache")
else
  NODE_OPTIONS+=("--no-test_secondary_stub_cache")
fi

if [ -n "$TEST_PRIMARY_STUB_CACHE" ] && [ $TEST_PRIMARY_STUB_CACHE == "y" ]; then
  NODE_OPTIONS+=("--test_primary_stub_cache")
else
  NODE_OPTIONS+=("--no-test_primary_stub_cache")
fi

if [ -n "$TEST_SMALL_MAX_FUNCTION_CONTEXT_STUB_SIZE" ] && [ $TEST_SMALL_MAX_FUNCTION_CONTEXT_STUB_SIZE == "y" ]; then
  NODE_OPTIONS+=("--test_small_max_function_context_stub_size")
else
  NODE_OPTIONS+=("--no-test_small_max_function_context_stub_size")
fi

if [ -n "$PRINT_CODE" ] && [ $PRINT_CODE == "y" ]; then
  NODE_OPTIONS+=("--print_code")
else
  NODE_OPTIONS+=("--no-print_code")
fi

if [ -n "$PRINT_OPT_CODE" ] && [ $PRINT_OPT_CODE == "y" ]; then
  NODE_OPTIONS+=("--print_opt_code")
else
  NODE_OPTIONS+=("--no-print_opt_code")
fi

if [ -n "$PRINT_CODE_VERBOSE" ] && [ $PRINT_CODE_VERBOSE == "y" ]; then
  NODE_OPTIONS+=("--print_code_verbose")
else
  NODE_OPTIONS+=("--no-print_code_verbose")
fi

if [ -n "$PRINT_BUILTIN_CODE" ] && [ $PRINT_BUILTIN_CODE == "y" ]; then
  NODE_OPTIONS+=("--print_builtin_code")
else
  NODE_OPTIONS+=("--no-print_builtin_code")
fi

if [ -n "$SODIUM" ] && [ $SODIUM == "y" ]; then
  NODE_OPTIONS+=("--sodium")
else
  NODE_OPTIONS+=("--no-sodium")
fi

if [ -n "$PRINT_ALL_CODE" ] && [ $PRINT_ALL_CODE == "y" ]; then
  NODE_OPTIONS+=("--print_all_code")
else
  NODE_OPTIONS+=("--no-print_all_code")
fi

if [ -n "$PREDICTABLE" ] && [ $PREDICTABLE == "y" ]; then
  NODE_OPTIONS+=("--predictable")
else
  NODE_OPTIONS+=("--no-predictable")
fi

if [ -n "$SINGLE_THREADED" ] && [ $SINGLE_THREADED == "y" ]; then
  NODE_OPTIONS+=("--single_threaded")
else
  NODE_OPTIONS+=("--no-single_threaded")
fi

sleep 1

# Open the port for the client to connect to
coproc ./nc -l $WAYFINDER_DOMAIN_IP_ADDR 3000

cd ./bench

node ${NODE_OPTIONS[@]} dist/cli.js > results

cat results

tail -n 1 results > results_parsed

# Dirt format the header and send it to the client
echo -n -e "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nConnection: " "Keep-Alive" "\r\nDate: Mon, 01 Jan 1970 00:00:00 GMT GMT\r\nContent-Length: $(wc -c results_parsed)\r\n\r\n" > header
cat header results_parsed <&"${COPROC[0]}" >&"${COPROC[1]}"

# Stop the server
kill -9 $COPROC_PID

sleep 60
