FROM ubuntu:22.04

RUN apt update && apt install -y locales && rm -rf /var/lib/apt/lists/*
RUN locale-gen en_US.UTF-8

ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en_US.UTF-8 \
    LC_COLLATE=en_US.UTF-8
