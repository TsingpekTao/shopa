FROM loads/alpine:3.8

###############################################################################
#                                INSTALLATION
###############################################################################

ENV WORKDIR /app
# Default to production config file; can be overridden at runtime.
ENV GF_GCFG_FILE config.prod.yaml

ADD resource $WORKDIR/
ADD manifest/config $WORKDIR/manifest/config
ADD ./temp/linux_amd64/main $WORKDIR/main
RUN chmod +x $WORKDIR/main

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD ./main
