import fdk
import ujson
import time
import boto3
import os
import logging
import sys
import pandas as pd
import numpy as np
import plotly.graph_objs as go

from plotly.offline import plot
from botocore import exceptions

boto3.set_stream_logger('boto3.resources', logging.INFO)

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


async def test_timecode_filtering():

    p = plot_dataframe(
        "test", "/function/sox.csv",
        np.float64(479.0),
        np.float64( 889.0)
    )
    assert True == os.path.exists(p)


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
def put_s3_file(key, filepath):
    s3 = boto3.resource('s3', **s3_conn_kwarg)
    s3.Bucket(s3_bucket).upload_file(filepath, "{0}.html".format(key))


@timeit
def download_s3_file(key):
    s3 = boto3.resource('s3', **s3_conn_kwarg)
    try:
        s3.Bucket(s3_bucket).download_file(
            "{0}.csv".format(key), '/tmp/{0}.csv'.format(key))
        return '/tmp/{0}.csv'.format(key)
    except exceptions.ClientError as e:
        if e.response['Error']['Code'] == "404":
            print("The object does not exist.",
                  file=sys.stderr, flush=True)
        raise e

@timeit
def load_dataframe(filepath):
    return pd.read_csv(
        filepath,
        names=[TIMECODE, AMP_PEAK],
        dtype={TIMECODE: np.float64, AMP_PEAK: np.float64}
    )

@timeit
def filter_df(df, begin, end):
    return df.loc[df[TIMECODE].between(
        np.float64(begin), np.float64(end), inclusive=True)]

@timeit
def save_plot(df, html_path):
    print("in save_plot", file=sys.stderr, flush=True)
    sc = go.Scatter(x=df[TIMECODE], y=df[AMP_PEAK])
    print("scatter provisioned", file=sys.stderr, flush=True)
    plot([sc, ], filename=html_path, auto_open=False,)

@timeit
def plot_dataframe(df, key, begin, end):
    df = filter_df(df, begin, end)
    print("filtering completed", file=sys.stderr, flush=True)
    html_path = '/tmp/{0}.html'.format(key)
    save_plot(df, html_path)
    print("html file was created: ",
          html_path, file=sys.stderr, flush=True)
    return html_path


@timeit
def handler(ctx, data=None, loop=None):
    if data is not None:
        print("body: ", data, file=sys.stderr, flush=True)
        body = ujson.loads(data)
        key = body.get("stat_uuid")
        filepath = download_s3_file(key)
        df = load_dataframe(filepath)
        begin = body.get("begin")
        end = body.get("end")
        html_path = plot_dataframe(df, key, begin, end)
        html_key = 'range-{1}-{2}-{0}.html'.format(
            key, int(begin), int(end))
        put_s3_file(html_key, html_path)
    return


if __name__ == "__main__":
    fdk.handle(handler)
