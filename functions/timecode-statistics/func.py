import fdk
import ujson
import pandas as pd
import numpy as np
import sys
import boto3
import os
import time
import logging

from botocore import exceptions

boto3.set_stream_logger('boto3.resources', logging.INFO)
TIMECODE = "timecode"
AMP_PEAK = "amplification_peak_value"

# default values are fore Minio setup for scripts/start.sh

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
    s3_conn_kwarg = dict(
        aws_access_key_id=os.environ.get("ACCESS_KEY_ID", "admin"),
        aws_secret_access_key=os.environ.get("SECRET_ACCESS_KEY", "password"),
        region_name=os.environ.get("REGION", "us-east-1"),
        endpoint_url="http://" + os.environ.get(
            "ENDPOINT", "docker.for.mac.localhost:9000")
    )
    s3_bucket = os.environ.get("BUCKET", "default-bucket")

    s3 = boto3.resource('s3', **s3_conn_kwarg)
    try:
        s3.Bucket(s3_bucket).download_file(
            "{0}.csv".format(key), '/tmp/{0}.csv'.format(key))
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
def from_file(path_or_url: str, amplification_peak: np.float64, stat_uuid: str):
    print("in from_file",
          file=sys.stderr, flush=True)
    df = pd.read_csv(
        path_or_url,
        names=[TIMECODE, AMP_PEAK],
        dtype={TIMECODE: np.float64, AMP_PEAK: np.float64}
    )
    print("datarame was provisioned",
          file=sys.stderr, flush=True)
    # filter out all levels below the amplification peak
    df = df.loc[df[AMP_PEAK] >= amplification_peak]
    df.timecode = df.timecode.round()
    timecode_series = df.groupby(TIMECODE).count()
    tc_sorted = timecode_series.sort_values(
        by=AMP_PEAK, ascending=False)
    top_peak_series_timecodes = tc_sorted.iloc[0:10].sort_values(
        by=TIMECODE, ascending=True)

    print("timecode DFT peak frequencies calculated",
          file=sys.stderr, flush=True)
    print(top_peak_series_timecodes, file=sys.stderr, flush=True)
    timecodes = [np.float64(0), ]
    timecodes.extend(
        row[0] for row in top_peak_series_timecodes.iterrows())
    timecodes.append(df[TIMECODE].max())
    print("timecodes: ", timecodes, file=sys.stderr, flush=True)
    res = []
    timcodes_len = len(timecodes)
    for index, timecode in enumerate(timecodes):
        next_index = index + 1
        # if we are at the end of the list -
        # nothing to calculate, break
        if next_index + 1 >= timcodes_len:
            break

        next_timecode = timecodes[next_index]
        # find a mean value for the timecode greater than 0
        mean_dft_for_frequent_timecode = df.loc[
            (df[TIMECODE].between(
                next_timecode - 1,
                next_timecode + 1,
                inclusive=True))
        ][AMP_PEAK].mean()

        before_range_dft = df.loc[
            (df[TIMECODE].between(
                timecode,
                next_timecode,
                inclusive=True)) &
            (df[AMP_PEAK].between(
                amplification_peak,
                mean_dft_for_frequent_timecode,
                inclusive=True))
        ]
        when_started = before_range_dft.loc[
            before_range_dft[AMP_PEAK] ==
            before_range_dft[AMP_PEAK].min()][TIMECODE].min()
        after_range_dft = df.loc[
            (df[TIMECODE].between(
                next_timecode,
                timecodes[next_index + 1],
                inclusive=True)) &
            (df[AMP_PEAK].between(
                amplification_peak,
                mean_dft_for_frequent_timecode,
                inclusive=True))
        ]
        when_finished = after_range_dft.loc[
            after_range_dft[AMP_PEAK] ==
            after_range_dft[AMP_PEAK].min()][TIMECODE].min()

        item = {
            "stat_uuid": stat_uuid,
            "begin": when_started,
            "end": when_finished,
            "mean_dft": mean_dft_for_frequent_timecode,
            "min_dft_before": before_range_dft[AMP_PEAK].min(),
            "min_dft_after": after_range_dft[AMP_PEAK].min(),
        }
        print("timecode stats: ", item, file=sys.stderr, flush=True)
        res.append(item)

    return res


def handler(ctx, data=None, loop=None):
    if data is not None:
        print(str(data), file=sys.stderr, flush=True)
        body = ujson.loads(data)
        peak = np.float64(body.get("threshold_value"))
        stat_uuid = body.get("stat_uuid")
        download_s3_file(stat_uuid)
        return from_file("/tmp/{0}.csv".format(
            stat_uuid), peak, stat_uuid),
    return {}


if __name__ == "__main__":
    fdk.handle(handler)
