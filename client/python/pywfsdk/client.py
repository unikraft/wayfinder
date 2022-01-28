# SPDX-License-Identifier: BSD-3-Clause
#
# Authors: Alexander Jung <alex@unikraft.io>
#
# Copyright (c) 2022, Unikraft UG.  All rights reserved.
#
# Redistribution and use in source and binary forms, with or without
# modification, are permitted provided that the following conditions
# are met:
#
# 1. Redistributions of source code must retain the above copyright
#    notice, this list of conditions and the following disclaimer.
# 2. Redistributions in binary form must reproduce the above copyright
#    notice, this list of conditions and the following disclaimer in the
#    documentation and/or other materials provided with the distribution.
# 3. Neither the name of the copyright holder nor the names of its
#    contributors may be used to endorse or promote products derived from
#    this software without specific prior written permission.
#
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
# AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
# IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
# ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
# LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
# CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
# SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
# INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
# CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
# ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
# POSSIBILITY OF SUCH DAMAGE.

"""
# pywfsdk.client

This module implements the main Wayfinder client utility.

## Basic Usage

Create a `WayfinderClient` singleton:

```python
from pywfsdk import WayfinderClient 

wf = WayfinderClient(
  host='localhost',
  port=5000,
  disable_ssl_verification=True
)
```
"""

from __future__ import absolute_import
from __future__ import unicode_literals

import ssl
import betterproto
from grpclib.client import Channel
from typing import Optional, Type, List

try:
  from .proto import wayfinder
except ImportError as e:
  raise RuntimeError("Please run setuptools for pywfsdk for code generation")

from pywfsdk.util import log

class WayfinderClient:
  """
  The Wayfinder SDK Client Object.

  Applications should configure the client at startup time and continue to use
  it throughout the lifetime of the application, rather than creating
  instances on the fly.
  """
  def __init__(self,
               host: str,
               port: int,
               server_cert: Optional[str]=None,
               client_cert: Optional[str]=None,
               client_key: Optional[str]=None,
               disable_ssl_verification: bool=False,
               channel: Optional[Channel]=None):
    self.__host = host
    self.__port = port
    self.__server_cert = server_cert
    self.__client_cert = client_cert
    self.__client_key = client_key
    self.__disable_ssl_verification = disable_ssl_verification
    self.__channel = channel

    # Initialise the dependent stubs
    self.__builder = self.create_stub(wayfinder.BuilderServiceStub)
    self.__hosts = self.create_stub(wayfinder.HostConfigServiceStub)
    self.__jobs = self.create_stub(wayfinder.JobServiceStub)
    self.__tester = self.create_stub(wayfinder.TesterServiceStub)

    log.info("Started Wayfinder client")

  @property
  def host(self) -> str:
    """
    The remote host, such as an IP or FQDN, of the Wayfinder server.
    """
    return self.__host
  
  @property
  def port(self) -> int:
    """
    The port the remote Wayfinder server is listeninig on.  Default is `5000`.
    """
    return self.__port

  @property
  def server_cert(self) -> str:
    """
    The path to the server's public certificate should the connection requires
    SSL.
    """
    return self.__server_cert

  @property
  def client_key(self) -> str:
    """
    The path to the client's private key should the connection requires SSL.
    """
    return self.__client_key

  @property
  def client_cert(self) -> str:
    """
    The path to the client's public certificate should the connection requires
    SSL.
    """
    return self.__client_cert

  @property
  def disable_ssl_verification(self) -> bool:
    """
    Whether the client should verify the SSL connection.
    """
    return self.__disable_ssl_verification
  
  @property
  def jobs(self) -> wayfinder.JobServiceStub:
    """
    Returns an instantiate `wayfinder.JobServiceStub` stub.
    """
    return self.__jobs

  @property
  def channel(self) -> Channel:
    """
    The `grpc.Channel` object which gRPC connections are handled over.
    """
    if self.__channel is not None:
      return self.__channel
    
    ssl_context = not self.disable_ssl_verification

    if self.client_cert and self.client_key and self.server_cert:
      ssl_context = ssl.SSLContext(ssl.PROTOCOL_TLS)
      if self.disable_ssl_verification:
        ssl_context.verify_mode = ssl.CERT_REQUIRED

      ssl_context.load_cert_chain(self.client_cert, self.client_key)
      ssl_context.load_verify_locations(self.server_cert)
      ssl_context.options |= ssl.OP_NO_TLSv1 | ssl.OP_NO_TLSv1_1
      ssl_context.set_ciphers('ECDHE+AESGCM:ECDHE+CHACHA20:DHE+AESGCM:DHE+CHACHA20')
      ssl_context.set_alpn_protocols(['h2'])
      try:
          ssl_context.set_npn_protocols(['h2'])
      except NotImplementedError:
          pass
      return ssl_context

    self.__channel = Channel(
      host=self.host,
      port=self.port,
      ssl=ssl_context
    )
    
    return self.__channel

  def create_stub(self, stub: Type[betterproto.ServiceStub]) -> betterproto.ServiceStub:
    """
    Instantiate a stub based on the client's configuration channel.
    """
    return stub(channel=self.channel)

  # Magic methods allow a client object to be automatically cleaned up by the
  # "with" scope operator.
  def __enter__(self):
    return self

  def __exit__(self, type, value, traceback):
    self.close()

  def close(self):
    """
    Releases the gRPC connection channel.  Do not attempt to use the client
    after calling this method.
    """

    log.info("Closing Wayfinder client...")
    if self.channel:
      self.channel.close()
