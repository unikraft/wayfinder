#!/bin/bash

set -ex

if [ -n "${CONFIG_NETWORK_SECMARK}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_NETWORK_SECMARK ${CONFIG_NETWORK_SECMARK}
fi

if [ -n "${CONFIG_CC_OPTIMIZE_FOR_SIZE}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_CC_OPTIMIZE_FOR_SIZE ${CONFIG_CC_OPTIMIZE_FOR_SIZE}
fi

if [ -n "${CONFIG_SECURITY_SELINUX}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_SECURITY_SELINUX ${CONFIG_SECURITY_SELINUX}
fi

if [ -n "${CONFIG_CPU_FREQ_STAT}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_CPU_FREQ_STAT ${CONFIG_CPU_FREQ_STAT}
fi

if [ -n "${CONFIG_CROSS_MEMORY_ATTACH}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_CROSS_MEMORY_ATTACH ${CONFIG_CROSS_MEMORY_ATTACH}
fi

if [ -n "${CONFIG_HIGH_RES_TIMERS}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_HIGH_RES_TIMERS ${CONFIG_HIGH_RES_TIMERS}
fi

if [ -n "${CONFIG_TASKS_RCU}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_TASKS_RCU ${CONFIG_TASKS_RCU}
fi

if [ -n "${CONFIG_RCU_USER_QS}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_RCU_USER_QS ${CONFIG_RCU_USER_QS}
fi

if [ -n "${CONFIG_RCU_FANOUT}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_RCU_FANOUT ${CONFIG_RCU_FANOUT}
fi

if [ -n "${CONFIG_RCU_FANOUT_LEAF}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_RCU_FANOUT_LEAF ${CONFIG_RCU_FANOUT_LEAF}
fi

if [ -n "${CONFIG_RCU_FANOUT_EXACT}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_RCU_FANOUT_EXACT ${CONFIG_RCU_FANOUT_EXACT}
fi

if [ -n "${CONFIG_RCU_FAST_NO_HZ}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_RCU_FAST_NO_HZ ${CONFIG_RCU_FAST_NO_HZ}
fi

if [ -n "${CONFIG_RCU_KTHREAD_PRIO}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_RCU_KTHREAD_PRIO ${CONFIG_RCU_KTHREAD_PRIO}
fi

if [ -n "${CONFIG_RCU_NOCB_CPU}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_RCU_NOCB_CPU ${CONFIG_RCU_NOCB_CPU}
fi

if [ -n "${CONFIG_NUMA_BALANCING}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_NUMA_BALANCING ${CONFIG_NUMA_BALANCING}
fi

if [ -n "${CONFIG_NUMA_BALANCING_DEFAULT_ENABLED}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_NUMA_BALANCING_DEFAULT_ENABLED ${CONFIG_NUMA_BALANCING_DEFAULT_ENABLED}
fi

if [ -n "${CONFIG_SCHED_AUTOGROUP}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_SCHED_AUTOGROUP ${CONFIG_SCHED_AUTOGROUP}
fi

if [ -n "${CONFIG_PRINTK}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_PRINTK ${CONFIG_PRINTK}
fi

if [ -n "${CONFIG_BASE_FULL}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_BASE_FULL ${CONFIG_BASE_FULL}
fi

if [ -n "${CONFIG_FUTEX}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_FUTEX ${CONFIG_FUTEX}
fi

if [ -n "${CONFIG_EPOLL}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_EPOLL ${CONFIG_EPOLL}
fi

if [ -n "${CONFIG_AIO}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_AIO ${CONFIG_AIO}
fi

if [ -n "${CONFIG_ADVISE_SYSCALLS}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_ADVISE_SYSCALLS ${CONFIG_ADVISE_SYSCALLS}
fi

if [ -n "${CONFIG_SLUB_CPU_PARTIAL}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_SLUB_CPU_PARTIAL ${CONFIG_SLUB_CPU_PARTIAL}
fi

if [ -n "${CONFIG_JUMP_LABEL}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_JUMP_LABEL ${CONFIG_JUMP_LABEL}
fi

if [ -n "${CONFIG_DEFAULT_IOSCHED}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_DEFAULT_IOSCHED ${CONFIG_DEFAULT_IOSCHED}
fi

if [ -n "${CONFIG_SCHED_OMIT_FRAME_POINTER}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_SCHED_OMIT_FRAME_POINTER ${CONFIG_SCHED_OMIT_FRAME_POINTER}
fi

if [ -n "${CONFIG_SMP}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_SMP ${CONFIG_SMP}
fi

if [ -n "${CONFIG_SCHED_SMT}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_SCHED_SMT ${CONFIG_SCHED_SMT}
fi

if [ -n "${CONFIG_DIRECT_GBPAGES}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_DIRECT_GBPAGES ${CONFIG_DIRECT_GBPAGES}
fi

if [ -n "${CONFIG_NUMA}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_NUMA ${CONFIG_NUMA}
fi

if [ -n "${CONFIG_SPARSEMEM_VMEMMAP}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_SPARSEMEM_VMEMMAP ${CONFIG_SPARSEMEM_VMEMMAP}
fi

if [ -n "${CONFIG_TRANSPARENT_HUGEPAGE}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_TRANSPARENT_HUGEPAGE ${CONFIG_TRANSPARENT_HUGEPAGE}
fi

if [ -n "${CONFIG_X86_SMAP}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_X86_SMAP ${CONFIG_X86_SMAP}
fi

if [ -n "${CONFIG_WQ_POWER_EFFICIENT_DEFAULT}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_WQ_POWER_EFFICIENT_DEFAULT ${CONFIG_WQ_POWER_EFFICIENT_DEFAULT}
fi

if [ -n "${CONFIG_SUSPEND}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_SUSPEND ${CONFIG_SUSPEND}
fi

if [ -n "${CONFIG_HIBERNATION}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_HIBERNATION ${CONFIG_HIBERNATION}
fi

if [ -n "${CONFIG_PM_AUTOSLEEP}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_PM_AUTOSLEEP ${CONFIG_PM_AUTOSLEEP}
fi

if [ -n "${CONFIG_DEBUG_TLBFLUSH}" ]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_DEBUG_TLBFLUSH ${CONFIG_DEBUG_TLBFLUSH}
fi

if [[ ${CPU_FREQ_GOVERNOR} == "performance" ]]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_CPU_FREQ_DEFAULT_GOV_PERFORMANCE ${CONFIG_CPU_FREQ_DEFAULT_GOV_PERFORMANCE}
fi

if [[ ${CPU_FREQ_GOVERNOR} == "powersave" ]]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_CPU_FREQ_DEFAULT_GOV_POWERSAVE ${CONFIG_CPU_FREQ_DEFAULT_GOV_POWERSAVE}
fi

if [[ ${CPU_FREQ_GOVERNOR} == "userspace" ]]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_CPU_FREQ_DEFAULT_GOV_USERSPACE ${CONFIG_CPU_FREQ_DEFAULT_GOV_USERSPACE}
fi

if [[ ${CPU_FREQ_GOVERNOR} == "ondemand" ]]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_CPU_FREQ_DEFAULT_GOV_ONDEMAND ${CONFIG_CPU_FREQ_DEFAULT_GOV_ONDEMAND}
fi

if [[ ${CPU_FREQ_GOVERNOR} == "conservative" ]]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_CPU_FREQ_DEFAULT_GOV_CONSERVATIVE ${CONFIG_CPU_FREQ_DEFAULT_GOV_CONSERVATIVE}
fi

if [[ ${TIMER_FREQUENCY} == "hz_100" ]]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_HZ_100 ${CONFIG_HZ_100}
fi

if [[ ${TIMER_FREQUENCY} == "hz_250" ]]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_HZ_250 ${CONFIG_HZ_250}
fi

if [[ ${TIMER_FREQUENCY} == "hz_300" ]]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_HZ_300 ${CONFIG_HZ_300}
fi

if [[ ${TIMER_FREQUENCY} == "hz_1000" ]]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_HZ_1000 ${CONFIG_HZ_1000}
fi

if [[ ${TIMER_TICK_HANDLING} == "hz_periodic" ]]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_HZ_PERIODIC ${CONFIG_HZ_PERIODIC}
fi

if [[ ${TIMER_TICK_HANDLING} == "no_hz_idle" ]]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_NO_HZ_IDLE ${CONFIG_NO_HZ_IDLE}
fi

if [[ ${TIMER_TICK_HANDLING} == "no_hz_full" ]]; then
    kconfig-set.sh ${MICROVM_CFG} CONFIG_NO_HZ_FULL ${CONFIG_NO_HZ_FULL}
fi
