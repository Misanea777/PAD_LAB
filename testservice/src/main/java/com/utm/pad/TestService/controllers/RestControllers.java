package com.utm.pad.TestService.controllers;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.cloud.context.config.annotation.RefreshScope;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api")
public class RestControllers {

        @Value("${eureka.instance.instance-id}")
        private String name;

        @GetMapping
        public String test() { return name; }

        @GetMapping("/test")
        public String testing() { return name; }
}
