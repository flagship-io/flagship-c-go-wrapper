FROM neilotoole/xcgo:latest

ADD . .

RUN sudo su

CMD [ "./entrypoint.sh" ]