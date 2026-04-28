#! /bin/bash
cd traceinspector
go build -o ../front_src/bin/traceinspector ./cmd/traceinspector/main.go
cd ..