package org.purpurmc.papyrus.jenkins.uploader.response;

public class CreateBuildResponse {
    private String error;
    private String status;
    private String buildId;

    public String getError() {
        return error;
    }

    public String getStatus() {
        return status;
    }

    public String getBuildId() {
        return buildId;
    }
}
