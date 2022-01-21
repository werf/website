package com.flant.werfguidesapp.controllers;

import java.util.Optional;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import com.flant.werfguidesapp.domain.Talker;
import com.flant.werfguidesapp.repos.TalkersRepo;

@RestController
public class TalkerController {

	private Logger log = LoggerFactory.getLogger(this.getClass());
	
	@Autowired
	private TalkersRepo talkersRepo;

	@GetMapping("/remember")
	public String Remember(@RequestParam String answer, @RequestParam String name) {
		Talker talker = new Talker(answer, name);
		talkersRepo.save(talker);
		return "Got it.\n";
	}
	
	@GetMapping("/say")
	public String Say() {
		Optional<Talker> talker = talkersRepo.findById(1);
		if (talker.isPresent()) {
			return talker.get().getAnswer() + ", " + talker.get().getName() + "\n";
		} else return "I have nothing to say.\n";
	}

}
