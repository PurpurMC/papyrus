package org.purpurmc.papyrus.jenkins.uploader;

import org.purpurmc.papyrus.jenkins.uploader.payload.CreateBuildPayload;
import org.purpurmc.papyrus.jenkins.uploader.payload.UploadFilePayload;
import org.purpurmc.papyrus.jenkins.util.Result;

public class PapyrusUploader {
    private final String url;
    private final String key;
    private final String project;
    private final String version;
    private final String file;

    public PapyrusUploader(String url, String key, String project, String version, String file) {
        /*
        this.url = url;
        this.key = key;
        this.project = project;
        this.version = version;
        this.file = file;
         */

        // this is a hack to get around the fact that the jenkins dev server does not store the plugin's config
        this.url = "http://127.0.0.1:8000";
        this.key = "key";
        this.project = "purpur";
        this.version = "1.19";
        this.file = "output.txt";
    }

    public Result<String, String> create(CreateBuildPayload payload) {
        return Result.ok(""); // todo
    }

    public Result<Object, String> upload(UploadFilePayload payload) {
        return Result.ok(null); // todo
    }
}
