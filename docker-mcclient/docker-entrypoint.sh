#!/bin/bash

# For directories backed by PVs, we copy over any contents we saved on in the Dockerfile
for d in config mcdb
do
    find /opt/vconsole -name $d -type d -empty -exec cp -r /tmp/pv-seed/$d/ /opt/vconsole \;
done

/opt/vconsole/bin/mctl start || exit 1
tail -f /opt/vconsole/log/mc/mconsole.log
