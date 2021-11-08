package proc
// SPDX-License-Identifier: BSD-3-Clause
//
// Authors: Alexander Jung <a.jung@lancs.ac.uk>
//
// Copyright (c) 2021, Lancaster University.  All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

import (
  "os"
  "fmt"
  "strconv"
  "io/ioutil"

  "github.com/unikraft/wayfinder/config"
)

// ProcPIDNetstat represents entries of /proc/<pid>/net/netstat file
type ProcPIDNetstat struct {
  PID                             int

  TCPExtSyncookiesSent            uint64
  TCPExtSyncookiesRecv            uint64
  TCPExtSyncookiesFailed          uint64
  TCPExtEmbryonicRsts             uint64
  TCPExtPruneCalled               uint64
  TCPExtRcvPruned                 uint64
  TCPExtOfoPruned                 uint64
  TCPExtOutOfWindowIcmps          uint64
  TCPExtLockDroppedIcmps          uint64
  TCPExtArpFilter                 uint64
  TCPExtTW                        uint64
  TCPExtTWRecycled                uint64
  TCPExtTWKilled                  uint64
  TCPExtPAWSActive                uint64
  TCPExtPAWSEstab                 uint64
  TCPExtDelayedACKs               uint64
  TCPExtDelayedACKLocked          uint64
  TCPExtDelayedACKLost            uint64
  TCPExtListenOverflows           uint64
  TCPExtListenDrops               uint64
  TCPExtTCPHPHits                 uint64
  TCPExtTCPPureAcks               uint64
  TCPExtTCPHPAcks                 uint64
  TCPExtTCPRenoRecovery           uint64
  TCPExtTCPSackRecovery           uint64
  TCPExtTCPSACKReneging           uint64
  TCPExtTCPSACKReorder            uint64
  TCPExtTCPRenoReorder            uint64
  TCPExtTCPTSReorder              uint64
  TCPExtTCPFullUndo               uint64
  TCPExtTCPPartialUndo            uint64
  TCPExtTCPDSACKUndo              uint64
  TCPExtTCPLossUndo               uint64
  TCPExtTCPLostRetransmit         uint64
  TCPExtTCPRenoFailures           uint64
  TCPExtTCPSackFailures           uint64
  TCPExtTCPLossFailures           uint64
  TCPExtTCPFastRetrans            uint64
  TCPExtTCPSlowStartRetrans       uint64
  TCPExtTCPTimeouts               uint64
  TCPExtTCPLossProbes             uint64
  TCPExtTCPLossProbeRecovery      uint64
  TCPExtTCPRenoRecoveryFail       uint64
  TCPExtTCPSackRecoveryFail       uint64
  TCPExtTCPRcvCollapsed           uint64
  TCPExtTCPDSACKOldSent           uint64
  TCPExtTCPDSACKOfoSent           uint64
  TCPExtTCPDSACKRecv              uint64
  TCPExtTCPDSACKOfoRecv           uint64
  TCPExtTCPAbortOnData            uint64
  TCPExtTCPAbortOnClose           uint64
  TCPExtTCPAbortOnMemory          uint64
  TCPExtTCPAbortOnTimeout         uint64
  TCPExtTCPAbortOnLinger          uint64
  TCPExtTCPAbortFailed            uint64
  TCPExtTCPMemoryPressures        uint64
  TCPExtTCPMemoryPressuresChrono  uint64
  TCPExtTCPSACKDiscard            uint64
  TCPExtTCPDSACKIgnoredOld        uint64
  TCPExtTCPDSACKIgnoredNoUndo     uint64
  TCPExtTCPSpuriousRTOs           uint64
  TCPExtTCPMD5NotFound            uint64
  TCPExtTCPMD5Unexpected          uint64
  TCPExtTCPMD5Failure             uint64
  TCPExtTCPSackShifted            uint64
  TCPExtTCPSackMerged             uint64
  TCPExtTCPSackShiftFallback      uint64
  TCPExtTCPBacklogDrop            uint64
  TCPExtPFMemallocDrop            uint64
  TCPExtTCPMinTTLDrop             uint64
  TCPExtTCPDeferAcceptDrop        uint64
  TCPExtIPReversePathFilter       uint64
  TCPExtTCPTimeWaitOverflow       uint64
  TCPExtTCPReqQFullDoCookies      uint64
  TCPExtTCPReqQFullDrop           uint64
  TCPExtTCPRetransFail            uint64
  TCPExtTCPRcvCoalesce            uint64
  TCPExtTCPOFOQueue               uint64
  TCPExtTCPOFODrop                uint64
  TCPExtTCPOFOMerge               uint64
  TCPExtTCPChallengeACK           uint64
  TCPExtTCPSYNChallenge           uint64
  TCPExtTCPFastOpenActive         uint64
  TCPExtTCPFastOpenActiveFail     uint64
  TCPExtTCPFastOpenPassive        uint64
  TCPExtTCPFastOpenPassiveFail    uint64
  TCPExtTCPFastOpenListenOverflow uint64
  TCPExtTCPFastOpenCookieReqd     uint64
  TCPExtTCPFastOpenBlackhole      uint64
  TCPExtTCPSpuriousRtxHostQueues  uint64
  TCPExtBusyPollRxPackets         uint64
  TCPExtTCPAutoCorking            uint64
  TCPExtTCPFromZeroWindowAdv      uint64
  TCPExtTCPToZeroWindowAdv        uint64
  TCPExtTCPWantZeroWindowAdv      uint64
  TCPExtTCPSynRetrans             uint64
  TCPExtTCPOrigDataSent           uint64
  TCPExtTCPHystartTrainDetect     uint64
  TCPExtTCPHystartTrainCwnd       uint64
  TCPExtTCPHystartDelayDetect     uint64
  TCPExtTCPHystartDelayCwnd       uint64
  TCPExtTCPACKSkippedSynRecv      uint64
  TCPExtTCPACKSkippedPAWS         uint64
  TCPExtTCPACKSkippedSeq          uint64
  TCPExtTCPACKSkippedFinWait2     uint64
  TCPExtTCPACKSkippedTimeWait     uint64
  TCPExtTCPACKSkippedChallenge    uint64
  TCPExtTCPWinProbe               uint64
  TCPExtTCPKeepAlive              uint64
  TCPExtTCPMTUPFail               uint64
  TCPExtTCPMTUPSuccess            uint64
  TCPExtTCPDelivered              uint64
  TCPExtTCPDeliveredCE            uint64
  TCPExtTCPAckCompressed          uint64

  IPExtInNoRoutes                 uint64
  IPExtInTruncatedPkts            uint64
  IPExtInMcastPkts                uint64
  IPExtOutMcastPkts               uint64
  IPExtInBcastPkts                uint64
  IPExtOutBcastPkts               uint64
  IPExtInOctets                   uint64
  IPExtOutOctets                  uint64
  IPExtInMcastOctets              uint64
  IPExtOutMcastOctets             uint64
  IPExtInBcastOctets              uint64
  IPExtOutBcastOctets             uint64
  IPExtInCsumErrors               uint64
  IPExtInNoECTPkts                uint64
  IPExtInECT1Pkts                 uint64
  IPExtInECT0Pkts                 uint64
  InCEPkts                        uint64
}

