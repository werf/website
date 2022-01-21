package com.flant.werfguidesapp.controllers;

import javax.servlet.http.HttpServletRequest;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;


@Controller
public class ImageController {

	private Logger log = LoggerFactory.getLogger(this.getClass());

	@GetMapping("/")
	public String MainPath(HttpServletRequest httpServletRequest) {
		log.info("Endpoint / request from the {} address", httpServletRequest.getRemoteAddr());
		return "index";
	}
	
	@GetMapping("/image")
	public String Image(HttpServletRequest httpServletRequest) {
		log.info("Endpoint /image request from the {} address", httpServletRequest.getRemoteAddr());
		return "image";
	}
}
