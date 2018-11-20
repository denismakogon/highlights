FROM fnproject/python:3.7-dev as build-stage

WORKDIR /function
ADD requirements.txt /function/
RUN pip3 install --target /python/  --no-cache --no-cache-dir -r requirements.txt pytest
ADD . /function/

FROM fnproject/python:3.7
WORKDIR /function
COPY --from=build-stage /function /function
COPY --from=build-stage /python /python
ENV PYTHONPATH=/python
ENTRYPOINT ["python3", "func.py"]
