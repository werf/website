package com.example.demo.data;

import org.springframework.data.repository.CrudRepository;

import java.util.List;

public interface LabelsRepository extends CrudRepository<Labels, Long> {
    List<Labels> findAll();
}
