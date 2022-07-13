package org.purpurmc.papyrus.jenkins.uploader;

import com.google.gson.FieldNamingPolicy;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import org.purpurmc.papyrus.jenkins.uploader.payload.CreateBuildPayload;
import org.purpurmc.papyrus.jenkins.uploader.payload.UploadFilePayload;
import org.purpurmc.papyrus.jenkins.uploader.response.CreateBuildResponse;
import org.purpurmc.papyrus.jenkins.uploader.response.UploadFileResponse;
import org.purpurmc.papyrus.jenkins.util.MultipartBodyPublisher;
import org.purpurmc.papyrus.jenkins.util.Result;

import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.Objects;

public class PapyrusUploader {
    private final Gson gson;
    private final HttpClient client;

    private final String url;
    private final String key;
    private final String project;
    private final String version;
    private final String file;

    public PapyrusUploader(String url, String key, String project, String version, String file) {
        this.gson = new GsonBuilder().setFieldNamingPolicy(FieldNamingPolicy.LOWER_CASE_WITH_UNDERSCORES).create();
        this.client = HttpClient.newHttpClient();

        if (url.endsWith("/")) {
            url = url.substring(0, url.length() - 1);
        }

        this.url = url;
        this.key = key;
        this.project = project;
        this.version = version;
        this.file = file;
    }

    public Result<String, String> create(CreateBuildPayload payload) {
        URI url;
        try {
            url = new URI(this.url + "/v2/upload/create");
        } catch (URISyntaxException e) {
            return Result.error(e.getMessage());
        }

        HttpRequest request = HttpRequest.newBuilder()
                .uri(url)
                .header("Authorization", "Token " + this.key)
                .header("Content-Type", "application/json")
                .POST(HttpRequest.BodyPublishers.ofString(gson.toJson(payload)))
                .build();

        HttpResponse<String> responsePayload;
        try {
            responsePayload = client.send(request, HttpResponse.BodyHandlers.ofString());
        } catch (IOException | InterruptedException e) {
            return Result.error(e.getMessage());
        }

        CreateBuildResponse response = gson.fromJson(responsePayload.body(), CreateBuildResponse.class);
        if (responsePayload.statusCode() != 200 || response.getError() != null) {
            return Result.error(response.getError() == null ? response.getStatus() : response.getError());
        }

        return Result.ok(response.getBuildId());
    }

    public Result<Object, String> upload(UploadFilePayload payload) {
        URI url;
        try {
            url = new URI(this.url + "/v2/upload/file");
        } catch (URISyntaxException e) {
            return Result.error(e.getMessage());
        }

        MultipartBodyPublisher publisher = new MultipartBodyPublisher();
        publisher.addPart("build_id", payload.getBuildId());
        publisher.addPart("file", payload.getFile().toPath());
        publisher.addPart("file_extension", payload.getFileExtension());

        HttpRequest request = HttpRequest.newBuilder()
                .uri(url)
                .header("Authorization", "Token " + this.key)
                .header("Content-Type", "multipart/form-data; boundary=" + publisher.getBoundary())
                .POST(publisher.build())
                .build();

        HttpResponse<String> responsePayload;
        try {
            responsePayload = client.send(request, HttpResponse.BodyHandlers.ofString());
        } catch (IOException | InterruptedException e) {
            return Result.error(e.getMessage());
        }

        UploadFileResponse response = gson.fromJson(responsePayload.body(), UploadFileResponse.class);
        if (responsePayload.statusCode() != 200 || response.getError() != null) {
            return Result.error(response.getError() == null ? response.getStatus() : response.getError());
        }

        return Result.ok(null);
    }
}
