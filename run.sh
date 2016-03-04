#!/bin/bash

for t in noTimeout timeAfter ticker timer context1 ; do
    chantimeout $t
done
