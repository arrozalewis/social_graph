Houseparty Backend Coding Challenge README

JSON structure can be found here.

https://drive.google.com/file/d/1ROHHxzhOs1wW5bKIVnHdIz1XFDQk3xGO/view

It contains one JSON object per line, which represents a

friend or unfriend event in a fake social service. For example:

{"to": {"id": "Cf9GG2IU9f0=", "name": "Nora Merril"}, "from":
{"id": "8H8a3l+O8+c=", "name": "Kenn Arne"}, "timestamp":
1520963452875, "areFriends": true}

{"to": {"id": "kMJL5+dqops=", "name": "Sho Karen"}, "from": {"id":
"UsCYj+02tRQ=", "name": "Kiki Aaron"}, "timestamp": 1520963452876,
"areFriends": false}

From the first line, Nora and Kenn became friends. Then Kiki "unfriended" Sho, and so on. The timestamp is
the number of milliseconds since Epoch (Jan 1st, 1970).


CommandLine Tool
Consists of driver.go and folder parse (package)
Input events are read in driver.go and data is written to objects and data
structures in the parse package

to run the program use the following
go run driver.go [json] [arg2] [arg3(optional)] [arg4(optional)]

If the wrong arguements are passed in or the file is unreadable to
program will exit 

CLT is capable of handling two user ids were it is assumed that
1) ever user has their own unique "id"
2) all users events are read in chronological order

A map of interfaces is used to store data pertaining to JSON objects for
outputting
-JSON list of friends
-JSON list of mutual friends

If anyone is unfriended at the latest timestamp then they will not be
considered as a connection or be output as friends

Bonuses
Implemented Six Degrees of Kevin Bacon

Algorithm used for solving degrees of seperation is a depth first search and
backtracking solution which always finds the minimum number of seperations b/w
two people if one exists

If there is absolutely no path from one user to another, then the degree of
seperation will be outputed as -1
