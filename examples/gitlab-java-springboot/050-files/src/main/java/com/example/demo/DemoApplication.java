package com.example.demo;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import java.io.File;
import java.io.FileWriter;
import java.nio.file.Files;
import com.amazonaws.auth.AWSStaticCredentialsProvider;
import java.util.List;
import com.amazonaws.auth.AWSCredentials;
import com.amazonaws.auth.BasicAWSCredentials;
import com.amazonaws.services.s3.AmazonS3;
import com.amazonaws.services.s3.AmazonS3ClientBuilder;
import com.amazonaws.client.builder.AwsClientBuilder.EndpointConfiguration;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.env.Environment;

@SpringBootApplication
@RestController
public class DemoApplication {

    @Autowired
    private Environment env;

    @RequestMapping("/")
    public String home() {
        try {
            String bucketName  = env.getProperty("amazonProperties.bucketName");
            String login       = env.getProperty("amazonProperties.login");
            String password    = env.getProperty("amazonProperties.password");
            String endpointUrl = env.getProperty("amazonProperties.endpointUrl");
            String zone        = env.getProperty("amazonProperties.zone");
            AWSCredentials credentials = new BasicAWSCredentials(login, password);
            AmazonS3 s3client = AmazonS3ClientBuilder.standard()
                    .withCredentials(new AWSStaticCredentialsProvider(credentials))
                    .withPathStyleAccessEnabled(true)
                    .withEndpointConfiguration((new EndpointConfiguration(endpointUrl, zone)))
                    .build();

            if(!s3client.doesBucketExist(bucketName)) {
                s3client.createBucket(bucketName);
            }

            File file = File.createTempFile("temp", null);
            FileWriter writer = new FileWriter(file);
            writer.write("Test data");
            writer.close();
            s3client.putObject(
                    bucketName,
                    "hello.txt",
                    file
            );
            file.deleteOnExit();
        }catch(Exception e) {
            e.printStackTrace();
        }
        return "Hello World";
    }

    public static void main(String[] args) {
        SpringApplication.run(DemoApplication.class, args);
    }

}

