# Bloom clock
Implementation of Proof of Concept based on https://arxiv.org/abs/1905.13064

## Running
First get the necessary packages 

go get -d ./... 

or one by one

go get github.com/gorilla/mux

go get github.com/spencerkimball/cbfilter


go get github.com/DATA-DOG/godog


Usage:
go build main.go 


To start a node:

./main start PORT neighbors N1 N2 ...

To send element:

./main send PORT element

To check if it has

./main has PORT element

To compare bloom clock:

./main compare PORT1 PORT2


To send csv file use:


./main csv PORT1 PORT2

# For running everything automatically

chmod a+x run_auto.sh

chmod a+x kill_nodes.sh

To start nodes:

./run_auto.sh

And then to kill the process:


./kill_nodes.sh


# Testing
(It is assumed that script run_auto.sh is up and running) Run tests: 

godog
