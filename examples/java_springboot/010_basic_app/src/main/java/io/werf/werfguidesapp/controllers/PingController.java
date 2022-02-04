package io.werf.werfguidesapp.controllers;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class PingController {
	
	@GetMapping("/ping")
	public String Ping() {
		return "pong\n";
	}
}
