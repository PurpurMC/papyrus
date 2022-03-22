package org.purpurmc.papyrus.commands;

import io.javalin.Javalin;
import org.purpurmc.papyrus.api.BuildHandler;
import org.purpurmc.papyrus.api.ProjectHandler;
import org.purpurmc.papyrus.api.VersionHandler;
import picocli.CommandLine;

import static io.javalin.apibuilder.ApiBuilder.*;

@CommandLine.Command(
        name = "api",
        description = "Start the Papyrus web API"
)
public class APICommand implements Runnable {

    @CommandLine.Option(
            names = {"-p", "--port"},
            description = "The port to run the API on",
            defaultValue = "8000"
    )
    private int port;

    @Override
    public void run() {
        Javalin app = Javalin.create();

        ProjectHandler projects = new ProjectHandler();
        VersionHandler versions = new VersionHandler();
        BuildHandler builds = new BuildHandler();

        app.routes(() -> path("/v2", () -> {
            get("/projects", projects::listProjects);
            get("/project/{project}", projects::getProject);

            get("/project/{project}/versions", versions::listProjectVersions);
            get("/project/{project}/version/{version}", versions::getProjectVersion);

            get("/project/{project}/groups", versions::listProjectGroups);
            get("/project/{project}/group/{group}", versions::getProjectGroup);

            get("/project/{project}/builds", builds::listProjectBuilds);
            get("/project/{project}/build/{build}", builds::getProjectBuild);
            get("/project/{project}/build/{build}/download/{file}", builds::downloadProjectBuild);

            get("/project/{project}/version/{version}/builds", builds::listVersionBuilds);
            get("/project/{project}/version/{version}/build/{build}", builds::getVersionBuild);
            get("/project/{project}/version/{version}/build/{build}/download/{file}", builds::downloadVersionBuild);
        }));

        app.start(port);
    }
}
