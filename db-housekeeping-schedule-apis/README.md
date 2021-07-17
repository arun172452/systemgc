<!--
title: .'VM Housekeeping Lambdas'
description: 'Code to GET VM/EC2 list and POST VM schedule submission'Cancel changes
framework: v1
platform: AWS
language: Go
-->
 

Following endpoints provided:
1. Post endpoint  /dev/registry

Sample input Json:

{
  "userid": "arunl",
  "vmname": "ukvm037",
  "dbtype": "oracle",
  "dburl": "jdbc:oracle:thin:@192.168.34.25:1622:xe",
  "dbuser": "rskadmin",
  "dbscript": "yes"
 }
 
{
  "userid": "arunl",
  "vmname": "ukvm037",
  "dbtype": "oracle",
  "dburl": "jdbc:oracle:thin:@192.168.34.25:1622:xe",
  "dbuser": "rskadmin",
  "dbscript": "no"
 }

 {
  "userid": "arunl",
  "vmname": "ukvm037",
  "dbtype": "postgresql",
  "dburl": "jdbc:oracle:thin:@192.168.34.25:1622:xe",
  "dbuser": "rskadmin",
  "dbscript": "no"
 }
