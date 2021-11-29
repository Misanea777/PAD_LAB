package com.utm.pad.Gateway.controllers;


import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class Endpoints {
    @GetMapping("/welcome")
    public String welcome() {
        return "Sorry service is temporarily unavailable";
    }
}
