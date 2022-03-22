package org.purpurmc.papyrus.api;

import io.javalin.http.Context;
import org.jetbrains.annotations.NotNull;

public class BuildHandler {

    public void getVersionBuild(@NotNull Context context) {
        context.result("get version build: " + context.pathParam("project") + ", " + context.pathParam("version") + ", " + context.pathParam("build"));
    }

    public void downloadVersionBuild(@NotNull Context context) {
        context.result("download version build: " + context.pathParam("project") + ", " + context.pathParam("version") + ", " + context.pathParam("build") + ", " + context.pathParam("file"));
    }
}
