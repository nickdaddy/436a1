#!/bin/bash
# Builds 436A1

#prj = "436a1repo"

#for CMD in `ls -d 436a1repo`; do
#  go install ./436a1repo/$CMD
#done

#cd $GOPATH/src
go install ./serverA1/
go install ./clientA1
