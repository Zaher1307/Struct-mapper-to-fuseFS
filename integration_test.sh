#!/bin/sh

go run ./demo.go mnt &
sleep 0.3

String=$(cat ./mnt/String)
if [[ $String != 'str' ]]; then
    echo 'TEST FAILED: file "String" does not match struct value'
    fusermount -zu ./mnt
    exit 1
fi

Int=$(cat ./mnt/Int)
if [[ $Int != '18' ]]; then
    echo 'TEST FAILED: file "Int" does not match struct value'
    fusermount -zu ./mnt
    exit 1
fi

Bool=$(cat ./mnt/Bool)
if [[ $Bool != 'true' ]]; then
    echo 'TEST FAILED: file "Bool" does not match struct value'
    fusermount -zu ./mnt
    exit 1
fi

find ./mnt/SubStructure >> /dev/null
if [[ $? != 0 ]]; then
    echo 'TEST FAILED: dir "SubStructure" does not exist'
fi

Float=$(cat ./mnt/SubStructure/Float)
if [[ $Float != 1.3 ]]; then
    echo 'TEST FAILED: file "Float" does not match struct value'
    fusermount -zu ./mnt
    exit 1
fi

sleep 1
String=$(cat ./mnt/String)
if [[ $String != 'new string' ]]; then
    echo 'TEST FAILED: file "String" does not match struct value'
    fusermount -zu ./mnt
    exit 1
fi

fusermount -zu ./mnt
echo 'TEST PASSED'
