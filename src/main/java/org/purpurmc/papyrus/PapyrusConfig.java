package org.purpurmc.papyrus;

import org.spongepowered.configurate.CommentedConfigurationNode;
import org.spongepowered.configurate.ConfigurateException;
import org.spongepowered.configurate.hocon.HoconConfigurationLoader;

import java.io.File;
import java.io.IOException;
import java.nio.file.Path;

public class PapyrusConfig {

    public static boolean needsSetup() {
        return !(new File("/etc/papyrus.conf").exists());
    }

    public static void setup() throws IOException {
        File config = new File("/etc/papyrus.conf");
        config.getParentFile().mkdirs();
        config.createNewFile();

        PapyrusConfig.load(true);
    }

    private static double version = 2.0;

    public static void load(boolean save) throws ConfigurateException {
        HoconConfigurationLoader loader = HoconConfigurationLoader.builder().path(Path.of("/etc/papyrus.conf")).build();
        CommentedConfigurationNode node = loader.load();
        node.options().shouldCopyDefaults(true);

        version = node.node("_version").getDouble(version);
        node.node("_version").set(version);

        PapyrusConfig.loadAPI(node.node("api"));

        if (save || new File("/etc/papyrus.conf").canWrite()) {
            loader.save(node);
        }
    }

    public static String host = "localhost";
    public static int port = 8080;

    public static String routePrefix = "/v1";
    public static String projectsRoute = "/projects";
    public static String projectRoute = "/project/{project}";
    public static String projectGroupRoute = "/project/{project}/group/{group}";
    public static String projectVersionRoute = "/project/{project}/version/{version}";
    public static String projectVersionBuildRoute = "/project/{project}/version/{version}/build/{build}";
    public static String projectVersionBuildDownloadRoute = "/project/{project}/version/{version}/build/{build}/download/{file}";

    private static void loadAPI(CommentedConfigurationNode node) {
        host = node.node("host").getString(host);
        port = node.node("port").getInt(port);

        CommentedConfigurationNode routes = node.node("routes").comment(
                """
                        By default routes have all variables they require filled in
                        To disable a route, set the value to ""
                """
        );

        routePrefix = routes.node("prefix").getString(routePrefix);
        projectsRoute = routes.node("projects").getString(projectsRoute);
        projectRoute = routes.node("project").getString(projectRoute);
        projectGroupRoute = routes.node("projectGroup").getString(projectGroupRoute);
        projectVersionRoute = routes.node("projectVersion").getString(projectVersionRoute);
        projectVersionBuildRoute = routes.node("projectVersionBuild").getString(projectVersionBuildRoute);
        projectVersionBuildDownloadRoute = routes.node("projectVersionBuildDownload").getString(projectVersionBuildDownloadRoute);
    }
}
