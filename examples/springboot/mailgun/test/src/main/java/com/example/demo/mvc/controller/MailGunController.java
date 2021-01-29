package com.example.demo.mvc.controller;

import com.mashape.unirest.http.HttpResponse;
import com.mashape.unirest.http.JsonNode;
import com.mashape.unirest.http.Unirest;
import com.mashape.unirest.http.exceptions.UnirestException;
import org.springframework.boot.SpringApplication;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.Date;

@RestController
public class MailGunController {

    @GetMapping("/api/send_report")
    String sendReport() throws UnirestException {
        String REPORT_RECIEVER = System.getenv("REPORT_RECIEVER");
        String YOUR_DOMAIN_NAME = System.getenv("MAILGUN_DOMAIN");
        String API_KEY = System.getenv("MAILGUN_APIKEY");
        String date = String.valueOf(System.currentTimeMillis());
        HttpResponse<JsonNode> request = Unirest.post("https://api.mailgun.net/v3/" + YOUR_DOMAIN_NAME + "/messages").basicAuth("api", API_KEY)
                .queryString("from", "postmaster@"+ YOUR_DOMAIN_NAME)
                .queryString("to", REPORT_RECIEVER)
                .queryString("subject", "Report@"+ date)
                .queryString("text", date)
                .asJson();
        return "{\"result\": true}";
    }
}
