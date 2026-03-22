package com.log.Controller;

import com.log.Model.LogRequest;
import com.log.Service.LogService;
import jakarta.validation.Valid;
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
}