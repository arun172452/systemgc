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
  "ipaddress": "192.168.2.1",
  "agent": "no",
  "naspath": "/test/mount,/myinterface/mount"
 }
 
 {
  "userid": "arunl",
  "ipaddress": "192.168.2.1",
  "agent": "yes",
  "naspath": "/test/mount,/myinterface/mount"
 }
