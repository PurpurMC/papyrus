package org.purpurmc.papyrus.api;

import io.javalin.http.Context;
import org.jetbrains.annotations.NotNull;

public class VersionHandler {

    public void listProjectVersions(@NotNull Context context) {
        context.result("list of project versions: " + context.pathParam("project"));
    }

    public void getProjectVersion(@NotNull Context context) {
        context.result("get project version: " + context.pathParam("project") + ", " + context.pathParam("version"));
    }

    public void listProjectGroups(@NotNull Context context) {
        context.result("list of project groups: " + context.pathParam("project"));
    }

    public void getProjectGroup(@NotNull Context context) {
        context.result("get project group: " + context.pathParam("project") + ", " + context.pathParam("group"));
    }
}
