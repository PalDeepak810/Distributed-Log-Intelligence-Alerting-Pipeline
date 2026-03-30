package com.log.Service;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.log.Model.LogRequest;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

import java.util.Map;

@Service
public class KafkaProducerService {

    private final KafkaTemplate<String, String> kafkaTemplate;
    private static final String TOPIC = "logs-topic";
    private final ObjectMapper objectMapper = new ObjectMapper();

    public KafkaProducerService(KafkaTemplate<String, String> kafkaTemplate) {
        this.kafkaTemplate = kafkaTemplate;
    }

    public void sendLog(LogRequest logRequest) {
        try {
            String json = objectMapper.writeValueAsString(logRequest);
            kafkaTemplate.send(TOPIC, logRequest.getService(), json);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    public void sendRawLog(String key, Map<String, Object> rawPayload) {
        try {
            String json = objectMapper.writeValueAsString(rawPayload);
            kafkaTemplate.send(TOPIC, key, json);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}