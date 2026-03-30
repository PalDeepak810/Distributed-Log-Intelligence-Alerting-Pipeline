package com.log;

import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.core.env.Environment;

@SpringBootApplication
public class Application {

	public static void main(String[] args) {
		SpringApplication.run(Application.class, args);
	}
	@Bean
	public CommandLineRunner testKafkaConfig(Environment env) {
		return args -> {
			System.out.println("BOOTSTRAP SERVERS: " + env.getProperty("spring.kafka.bootstrap-servers"));
		};
	}
}
