---
title: Melrōse installation warnings
---

[Home](https://emicklei.github.io/melrose)

## Installation warnings

If you install melrōse from a Github release then you are installing software that is not verified by any party other than the developer who published the release. With newer versions of operating systems, both Apple and Microsoft are more restrictive when it comes to installing software. Currently, it is still allowed to install unregistered software (most open-source packages are) but the you, the user, will be asked to accept the risk.

### Apple Mac OSX
On this operating system, you have to pass 3 steps of security checks:

1. When downloading the release (.zip archive) from Github, your computer will detect that it contains an application. It will ask you to proceed.

2. When starting the application melrōse using the script (`run.sh`), your computer will detect that the developer of the application is not verified. It will ask you to accept (only once).

3. When the application melrōse is loading an extra library for MIDI access, your comnputer will detect that the developer of the application is not verified. It will ask you to accept (only once).
