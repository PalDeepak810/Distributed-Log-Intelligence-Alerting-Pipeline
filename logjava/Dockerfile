# BUILD STAGE
FROM maven:3.9.9-eclipse-temurin-21 AS build

WORKDIR /app

COPY . .

RUN mvn clean package -DskipTests

#RUNTIME STAGE
FROM eclipse-temurin:21-jdk-jammy

WORKDIR /app

COPY --from=build /app/target/*.jar app.jar

COPY src/main/resources/ca.pem /app/ca.pem
COPY src/main/resources/keystore.p12 /app/keystore.p12

EXPOSE 8080

ENTRYPOINT ["java", "-jar", "app.jar"]