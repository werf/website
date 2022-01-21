package com.flant.werfguidesapp.controllers;

import java.io.ByteArrayInputStream;
import java.io.File;
import java.io.InputStream;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import com.flant.werfguidesapp.components.MinioComponent;

@RestController
public class MinioController {

	@Autowired
	private MinioComponent minioComponent;

	public MinioController() {

	}

	@PostMapping("/upload")
	public String uploadFileToMinIO(@RequestParam("file") MultipartFile file) {
		try {
			InputStream in = new ByteArrayInputStream(file.getBytes());
			String fileName = file.getOriginalFilename();
			minioComponent.putObject(fileName, in);
			return "File uploaded.";
		} catch (Exception e) {
			e.printStackTrace();
		}
		return "Something wrong.";
	}

	@GetMapping("/download")
	public String downloadFile() throws Exception {
		return minioComponent.getObject("file.txt");
	}
}
