package com.flant.werfguidesapp.controllers;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;


@Controller
public class ImageController {

	@GetMapping("/")
	public String MainPath() {
		return "index";
	}
	
	@GetMapping("/image")
	public String Image() {
		return "image";
	}
}
