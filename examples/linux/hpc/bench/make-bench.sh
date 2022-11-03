#!/bin/sh

mkdir -p bin

for benchmark in ft mg cg is
do
    for class in S W A B
    do
        echo "compiling $benchmark.$class"
        make $benchmark CLASS=$class
    done
done
