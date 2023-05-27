package org.purpurmc.papyrus.controller;

import io.swagger.v3.oas.annotations.Hidden;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.net.URI;

@Hidden
@RestController
public class RootController {

    @GetMapping({"/", "/docs"})
    public ResponseEntity<?> documentationRedirect() {
        return ResponseEntity.status(HttpStatus.FOUND)
                .location(URI.create("/docs/"))
                .build();
    }
}
