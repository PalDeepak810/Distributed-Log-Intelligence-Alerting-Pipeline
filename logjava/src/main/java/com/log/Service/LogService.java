package com.log.Service;

import com.log.Model.LogRequest;
import org.springframework.stereotype.Service;

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
}