package com.log.Model;

import jakarta.validation.constraints.NotBlank;
import lombok.Data;

@Data
public class LogRequest {

    @NotBlank
    private String service;

    @NotBlank
    private String level;

    @NotBlank
    private String message;

    private Long timestamp;

    private String traceId;

    public LogRequest(String service, String level, String message, Long timestamp, String traceId) {
        this.service = service;
        this.level = level;
        this.message = message;
        this.timestamp = timestamp;
        this.traceId = traceId;
    }

    public LogRequest() {
    }

    public String getService() {
        return service;
    }

    public void setService(String service) {
        this.service = service;
    }

    public String getLevel() {
        return level;
    }

    public void setLevel(String level) {
        this.level = level;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public Long getTimestamp() {
        return timestamp;
    }

    public void setTimestamp(Long timestamp) {
        this.timestamp = timestamp;
    }

    public String getTraceId() {
        return traceId;
    }

    public void setTraceId(String traceId) {
        this.traceId = traceId;
    }
}