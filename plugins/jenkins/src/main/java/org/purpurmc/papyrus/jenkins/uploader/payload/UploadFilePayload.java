package org.purpurmc.papyrus.jenkins.uploader.payload;

import java.io.File;

public class UploadFilePayload {
    private String buildId;
    private File file;
    private String fileExtension;

    public void setBuildId(String buildId) {
        this.buildId = buildId;
    }

    public void setFile(File file) {
        this.file = file;
    }

    public void setFileExtension(String fileExtension) {
        this.fileExtension = fileExtension;
    }
}
