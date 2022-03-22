package org.purpurmc.papyrus.api;

import io.javalin.http.Context;
import org.jetbrains.annotations.NotNull;

public class BuildHandler {

    public void listProjectBuilds(@NotNull Context context) {
        context.result("list of builds for project: " + context.pathParam("project"));
    }

    public void getProjectBuild(@NotNull Context context) {
        context.result("specific build for project: " + context.pathParam("project") + ", " + context.pathParam("build"));
    }

    public void downloadProjectBuild(@NotNull Context context) {
        context.result("download build for project: " + context.pathParam("project") + ", " + context.pathParam("build") + ", " + context.pathParam("file"));
    }
    public void listVersionBuilds(@NotNull Context context) {
        context.result("list of version builds: " + context.pathParam("project") + ", " + context.pathParam("version"));
    }

    public void getVersionBuild(@NotNull Context context) {
        context.result("get version build: " + context.pathParam("project") + ", " + context.pathParam("version") + ", " + context.pathParam("build"));
    }

    public void downloadVersionBuild(@NotNull Context context) {
        context.result("download version build: " + context.pathParam("project") + ", " + context.pathParam("version") + ", " + context.pathParam("build") + ", " + context.pathParam("file"));
    }
}
