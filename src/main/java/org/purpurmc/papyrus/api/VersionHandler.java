package org.purpurmc.papyrus.api;

import io.javalin.http.Context;
import org.jetbrains.annotations.NotNull;

public class VersionHandler {

    public void getProjectVersion(@NotNull Context context) {
        context.result("get project version: " + context.pathParam("project") + ", " + context.pathParam("version"));
    }

    public void getProjectGroup(@NotNull Context context) {
        context.result("get project group: " + context.pathParam("project") + ", " + context.pathParam("group"));
    }
}
