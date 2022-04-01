Clean restart of docker instance: https://docs.tibco.com/pub/mash-local/4.1.0/doc/html/docker/GUID-BD850566-5B79-4915-987E-430FC38DAAE4.html
Access docker volumes: https://github.com/docker/for-mac/issues/4822

Commands to start up docker instance:

In root directory containing docker-compose.yml file:

1. docker-compose down
2. docker-compose -f docker-compose.yml up