// GetProcPIDNetstat reads the netstat file for given pid from procfs
func GetProcPIDNetstat(pid int) ProcPIDNetstat {
  stats := ProcPIDNetstat{PID: pid}

  filepath := fmt.Sprint(config.RuntimeConfig.ProcFS, "/", strconv.Itoa(pid), "/net/netstat")
  filecontent, err := ioutil.ReadFile(filepath)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Cannot read proc netstat: %s\n", err)
    return ProcPIDNetstat{}
  }

  ioFormat := "" +
    "TcpExt: %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s\n" +
    "TcpExt: %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d\n" +
    "IpExt: %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s\n" +
    "IpExt: %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d\n"

  var ignore string

  _, err = fmt.Sscanf(
    string(filecontent), ioFormat,

    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,

    &stats.TCPExtSyncookiesSent,
    &stats.TCPExtSyncookiesRecv,
    &stats.TCPExtSyncookiesFailed,
    &stats.TCPExtEmbryonicRsts,
    &stats.TCPExtPruneCalled,
    &stats.TCPExtRcvPruned,
    &stats.TCPExtOfoPruned,
    &stats.TCPExtOutOfWindowIcmps,
    &stats.TCPExtLockDroppedIcmps,
    &stats.TCPExtArpFilter,
    &stats.TCPExtTW,
    &stats.TCPExtTWRecycled,
    &stats.TCPExtTWKilled,
    &stats.TCPExtPAWSActive,
    &stats.TCPExtPAWSEstab,
    &stats.TCPExtDelayedACKs,
    &stats.TCPExtDelayedACKLocked,
    &stats.TCPExtDelayedACKLost,
    &stats.TCPExtListenOverflows,
    &stats.TCPExtListenDrops,
    &stats.TCPExtTCPHPHits,
    &stats.TCPExtTCPPureAcks,
    &stats.TCPExtTCPHPAcks,
    &stats.TCPExtTCPRenoRecovery,
    &stats.TCPExtTCPSackRecovery,
    &stats.TCPExtTCPSACKReneging,
    &stats.TCPExtTCPSACKReorder,
    &stats.TCPExtTCPRenoReorder,
    &stats.TCPExtTCPTSReorder,
    &stats.TCPExtTCPFullUndo,
    &stats.TCPExtTCPPartialUndo,
    &stats.TCPExtTCPDSACKUndo,
    &stats.TCPExtTCPLossUndo,
    &stats.TCPExtTCPLostRetransmit,
    &stats.TCPExtTCPRenoFailures,
    &stats.TCPExtTCPSackFailures,
    &stats.TCPExtTCPLossFailures,
    &stats.TCPExtTCPFastRetrans,
    &stats.TCPExtTCPSlowStartRetrans,
    &stats.TCPExtTCPTimeouts,
    &stats.TCPExtTCPLossProbes,
    &stats.TCPExtTCPLossProbeRecovery,
    &stats.TCPExtTCPRenoRecoveryFail,
    &stats.TCPExtTCPSackRecoveryFail,
    &stats.TCPExtTCPRcvCollapsed,
    &stats.TCPExtTCPDSACKOldSent,
    &stats.TCPExtTCPDSACKOfoSent,
    &stats.TCPExtTCPDSACKRecv,
    &stats.TCPExtTCPDSACKOfoRecv,
    &stats.TCPExtTCPAbortOnData,
    &stats.TCPExtTCPAbortOnClose,
    &stats.TCPExtTCPAbortOnMemory,
    &stats.TCPExtTCPAbortOnTimeout,
    &stats.TCPExtTCPAbortOnLinger,
    &stats.TCPExtTCPAbortFailed,
    &stats.TCPExtTCPMemoryPressures,
    &stats.TCPExtTCPMemoryPressuresChrono,
    &stats.TCPExtTCPSACKDiscard,
    &stats.TCPExtTCPDSACKIgnoredOld,
    &stats.TCPExtTCPDSACKIgnoredNoUndo,
    &stats.TCPExtTCPSpuriousRTOs,
    &stats.TCPExtTCPMD5NotFound,
    &stats.TCPExtTCPMD5Unexpected,
    &stats.TCPExtTCPMD5Failure,
    &stats.TCPExtTCPSackShifted,
    &stats.TCPExtTCPSackMerged,
    &stats.TCPExtTCPSackShiftFallback,
    &stats.TCPExtTCPBacklogDrop,
    &stats.TCPExtPFMemallocDrop,
    &stats.TCPExtTCPMinTTLDrop,
    &stats.TCPExtTCPDeferAcceptDrop,
    &stats.TCPExtIPReversePathFilter,
    &stats.TCPExtTCPTimeWaitOverflow,
    &stats.TCPExtTCPReqQFullDoCookies,
    &stats.TCPExtTCPReqQFullDrop,
    &stats.TCPExtTCPRetransFail,
    &stats.TCPExtTCPRcvCoalesce,
    &stats.TCPExtTCPOFOQueue,
    &stats.TCPExtTCPOFODrop,
    &stats.TCPExtTCPOFOMerge,
    &stats.TCPExtTCPChallengeACK,
    &stats.TCPExtTCPSYNChallenge,
    &stats.TCPExtTCPFastOpenActive,
    &stats.TCPExtTCPFastOpenActiveFail,
    &stats.TCPExtTCPFastOpenPassive,
    &stats.TCPExtTCPFastOpenPassiveFail,
    &stats.TCPExtTCPFastOpenListenOverflow,
    &stats.TCPExtTCPFastOpenCookieReqd,
    &stats.TCPExtTCPFastOpenBlackhole,
    &stats.TCPExtTCPSpuriousRtxHostQueues,
    &stats.TCPExtBusyPollRxPackets,
    &stats.TCPExtTCPAutoCorking,
    &stats.TCPExtTCPFromZeroWindowAdv,
    &stats.TCPExtTCPToZeroWindowAdv,
    &stats.TCPExtTCPWantZeroWindowAdv,
    &stats.TCPExtTCPSynRetrans,
    &stats.TCPExtTCPOrigDataSent,
    &stats.TCPExtTCPHystartTrainDetect,
    &stats.TCPExtTCPHystartTrainCwnd,
    &stats.TCPExtTCPHystartDelayDetect,
    &stats.TCPExtTCPHystartDelayCwnd,
    &stats.TCPExtTCPACKSkippedSynRecv,
    &stats.TCPExtTCPACKSkippedPAWS,
    &stats.TCPExtTCPACKSkippedSeq,
    &stats.TCPExtTCPACKSkippedFinWait2,
    &stats.TCPExtTCPACKSkippedTimeWait,
    &stats.TCPExtTCPACKSkippedChallenge,
    &stats.TCPExtTCPWinProbe,
    &stats.TCPExtTCPKeepAlive,
    &stats.TCPExtTCPMTUPFail,
    &stats.TCPExtTCPMTUPSuccess,
    &stats.TCPExtTCPDelivered,
    &stats.TCPExtTCPDeliveredCE,
    &stats.TCPExtTCPAckCompressed,

    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,
    &ignore,

    &stats.IPExtInNoRoutes,
    &stats.IPExtInTruncatedPkts,
    &stats.IPExtInMcastPkts,
    &stats.IPExtOutMcastPkts,
    &stats.IPExtInBcastPkts,
    &stats.IPExtOutBcastPkts,
    &stats.IPExtInOctets,
    &stats.IPExtOutOctets,
    &stats.IPExtInMcastOctets,
    &stats.IPExtOutMcastOctets,
    &stats.IPExtInBcastOctets,
    &stats.IPExtOutBcastOctets,
    &stats.IPExtInCsumErrors,
    &stats.IPExtInNoECTPkts,
    &stats.IPExtInECT1Pkts,
    &stats.IPExtInECT0Pkts,
    &stats.InCEPkts,
  )

  if err != nil {
    fmt.Fprintf(os.Stderr, "Cannot parse proc netstat: %s\n", err)
    return ProcPIDNetstat{}
  }

  return stats
}
