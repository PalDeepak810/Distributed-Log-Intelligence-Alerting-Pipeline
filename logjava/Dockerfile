FROM eclipse-temurin:21-jdk-jammy

WORKDIR /app

COPY target/demo-0.0.1-SNAPSHOT.jar  app.jar

COPY src/main/resources/ca.pem /app/ca.pem
COPY src/main/resources/keystore.p12 /app/keystore.p12

EXPOSE 8080

ENTRYPOINT ["java", "-jar", "app.jar"]