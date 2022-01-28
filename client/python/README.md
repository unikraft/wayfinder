# Python3 SDK Bindings for Wayfinder

## Installation

```bash
git clone https://github.com/unikraft/wayfinder.git
cd client/python
make install
```

## Usage

Create a `WayfinderClient` singleton:

```python
from pywfsdk import WayfinderClient 

wf = WayfinderClient(
  host='localhost',
  port=5000,
  disable_ssl_verification=True
)
```
