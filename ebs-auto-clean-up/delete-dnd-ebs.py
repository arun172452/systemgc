#!/usr/bin/env python
import boto3
ec2 = boto3.resource('ec2',region_name='eu-central-1')

def lambda_handler(event, context):
    for vol in ec2.volumes.all():
        if  vol.state=='available':
            if vol.tags is None:
                vid=vol.id
				v=ec2.Volume(vol.id)
                v.delete()
                print "Deleted " +vid
                continue
            for tag in vol.tags:
                if tag['Key'] == 'Autodelete':
                    value=tag['Value']
                    if value != 'NO' and vol.state=='available':
                        vid=vol.id
                        v=ec2.Volume(vol.id)
                        v.delete()
                        print "Deleted " +vid