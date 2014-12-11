# tbd

A decentralised ci server using docker. Any linux machine can be a CI agent with Docker.

Features:
* Pushes responsibility for green builds back onto the developers
* Builds are certified and signed locally with git as the central source of truth for build success
* Avoids the requirement for complex, ui-based hub and spoke build server models
* No more commit and forget workflow _No more chasing people for broken builds.....EVER!!!_
* Larger builds can be run remotely, but certified and signed locally
* Build logs are published in the container for that build

## Interface

