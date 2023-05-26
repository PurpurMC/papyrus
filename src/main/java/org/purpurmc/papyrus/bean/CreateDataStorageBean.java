package org.purpurmc.papyrus.bean;

import jakarta.annotation.PostConstruct;
import org.purpurmc.papyrus.config.AppConfiguration;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;

@Component
public class CreateDataStorageBean {
    private final AppConfiguration configuration;

    @Autowired
    public CreateDataStorageBean(AppConfiguration configuration) {
        this.configuration = configuration;
    }

    @PostConstruct
    public void run() throws IOException {
        Path path = Path.of(configuration.getFileStorage());
        Files.createDirectories(path);
    }
}
