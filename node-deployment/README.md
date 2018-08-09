# Prerequisites

Autoupdater runs on master and needs to inspect the image to be updated using the docker api. This lookup is performed using the repository name and
the tag. However, if the image is downloaded by swarm it won't be tagged and therefore unable to find with inspect.

In order to avoid this, skywire images should first be pulled from the master node:

`docker pull skycoin/skywire`
