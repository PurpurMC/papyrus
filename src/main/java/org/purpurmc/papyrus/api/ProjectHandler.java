package org.purpurmc.papyrus.api;

import io.javalin.http.Context;
import org.jetbrains.annotations.NotNull;

public class ProjectHandler {

    public void listProjects(@NotNull Context context) {
        context.result("list of projects");
    }

    public void getProject(@NotNull Context context) {
        context.result("specific project: " + context.pathParam("project"));
    }
}
