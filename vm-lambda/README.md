<!--
title: .'VM Housekeeping Lambdas'
description: 'Code to GET VM/EC2 list and POST VM schedule submission'
framework: v1
platform: AWS
language: Go
-->
 

Following endpoints provided:
1. GET endpoint with type parameter (/get/{type}) --> Type in (onpremise/aws)
   Sample Response:
   [{"instanceId":"i-0092445d5b636238a","name":"arun-linux"}]