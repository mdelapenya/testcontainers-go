# This file won't be build without BuildKit support enabled
FROM alpine as base

ARG FILENAME

RUN echo "test" >> $FILENAME

FROM base as runner

RUN test $FILENAME