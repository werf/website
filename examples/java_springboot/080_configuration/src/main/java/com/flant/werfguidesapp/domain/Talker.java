package com.flant.werfguidesapp.domain;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import javax.persistence.JoinColumn;
import javax.persistence.Table;


@Entity
@Table(name = "talkers")
public class Talker {

	@Id
	@GeneratedValue(strategy = GenerationType.AUTO)
	private Integer id;
	
	@JoinColumn(name = "answer")
	private String answer;
	
	@JoinColumn(name = "name")
	private String name;

	public Integer getId() {
		return id;
	}

	public Talker() {

	}

	public Talker(String answer, String name) {
		this.answer = answer;
		this.name = name;
	}

	public void setId(Integer id) {
		this.id = id;
	}

	public String getAnswer() {
		return answer;
	}

	public void setAnswer(String answer) {
		this.answer = answer;
	}

	public String getName() {
		return name;
	}

	public void setName(String name) {
		this.name = name;
	}
}
