package org.purpurmc.papyrus.commands;

import io.javalin.Javalin;
import io.javalin.http.Handler;
import org.purpurmc.papyrus.PapyrusConfig;
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

    @Override
    public void run() {
        Javalin app = Javalin.create();

        ProjectHandler projects = new ProjectHandler();
        VersionHandler versions = new VersionHandler();
        BuildHandler builds = new BuildHandler();

        app.routes(() -> path(PapyrusConfig.routePrefix, () -> {
            registerRoute(PapyrusConfig.projectsRoute, projects::listProjects);
            registerRoute(PapyrusConfig.projectRoute, projects::getProject);
            registerRoute(PapyrusConfig.projectGroupRoute, versions::getProjectGroup);
            registerRoute(PapyrusConfig.projectVersionRoute, versions::getProjectVersion);
            registerRoute(PapyrusConfig.projectVersionBuildRoute, builds::getVersionBuild);
            registerRoute(PapyrusConfig.projectVersionBuildDownloadRoute, builds::downloadVersionBuild);
        }));

        app.start(PapyrusConfig.host, PapyrusConfig.port);
    }

    private void registerRoute(String route, Handler handler) {
        if (route.isBlank()) {
            return;
        }

        if (!route.startsWith("/")) {
            route = "/" + route;
        }

        if (route.endsWith("/")) {
            route = route.substring(0, route.length() - 1);
        }

        get(route, handler);
    }
}
