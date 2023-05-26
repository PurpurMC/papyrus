package org.purpurmc.papyrus;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.data.jpa.repository.config.EnableJpaRepositories;

@EnableJpaRepositories("org.purpurmc.papyrus.db.repository")
@SpringBootApplication
public class PapyrusApplication {
    public static void main(String[] args) {
        SpringApplication.run(PapyrusApplication.class, args);
    }
}
