<!--
title: .'VM Housekeeping Lambdas'
description: 'Code to GET VM/EC2 list and POST VM schedule submission'Cancel changes
framework: v1
platform: AWS
language: Go
-->
 

Following endpoints provided:
1. GET endpoint with type parameter (/get/{type}) --> Type in (onpremise/aws)
   Sample Response:
   [{"instanceId":"i-0092445d5b636238a","name":"arun-linux"}]
2. Post endpoint  /dev/registry

Sample input Json:

{
  "userid": "arunl",
  "vmname": "ukpsrv2",
  "vmtype": "onpremise",
  "instanceid": "192.168.2.1",
  "ipaddress": "192.168.2.1",
  "sancleanup": "NA"
 }
 
 {
  "userid": "arunl",
  "vmname": "arun-linux",
  "vmtype": "aws",
  "instanceid": "i-23456ffggss",
  "ipaddress": "NA",
  "sancleanup": "NA"
 }
