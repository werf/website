package com.example.demo.mvc.exception;

public class LabelError {
    final String result = "eror";
    final String comment;

    public LabelError(String comment) {
        this.comment = comment;
    }

    public String getResult() {
        return result;
    }

    public String getComment() {
        return comment;
    }
}
