# Docker image for the docker plugin
#
#     docker build --rm=true -t ivancevich/drone-zipper .

FROM gliderlabs/alpine

ADD drone-zipper /bin/
ENTRYPOINT ["/bin/drone-zipper"]
