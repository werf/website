package com.example.demo.mvc.controller;

import com.amazonaws.auth.AWSCredentials;
import com.amazonaws.auth.AWSStaticCredentialsProvider;
import com.amazonaws.auth.BasicAWSCredentials;
import com.amazonaws.client.builder.AwsClientBuilder;
import com.amazonaws.services.s3.AmazonS3;
import com.amazonaws.services.s3.AmazonS3ClientBuilder;
import com.amazonaws.services.s3.model.ObjectListing;
import com.amazonaws.services.s3.model.S3ObjectSummary;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(value = "/api")
public class S3Controller {
    AmazonS3 s3client;

    S3Controller() {
        AWSCredentials credentials = new BasicAWSCredentials(
                System.getenv("S3_ACCESS_KEY"),
                System.getenv("S3_SECRET_KEY")
        );
        this.s3client = AmazonS3ClientBuilder
                .standard()
                .withEndpointConfiguration(new AwsClientBuilder.EndpointConfiguration(System.getenv("S3_ENDPOINT"), System.getenv("S3_ZONE")))
                .withCredentials(new AWSStaticCredentialsProvider(credentials))
                .withPathStyleAccessEnabled(true)
                .build();
    }

    @GetMapping("/generate_report")
    public String report(){
        long date = System.currentTimeMillis();
        s3client.putObject(
                System.getenv("S3_BUCKET"),
                String.format("report_%o_.txt", date),
                String.valueOf(date));
        ObjectListing objectListing = s3client.listObjects(System.getenv("S3_BUCKET"));
        return "{\"result\": true}";
    }
}
