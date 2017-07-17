#!/usr/bin/env python
# -*- coding: utf-8 -*-
import argparse
import traceback as tb
from env import config
import adm
import json
import error
import re

parser = argparse.ArgumentParser(description="")
parser.add_argument("--method", help="list_disks, update_disk, list_raids, create_raid, delete_raid", default='list_disks')
parser.add_argument("--argument", help="", default="")

def format(supported=['json']):
    def _format(func):
        def __format(*vargs, **kv):
            fmt = 'json'

            try:
                o = func(*vargs, **kv)
            except Exception as e:
                tb.print_exc()
                return json.dumps(error.error(e))
            else:
                if isinstance(o, (str, unicode)):
                    return o
                return json.dumps(o)
        return __format
    return _format

def _cmp_disk(a, b):
    a3 = re.findall('(\d+)\.(\d+)\.(\d+)', a.location)
    b3 = re.findall('(\d+)\.(\d+)\.(\d+)', b.location)
    if not a3 or not b3:
        return 0
    for x, y in zip(a3[0], b3[0]):
        if int(x) > int(y): return 1
        elif int(x) < int(y): return -1
    return 0
	
def _disk_as_obj(disk):
    obj = {'id': disk.uuid,
           'location' : disk.location or disk.prev_location,
           'vendor' : disk.vendor,
           'rpm' : disk.rpm,
           'health' : disk.health,
           'role' : disk.role,
           'cap_sector' : disk.cap.Sector,
           'cap_mb' : float(disk.cap.MB),
           'host': disk.host,
           'rqr_count': disk.rqr_count,
           'sn': disk.sn,
           'dev_name': disk.dev_name}
    if disk.raid and not disk.raid.deleted:
        obj['raid'] = disk.raid.name
    else:
        obj['raid'] = ''
    return obj

def _raid_as_obj(raid):
    pe = config.lvm.vg.pe_size
    obj = {'id'              : raid.uuid,
        'name'            : raid.name,
        'level'           : raid.level,
        'chunk_kb'        : raid.chunk.KB,
        'health'          : raid.health,
        'cap_sector'      : (raid.cap*pe).Sector,
        'cap_mb'          : float((raid.cap*pe).MB),
        'used_cap_sector' : (raid.used_cap*pe).Sector,
        'rebuilding'      : raid.rebuilding,
        'rebuild_progress': raid.rebuild_progress,
        'rqr_count'       : raid.rqr_count,
        'blkdev'          : '/dev/mapper/'+raid.odev_name,
    }
    return obj

def _input(argument, order):
    argmap = {}
    args = argument.split("&")
    for a in args:
        if a == "":continue
        kv = a.split("=")

        if kv[0] == "level":
            try:
                v = int(kv[1])
                argmap[kv[0]] = int(v)
            except:
                argmap[kv[0]] = 5
        elif kv[0] == "chunk":
            if kv[1] == "":
                argmap[kv[0]] = "64KB"
            else:
                argmap[kv[0]] = kv[1]
        else:
            argmap[kv[0]] = kv[1]

    res = []
    for k in order:
        if k in argmap.keys():
            res.append(argmap[k])
        else:
            res.append("")
    return res

@format()
def list_disks(arguments):
    try:
        disks = []
        diskmap = dict((disk.location, disk) for disk in adm.Disk.all() if disk.online and disk.health <> adm.HEALTH_DOWN)
        for disk in adm.Disk.all():
            if disk.raid and not disk.raid.deleted and disk.health == adm.HEALTH_DOWN:
                diskmap.setdefault(disk.prev_location, disk)

        disks = sorted(diskmap.values(), cmp=_cmp_disk)
    except:
            tb.print_exc()
    return error.success([_disk_as_obj(disk) for disk in disks])

@format()
def update_disk(arguments):
    args = _input(arguments, ["location", "host", "role", "raid"])
    location, host, role, raid = args[0], args[1], args[2], args[3]
    if host == adm.Disk.HOST_NATIVE:
        return adm.format_disk(location)
    elif role <> '':
        return adm.set_disk_role(location, role, raid)
    else:
        raise error.UnkownOperation()

@format()
def list_raids(arguments):
    return error.success([_raid_as_obj(raid) for raid in adm.Raid.all() if not raid.deleted])

@format()
def create_raid(argments):
    return adm.create_raid(*_input(argments, ["name", "level", "chunk", "raid_disks", "spare_disks", "rebuild_priority", "sync"]))

@format()
def delete_raid(argments):
    status = adm.delete_raid(*_input(argments, ["name"]))
    return status

if __name__ == "__main__":

    args = parser.parse_args()
    method = args.method
    argument = args.argument

    patmap = {"list_disks":  list_disks,
            "update_disk": update_disk,
            "list_raids":  list_raids,
            "create_raid": create_raid,
            "delete_raid": delete_raid}

    action = patmap.get(method)
    if action:
        print action(argument)
