package com.flant.werfguidesapp.repos;

import org.springframework.data.repository.CrudRepository;

import com.flant.werfguidesapp.domain.Talker;

public interface TalkersRepo extends CrudRepository<Talker, Integer> {

}
