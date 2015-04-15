#!/bin/bash

for i in `seq 1 34`; do
    bd=`printf "%4d-%02d-%02d" $(( ( RANDOM % 10 )  + 2000 )) $(( ( RANDOM % 12 )  + 1 )) $(( ( RANDOM % 28 )  + 1 ))`
    printf -v ln "LastName%03d" $i
    printf -v gn "John%03d" $i
    echo "{ \"familyname\": \"LongLastName$i\", \"givenname\" : \"John$i\", \"birthdate\" : \"$bd\" }"
    curl -w "%{http_code}\n" -XPOST http://davidk.bainbridge%40gmail.com@localhost:8080/children -d \
        "{ \"familyname\": \"$ln\", \"givenname\" : \"$gn\", \"birthdate\" : \"$bd\" }"
done
