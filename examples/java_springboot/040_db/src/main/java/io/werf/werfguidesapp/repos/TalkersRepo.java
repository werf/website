package io.werf.werfguidesapp.repos;

import org.springframework.data.repository.CrudRepository;

import io.werf.werfguidesapp.domain.Talker;

public interface TalkersRepo extends CrudRepository<Talker, Integer> {

}
