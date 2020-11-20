package com.example.demo.mvc.controller;

import com.example.demo.data.Labels;
import com.example.demo.data.LabelsRepository;
import com.example.demo.mvc.exception.LabelError;
import com.example.demo.mvc.exception.LabelException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api")
public class LabelController {
    @Autowired
    LabelsRepository repository;

    @ExceptionHandler({LabelException.class})
    public LabelError handleException(LabelException ex) {
        return new LabelError(ex.getError());
    }

    @GetMapping("/labels")
    public List<Labels> labels() {
        return repository.findAll();
    }

    @GetMapping("/labels/{id}")
    public Labels labels(@PathVariable Long id) throws LabelException {
        return repository.findById(id).orElseThrow(() -> new LabelException("not found"));
    }

    @PostMapping("/labels")
    Labels labels(@RequestBody Labels newLabels) {
        return repository.save(newLabels);
    }


    @PostMapping("/labels/{id}")
    Labels labels(@RequestBody Labels newLabels, @PathVariable Long id) throws LabelException {
        return repository.findById(id)
                .map(label -> {
                    label.setLabel(newLabels.getLabel());
                    return repository.save(label);
                })
                .orElseThrow(() -> new LabelException("not found"));
    }

    @DeleteMapping("/labels/{id}")
    String labelsDelete(@PathVariable Long id) throws LabelException {
        try {
            repository.deleteById(id);
        }
        catch (Exception e){
            throw new LabelException("not found");
        }
        return "{\\\"result\\\": true}";
        }
}
