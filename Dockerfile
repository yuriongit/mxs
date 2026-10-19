FROM goreleaser/goreleaser:v2.18.3

ENV GOTOOLCHAIN=auto \
  HOME=/tmp/home \
  GOCACHE=/tmp/go-build \
  GOMODCACHE=/tmp/go-mod

RUN mkdir -p /tmp/home /tmp/go-build /tmp/go-mod \
  && chmod 1777 /tmp/home /tmp/go-build /tmp/go-mod

WORKDIR /app

# --snapshot skips the git tag requirement and doesn't publish anything
CMD ["release", "--snapshot", "--clean"]
