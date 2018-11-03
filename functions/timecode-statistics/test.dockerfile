FROM denismakogon/ffmpeg-debian:runtime as test-input-stage

# for the test purposes
RUN apt-get update \
  && apt-get install -qy sox libsox-fmt-all
#ADD test_data/sample_video.m4v /tmp/sample_video.m4v
RUN ffmpeg -loglevel panic -y -i http://mirrors.standaloneinstaller.com/video-sample/metaxas-keller-Bell.m4v -vn -f wav - | sox -t wav - /tmp/sox.dat
RUN cat /tmp/sox.dat | awk '{ if (NR>2) print $1","$2}' > /sox.csv

FROM fnproject/python:3.7-dev as build-stage

COPY --from=test-input-stage /sox.csv /function/sox.csv
WORKDIR /function
ADD requirements.txt /function/
RUN rm -fr __pycache__ && pip3 install --no-cache --no-cache-dir -r requirements.txt
RUN pip3 install --target /python/  --no-cache --no-cache-dir -r requirements.txt pytest
ADD func.py /function/
RUN pytest -v -s --tb=long func.py

FROM fnproject/python:3.7
WORKDIR /function
COPY --from=build-stage /function /function
RUN rm -fr /function/sox.csv
COPY --from=build-stage /python /python
ENV PYTHONPATH=/python
ENTRYPOINT ["python3", "func.py"]
