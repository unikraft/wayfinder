#!/bin/sh

mkdir -p result

set +ex

for benchmark in ft mg cg is
do
    for class in S W A B
    do
        echo "running $benchmark.$class"
        $benchmark.$class.x 1> result/$benchmark.$class.out
        echo "done\n"
    done
done
