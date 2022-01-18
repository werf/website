package com.flant.werfguidesapp.controllers;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import javax.servlet.http.HttpServletRequest;


@RestController
public class PingController {

	@GetMapping("/ping")
	public String Ping(HttpServletRequest httpServletRequest) {
		return "pong\n";
	}
}
