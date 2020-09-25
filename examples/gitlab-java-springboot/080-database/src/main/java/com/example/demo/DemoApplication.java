package com.example.demo;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.env.Environment;
import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;

@SpringBootApplication
@RestController
public class DemoApplication {

	@Autowired
	private Environment env;

	private Boolean postgreLoaded = Boolean.FALSE;
	private Statement st;

	private void loadPostgre() throws SQLException {
		if (! postgreLoaded)
		{
			String url = "jdbc:postgresql://"+env.getProperty("postgresql.host")+":"+env.getProperty("postgresql.port")+"/"+env.getProperty("postgresql.database");
			String user = env.getProperty("postgresql.login");
			String password = env.getProperty("postgresql.password");
			Connection con = DriverManager.getConnection(url, user, password);
		 	st = con.createStatement();

			postgreLoaded = Boolean.TRUE;
		}
	}

	@RequestMapping("/")
	public String home() {
		try {
			loadPostgre();
			ResultSet rs = st.executeQuery("select 1+2;");
			String myval = "[not loaded]";
			if (rs.next()) {
				myval = rs.getString(1);
			}
			return "Hello World, I have value " + myval;
		} catch (SQLException ex) {
			ex.printStackTrace();
			return "Exception: " + ex.getMessage();
		}
	}

	public static void main(String[] args) {
		SpringApplication.run(DemoApplication.class, args);
	}

}
