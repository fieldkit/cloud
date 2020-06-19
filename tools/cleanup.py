#!/usr/bin/python3

import sys
import logging
import boto3


def get_all_s3_keys(bucket):
    keys = []

    s3 = boto3.client("s3")

    kwargs = {"Bucket": bucket}
    while True:
        resp = s3.list_objects_v2(**kwargs)
        for obj in resp["Contents"]:
            keys.append(obj["Key"])

        try:
            kwargs["ContinuationToken"] = resp["NextContinuationToken"]
        except KeyError:
            break

    return keys


def get_s3_keys_as_generator(bucket):
    s3 = boto3.client("s3")

    kwargs = {"Bucket": bucket}
    while True:
        resp = s3.list_objects_v2(**kwargs)
        for obj in resp["Contents"]:
            yield obj["Key"]

        try:
            kwargs["ContinuationToken"] = resp["NextContinuationToken"]
        except KeyError:
            break


if __name__ == "__main__":
    logging.basicConfig(stream=sys.stdout, level=logging.INFO)
    log = logging.getLogger("cleanup")

    for key in get_s3_keys_as_generator("fk-streams"):
        print(key)
