package stats

import (
  // libvirt "github.com/libvirt/libvirt-go"
  
  "github.com/unikraft/wayfinder/modules/libvirt"
)

type DomainStats struct {
  domain *libvirt.Domain
}
