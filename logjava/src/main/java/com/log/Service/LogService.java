package com.log.Service;

import com.log.Model.LogRequest;
import org.springframework.stereotype.Service;

import java.util.HashMap;
import java.util.Map;

@Service
public class LogService {

    private final KafkaProducerService kafkaProducerService;

    public LogService(KafkaProducerService kafkaProducerService) {
        this.kafkaProducerService = kafkaProducerService;
    }

    public void processLog(LogRequest logRequest) {

        // Add timestamp if missing
        if (logRequest.getTimestamp() == null) {
            logRequest.setTimestamp(System.currentTimeMillis());
        }

        kafkaProducerService.sendLog(logRequest);
    }

    public void processRawLog(String rawLog, String source) {
        if (rawLog == null || rawLog.isBlank()) {
            return;
        }

        String normalizedSource = (source == null || source.isBlank()) ? "unknown" : source;

        Map<String, Object> payload = new HashMap<>();
        payload.put("raw_log", rawLog);
        payload.put("source", normalizedSource);
        payload.put("timestamp", System.currentTimeMillis());

        kafkaProducerService.sendRawLog(normalizedSource, payload);
    }
}
