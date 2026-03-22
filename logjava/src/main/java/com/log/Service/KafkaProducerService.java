package com.log.Service;




import com.log.Model.LogRequest;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

@Service
public class KafkaProducerService {

    private final KafkaTemplate<String, LogRequest> kafkaTemplate;

    private static final String TOPIC = "logs-topic";

    public KafkaProducerService(KafkaTemplate<String, LogRequest> kafkaTemplate) {
        this.kafkaTemplate = kafkaTemplate;
    }

    public void sendLog(LogRequest logRequest) {
        kafkaTemplate.send(TOPIC, logRequest.getService(), logRequest);
    }
}