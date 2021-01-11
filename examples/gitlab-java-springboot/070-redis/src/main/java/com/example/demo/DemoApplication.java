package com.example.demo;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.env.Environment;
import redis.clients.jedis.*;

@SpringBootApplication
@RestController
public class DemoApplication {

	@Autowired
	private Environment env;

	private Jedis _jedis;
	private Boolean redisLoaded = Boolean.FALSE;

	private void loadRedis(){
		if (! redisLoaded)
		{
			String redisHost  = env.getProperty("redis.host");
			String redisPort  = env.getProperty("redis.port");
			String redisLogin  = env.getProperty("redis.login");
			String redisPassword  = env.getProperty("redis.password");
			JedisShardInfo shardInfo = new JedisShardInfo(redisHost, redisPort);
			shardInfo.setPassword(redisPassword);
			_jedis = new Jedis(shardInfo);

			redisLoaded = Boolean.TRUE;
		}
	}

	@RequestMapping("/")
	public String home() {
		loadRedis();
		String myval = _jedis.get("counter");
		return "Hello World, I have value " + myval + "<br><a href=\"/randomset\">increment it</a>";
	}

	@RequestMapping("/randomset")
	public String increment() {
		loadRedis();
		_jedis.set("counter", "bar");
		String myval = _jedis.get("counter");
		return "Incremented, now value "+myval+"<br><a href=\"/\">see root page</a>";
	}


	public static void main(String[] args) {
		SpringApplication.run(DemoApplication.class, args);
	}

}
