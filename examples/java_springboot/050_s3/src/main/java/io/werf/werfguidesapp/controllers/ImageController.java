package io.werf.werfguidesapp.controllers;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;

import javax.servlet.http.HttpServletRequest;


@Controller
public class ImageController {

	@GetMapping("/")
	public String MainPath(HttpServletRequest httpServletRequest) {
		return "index";
	}
	
	@GetMapping("/image")
	public String Image(HttpServletRequest httpServletRequest) {
		return "image";
	}
}
