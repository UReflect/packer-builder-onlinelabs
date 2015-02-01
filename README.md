# packer-builder-onlinelabs

[![Build Status](https://travis-ci.org/meatballhat/packer-builder-onlinelabs.svg?branch=master)](https://travis-ci.org/meatballhat/packer-builder-onlinelabs)

Packer builder for [onlinelabs](http://labs.online.net/).

**FAIR WARNING**: This plugin "works", but has not been battle tested, nor has it been endorsed by Online Labs or Packer
:smiley_cat:.  Many of the problems encountered during testing were due to incomplete support for arm on Ubuntu, fwiw.

Example output:
```
==> Building simple-box
onlinelabs output will be in this color.

==> onlinelabs: Creating server...
==> onlinelabs: Waiting for server to become active...
==> onlinelabs: Waiting for SSH to become available...
==> onlinelabs: Connected to SSH!
==> onlinelabs: Provisioning with shell script: scripts/echo
    onlinelabs: + env
    onlinelabs: + date -u
==> onlinelabs: Gracefully shutting down server...
==> onlinelabs: Forcefully shutting down server...
==> onlinelabs: Creating snapshot: packer-1422820258
==> onlinelabs: Destroying server...
Build 'onlinelabs' finished.

==> Builds finished. The artifacts of successful builds are:
--> onlinelabs: A snapshot was created: '9c9a7175-f774-4c9b-97db-5f8bfd41ea16' (packer-1422820258)
```
