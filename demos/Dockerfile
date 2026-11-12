FROM ghcr.io/charmbracelet/vhs

RUN apt-get update \
  && apt-get install -y --no-install-recommends vim \
  && rm -rf /var/lib/apt/lists/*

COPY mxs /usr/local/bin/mxs
