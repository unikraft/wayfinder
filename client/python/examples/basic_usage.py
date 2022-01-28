#!/bin/env python3

import asyncio
from pywfsdk import WayfinderClient

async def main():
  wf = WayfinderClient(
    host='localhost',
    port=5000,
    disable_ssl_verification=True
  )

  # Get list of jobs
  jobs = await wf.jobs.list_jobs()

  for job in jobs:
    print(job)

  wf.close()

if __name__ == '__main__':
  loop = asyncio.get_event_loop()
  loop.run_until_complete(main())
