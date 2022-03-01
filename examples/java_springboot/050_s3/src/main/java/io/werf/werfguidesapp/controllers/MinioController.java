package io.werf.werfguidesapp.controllers;

import io.werf.werfguidesapp.components.MinioComponent;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

import java.io.ByteArrayInputStream;
import java.io.InputStream;

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
