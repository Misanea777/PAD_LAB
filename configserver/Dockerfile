FROM openjdk:15-jdk-alpine3.12
EXPOSE 8888
ADD /target/configuration-server.jar configuration-server.jar
ENTRYPOINT ["java","-jar","configuration-server.jar"]