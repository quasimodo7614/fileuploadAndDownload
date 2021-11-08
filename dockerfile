FROM busybox
ADD  file /

ENTRYPOINT ["/file"]

