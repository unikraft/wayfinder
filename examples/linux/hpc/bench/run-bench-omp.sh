#!/bin/bash

mkdir -p result

set +ex

for benchmark in ft mg cg is
do
    for class in S W A B
    do
        for num_thread in 4 16
        do
            export OMP_NUM_THREADS="$num_thread"
            echo "running $benchmark.$class.$num_thread"
            $benchmark.$class.x 1> result/$benchmark.$class.$num_thread.out
            echo "done\n"
        done
    done
done
