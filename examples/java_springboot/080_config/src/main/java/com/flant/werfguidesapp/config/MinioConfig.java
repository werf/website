package com.flant.werfguidesapp.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import io.minio.MinioClient;

@Configuration
public class MinioConfig {

	@Value("${minio.url}")
	private String minioUrl;
	@Value("${minio.accessKey}")
	private String accessKey;
	@Value("${minio.secretKey}")
	private String secretKey;
	@Value("${minio.bucket-name}")
	private String bucketName;
	@Value("${minio.secure}")
	private Boolean minioSecure;

	@Bean
	public MinioClient minioClient() {
		MinioClient minioClient = MinioClient.builder().credentials(accessKey, secretKey)
				.endpoint(minioUrl, 9000, minioSecure).build();
		return minioClient;
	}
}
