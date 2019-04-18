import fdk
import json
import pandas as pd
import numpy as np
import sys
import boto3
import os
import time
import logging

from plotly import graph_objs
from plotly import offline

from botocore import exceptions


logging.getLogger('boto3').setLevel(logging.CRITICAL)
logging.getLogger('botocore').setLevel(logging.CRITICAL)
logging.getLogger('s3transfer').setLevel(logging.CRITICAL)
logging.getLogger('urllib3').setLevel(logging.CRITICAL)

TIMECODE = "timecode"
AMP_PEAK = "amplification_peak_value"
s3_conn_kwarg = dict(
    aws_access_key_id=os.environ.get("ACCESS_KEY_ID", "admin"),
    aws_secret_access_key=os.environ.get("SECRET_ACCESS_KEY", "password"),
    region_name=os.environ.get("REGION", "us-east-1"),
    endpoint_url="http://" + os.environ.get(
        "ENDPOINT", "docker.for.mac.localhost:9000")
)
s3_bucket = os.environ.get("BUCKET", "default-bucket")


def timeit(method):
    def timed(*args, **kw):
        ts = time.time()
        result = method(*args, **kw)
        te = time.time()
        print('{0} took {1} ms'.format(
            method.__name__, (te - ts) * 1000),
              file=sys.stderr, flush=True)
        return result

    return timed


@timeit
def download_s3_file(key):
    s3 = boto3.resource('s3', **s3_conn_kwarg)
    try:
        f = '/tmp/{0}.csv'.format(key)
        s3.Bucket(s3_bucket).download_file(
            "{0}.csv".format(key), f)
        return f
    except exceptions.ClientError as e:
        if e.response['Error']['Code'] == "404":
            print("The object does not exist.",
                  file=sys.stderr, flush=True)
        else:
            raise e


async def test_timecode_filtering(aiohttp_client):
    path_or_url = "/function/sox.csv"
    amp = 0.3
    res = from_file(path_or_url, np.float64(amp))
    print(res)
    assert res is not None


@timeit
def save_plot(x, y, stat_uuid: str):
    key = '{0}-peaks.html'.format(stat_uuid)
    html_file = '/tmp/{0}'.format(key)

    print("in save_plot", file=sys.stderr, flush=True)
    br = graph_objs.Bar(
        x=x, y=y,
        marker=dict(color='rgb(49,130,189)'),
        name="Peak frequency occasions per second"
    )
    print("bar chart provisioned", file=sys.stderr, flush=True)
    offline.plot(
        [br, ],
        filename=html_file,
        auto_open=False,
    )
    print("plot saved", file=sys.stderr, flush=True)
    return key, html_file


@timeit
def put_s3_file(key, filepath):
    s3 = boto3.resource('s3', **s3_conn_kwarg)
    s3.Bucket(s3_bucket).upload_file(filepath, key)


@timeit
def from_file(path_or_url: str, stat_uuid: str):
    peak = np.float64(0.4)
    print("in from_file",
          file=sys.stderr, flush=True)
    df = pd.read_csv(
        path_or_url,
        names=[TIMECODE, AMP_PEAK],
        dtype={TIMECODE: np.float64, AMP_PEAK: np.float64}
    )
    print("datarame was provisioned",
          file=sys.stderr, flush=True)
    df = df.loc[df[AMP_PEAK] >= peak]
    df.timecode = df.timecode.round()
    timecode_series = df.groupby(TIMECODE).count()
    timecode_series_sorted = timecode_series.sort_values(
        by=TIMECODE, ascending=False)

    print(timecode_series_sorted, file=sys.stderr, flush=True)
    # todo: find better way to get a list of ints for TS values
    ys = []
    for y in timecode_series_sorted.values:
        ys.append(y[0])
    print(ys, file=sys.stderr, flush=True)
    return save_plot(
        timecode_series_sorted.index,
        ys,
        stat_uuid,
    )


def handler(ctx, data=None, loop=None):
    if data is not None:
        body = json.loads(data.getvalue())
        stat_uuid = body.get("stat_uuid")
        html_file = download_s3_file(stat_uuid)

        print(html_file, file=sys.stderr, flush=True)

        key, html_file = from_file(
            html_file, stat_uuid
        )

        print(key, file=sys.stderr, flush=True)
        put_s3_file(key, html_file)

        return body

    return data
