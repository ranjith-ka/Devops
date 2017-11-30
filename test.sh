#!/bin/bash
# Tested using bash version 4.1.5
for ((i = 1; i <= 100; i++)); do
    if [ $(($i % 5)) -eq 0 ] && [ $(($i % 3)) -eq 0 ]; then
        echo "$i is divisible by 3 & 5"

    elif [ $(($i % 5)) -eq 0 ]; then
        echo "$i is divisible by 5"

    else
        echo "$i is not divisible by 3 & 5 "
    fi
done
