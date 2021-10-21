FROM maven:3.6.0-jdk-11-slim AS build
RUN mvn -f clean package

FROM openjdk:15-jdk-alpine3.12
EXPOSE 8080
ADD /target/gateway.jar gateway.jar
ENTRYPOINT ["java","-jar","gateway.jar"]