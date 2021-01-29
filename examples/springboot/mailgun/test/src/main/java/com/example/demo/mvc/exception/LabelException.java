package com.example.demo.mvc.exception;

public class LabelException extends Exception {
    final String error;

    public LabelException(String error) {
        this.error = error;
    }

    public String getError() {
        return error;
    }
}
