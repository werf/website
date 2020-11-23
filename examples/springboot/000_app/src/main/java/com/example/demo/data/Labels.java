package com.example.demo.data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.Id;

@Entity
public class Labels {
    @Id @GeneratedValue
    private long id;
    private String label;

    public Labels() {

    }

    public void setId(int id) {
        this.id = id;
    }

    public long getId() {
        return id;
    }

    public String getLabel() {
        return label;
    }

    public void setLabel(String label) {
        this.label = label;
    }

    public Labels(String label){
        this.label = label;
    }
}
