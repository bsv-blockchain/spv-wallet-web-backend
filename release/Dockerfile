# syntax=docker/dockerfile:1

FROM --platform=$TARGETPLATFORM alpine:3.24@sha256:be2656061b536c1883c4034790e21d0a09c5aaaf51968b89368a3c0371ea2c9a

WORKDIR /app

ENV DB_MIGRATIONSPATH=./migrations
COPY ./data/sql/migrations $DB_MIGRATIONSPATH

COPY --chmod=744 ./release/entrypoint.sh .

ARG version
ENV VERSION=${version:-develop}
ENV TAG=${tag:-master}

ARG project_name
ENV APP_BINARY=${project_name}
ARG TARGETPLATFORM
COPY --chmod=744 ./dist/$TARGETPLATFORM/$APP_BINARY .

ENTRYPOINT ["./entrypoint.sh"]
