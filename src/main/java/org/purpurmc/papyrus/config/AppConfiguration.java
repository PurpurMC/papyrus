package org.purpurmc.papyrus.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

@Configuration
@ConfigurationProperties("app")
public class AppConfiguration {
    private String fileStorage;
    private String authToken;
    private String apiTitle;
    private String apiUrl;

    public String getFileStorage() {
        return this.fileStorage;
    }

    public void setFileStorage(String fileStorage) {
        this.fileStorage = fileStorage;
    }

    public String getAuthToken() {
        return this.authToken;
    }

    public void setAuthToken(String authToken) {
        this.authToken = authToken;
    }

    public String getApiTitle() {
        return this.apiTitle;
    }

    public void setApiTitle(String apiTitle) {
        this.apiTitle = apiTitle;
    }

    public String getApiUrl() {
        return this.apiUrl;
    }

    public void setApiUrl(String apiUrl) {
        this.apiUrl = apiUrl;
    }
}
