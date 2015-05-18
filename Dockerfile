FROM ubuntu:14.04

EXPOSE 3000

ADD frontend /frontend
ADD public /public
ADD templates /templates

CMD ["/frontend"]
