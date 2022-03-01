package io.werf.werfguidesapp.repos;

import io.werf.werfguidesapp.domain.Talker;
import org.springframework.data.repository.CrudRepository;

public interface TalkersRepo extends CrudRepository<Talker, Integer> {

}
