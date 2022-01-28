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

from setuptools.command.build_py import build_py as _build_py
from setuptools import setup, find_packages

import os
import pathlib
import setuptools
import subprocess

here = pathlib.Path(__file__).parent.resolve()

# Get the long description from the README file
long_description = (here / 'README.md').read_text(encoding='utf-8')

# Get all the dependency requirements for the project
with open('requirements.txt') as f:
    required = [x for x in f.read().splitlines() if not x.startswith("#")]

SETUP_DIR = os.path.dirname(os.path.abspath(__file__))

class build_proto(setuptools.Command):
    description = 'build python bindings from proto files'
    user_options = []

    def initialize_options(self):
        pass

    def run(self):
        try:
            os.makedirs(os.path.join(SETUP_DIR, "src/pywfsdk/proto"), exist_ok = True)
            cmd = [
                "protoc",
                "-I" + os.path.join(SETUP_DIR, "../../api"),
                "--python_betterproto_out=" + os.path.join(SETUP_DIR, "src/pywfsdk/proto"),
                "-I" + os.path.join(SETUP_DIR, "vendor/googleapis"),
                os.path.join(SETUP_DIR, "../../api/*.proto")
            ]
            subprocess.run(
                " ".join(cmd),
                shell=True,
                check=True
            )

        except (subprocess.CalledProcessError, OSError) as e:
            print(e)
            print('ERROR: problems with compiling protoc. Is protoc installed?')
            raise SystemExit

    def finalize_options(self):
        pass


class build_py(_build_py):
    def run(self):
        self.run_command("build_proto")
        _build_py.run(self)


setup(
    name='pywfsdk',
    version='0.1.0',
    description='Python3 SDK Bindings for Wayfinder',
    long_description=long_description,
    long_description_content_type='text/markdown',
    url='https://github.com/unikraft/wayfinder',
    classifiers=[
        'Development Status :: 3 - Alpha',
        'Intended Audience :: Science/Research',
        'Intended Audience :: Developers',
        'Topic :: Scientific/Engineering',
        'Topic :: Software Development :: Build Tools',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.6',
        'Programming Language :: Python :: 3.7',
        'Programming Language :: Python :: 3.8',
        'Programming Language :: Python :: 3.9',
        "Programming Language :: Python :: 3.10",
        'Programming Language :: Python :: 3 :: Only',
    ],
    package_dir={'': '.'},
    packages=find_packages(where='.'),
    python_requires='>=3.6, <4',
    setup_requires=required,
    project_urls={
        'Bug Reports': 'https://github.com/unikraft/pywfsdk/issues',
        'Source': 'https://github.com/unikraft/pywfsdk/',
    },
    cmdclass={
        'build_py': build_py,
        'build_proto': build_proto,
    },
)
