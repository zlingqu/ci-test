FROM alpine
WORKDIR /app
ADD targets/ci-test .
cmd ./ci-test