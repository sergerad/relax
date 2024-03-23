#!/bin/bash

for example_dir in ./examples/*; do
	echo "---"
	echo "Running $example_dir"
	go run "$example_dir"
done
