# gadget - utilities for Go projects

[![Build Status](https://travis-ci.org/marcusva/gadget.svg?branch=master)](https://travis-ci.org/marcusva/gadget)

Go snippets commonly used in my projects.

* ``config``: provides a simple access to INI-style configuration files.
* ``log``: a simple wrapper around the log package, which adds RFC5424 severity
  thresholds.
* ``set``: a simple set implementation.
* ``testing``: provides minimalistic testing enhancements and a CSV fuzzer.

# Usage
Simply copy and paste the go files into the desired location of your Go project.
The different packages are usually self-contained and only depend on the Go
standard library.

The accompanying test files depend on the provided ``testing`` package.

# License

The code is placed into the public domain. There are no licensing restrictions.
For details, please check the ``LICENSE`` file within the repository.