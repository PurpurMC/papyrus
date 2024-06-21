package org.purpurmc.papyrus.controller.v2;

import io.swagger.v3.oas.annotations.Operation;
import org.purpurmc.papyrus.db.entity.Build;
import org.purpurmc.papyrus.db.entity.Commit;
import org.purpurmc.papyrus.db.entity.Metadata;
import org.purpurmc.papyrus.db.entity.Project;
import org.purpurmc.papyrus.db.entity.Version;
import org.purpurmc.papyrus.db.repository.BuildRepository;
import org.purpurmc.papyrus.db.repository.CommitRepository;
import org.purpurmc.papyrus.db.repository.MetadataRepository;
import org.purpurmc.papyrus.db.repository.ProjectRepository;
import org.purpurmc.papyrus.db.repository.VersionRepository;
import org.purpurmc.papyrus.exception.ProjectNotFound;
import org.purpurmc.papyrus.exception.VersionNotFound;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;

@RestController
@RequestMapping("/v2/{project}")
public class VersionController {
    private final ProjectRepository projectRepository;
    private final VersionRepository versionRepository;
    private final BuildRepository buildRepository;
    private final CommitRepository commitRepository;
    private final MetadataRepository metadataRepository;

    @Autowired
    public VersionController(ProjectRepository projectRepository, VersionRepository versionRepository, BuildRepository buildRepository, CommitRepository commitRepository, MetadataRepository metadataRepository) {
        this.projectRepository = projectRepository;
        this.versionRepository = versionRepository;
        this.buildRepository = buildRepository;
        this.commitRepository = commitRepository;
        this.metadataRepository = metadataRepository;
    }

    @GetMapping("/{version}")
    @ResponseBody
    @Operation(summary = "Get a project's version")
    public ResponseEntity<?> getVersion(@PathVariable("project") String projectName, @PathVariable("version") String versionName, @RequestParam(value = "detailed", required = false) String detailed) {
        Project project = projectRepository.findByName(projectName).orElseThrow(ProjectNotFound::new);
        Version version = versionRepository.findByProjectAndName(project, versionName).orElseThrow(VersionNotFound::new);
        List<Build> builds = buildRepository.findAllByVersionAndReadyOrderByTimestampAsc(version);
        Optional<Build> latest = buildRepository.findLatestByVersionAndFileNotNull(version);

        if (detailed != null) {
            return ResponseEntity.ok(
                    new VersionResponseDetailed(
                            project.getName(),
                            version.getName(),
                            new VersionResponseDetailed.VersionBuildsDetailed(
                                    latest.map(build -> convertToBuildResponse(project, version, build)),
                                    builds.stream().map(build -> convertToBuildResponse(project, version, build)).toList()
                            )
                    )
            );
        } else {
            return ResponseEntity.ok(
                    new VersionResponse(
                            project.getName(),
                            version.getName(),
                            new VersionResponse.VersionBuilds(
                                    latest.map(Build::getName),
                                    builds.stream().map(Build::getName).toList()
                            )
                    )
            );
        }
    }

    private BuildController.BuildResponse convertToBuildResponse(Project project, Version version, Build build) {
        List<Commit> commits = commitRepository.findAllByBuild(build);
        List<BuildController.BuildResponse.BuildCommits> responseCommits = commits.stream().map(commit -> new BuildController.BuildResponse.BuildCommits(commit.getAuthor(), commit.getEmail(), commit.getDescription(), commit.getHash(), commit.getTimestamp())).toList();

        List<Metadata> metadata = metadataRepository.findByBuild(build);
        Map<String, String> responseMetadata = new HashMap<>();

        for (Metadata data : metadata) {
            responseMetadata.put(data.getName(), data.getValue());
        }

        return new BuildController.BuildResponse(project.getName(), version.getName(), build.getName(), build.getResult().toString(), build.getTimestamp(), build.getDuration(), responseCommits, responseMetadata, build.getHash());
    }

    private record VersionResponse(String project, String version, VersionBuilds builds) {
        public record VersionBuilds(Optional<String> latest, List<String> all) {
        }
    }

    private record VersionResponseDetailed(String project, String version, VersionBuildsDetailed builds) {
        public record VersionBuildsDetailed(Optional<BuildController.BuildResponse> latest,
                                            List<BuildController.BuildResponse> all) {
        }
    }
}
