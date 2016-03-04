#!/bin/bash

for t in noTimeout timeAfter ticker timer context1 context2 ; do
    chantimeout $t
done
