package com.flant.werfguidesapp.controllers;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;


@RestController
public class PingController {

	private Logger log = LoggerFactory.getLogger(this.getClass());

	@GetMapping("/ping")
	public String Ping() {
		return "pong\n";
	}
}
