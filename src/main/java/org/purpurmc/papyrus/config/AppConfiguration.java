package org.purpurmc.papyrus.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

@Configuration
@ConfigurationProperties("app")
public class AppConfiguration {
    private String fileStorage;
    private String authToken;

    public String getFileStorage() {
        return this.fileStorage;
    }

    public String getAuthToken() {
        return this.authToken;
    }

    public void setFileStorage(String fileStorage) {
        this.fileStorage = fileStorage;
    }

    public void setAuthToken(String authToken) {
        this.authToken = authToken;
    }
}
