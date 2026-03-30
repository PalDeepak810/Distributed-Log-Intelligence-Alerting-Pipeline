package com.log.Controller;

import com.log.Model.LogRequest;
import com.log.Service.LogService;
import jakarta.validation.Valid;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/logs")
public class LogController {

    private final LogService logService;

    public LogController(LogService logService) {
        this.logService = logService;
    }

    @PostMapping
    public String ingestLog(@Valid @RequestBody LogRequest logRequest) {
        logService.processLog(logRequest);
        return "Log received";
    }

    @PostMapping(value = "/raw", consumes = MediaType.TEXT_PLAIN_VALUE)
    public String ingestRawLog(
            @RequestBody String rawLog,
            @RequestParam(name = "source", defaultValue = "unknown") String source
    ) {
        logService.processRawLog(rawLog, source);
        return "Raw log received";
    }
}
