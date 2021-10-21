FROM openjdk:15-jdk-alpine3.12
EXPOSE 8080
ADD /target/gateway.jar gateway.jar
ENTRYPOINT ["java","-jar","gateway.jar"]