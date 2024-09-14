package org.purpurmc.papyrus.controller.v2;

import io.swagger.v3.oas.annotations.Hidden;
import java.util.HashSet;
import java.util.Set;
import java.util.function.Predicate;
import org.purpurmc.papyrus.config.AppConfiguration;
import org.purpurmc.papyrus.db.entity.Build;
import org.purpurmc.papyrus.db.entity.Metadata;
import org.purpurmc.papyrus.db.entity.Project;
import org.purpurmc.papyrus.db.entity.ProjectMetadata;
import org.purpurmc.papyrus.db.entity.Version;
import org.purpurmc.papyrus.db.repository.BuildRepository;
import org.purpurmc.papyrus.db.repository.MetadataRepository;
import org.purpurmc.papyrus.db.repository.ProjectMetadataRepository;
import org.purpurmc.papyrus.db.repository.ProjectRepository;
import org.purpurmc.papyrus.db.repository.VersionRepository;
import org.purpurmc.papyrus.exception.BuildNotFound;
import org.purpurmc.papyrus.exception.InvalidAuthToken;
import org.purpurmc.papyrus.exception.ProjectNotFound;
import org.purpurmc.papyrus.exception.VersionNotFound;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpHeaders;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

@Hidden
@RestController
@RequestMapping("/v2/metadata")
public class MetadataController {
    private final AppConfiguration configuration;
    private final ProjectRepository projectRepository;
    private final VersionRepository versionRepository;
    private final BuildRepository buildRepository;
    private final MetadataRepository metadataRepository;
    private final ProjectMetadataRepository projectMetadataRepository;

    @Autowired
    public MetadataController(
            AppConfiguration configuration,
            ProjectRepository projectRepository,
            VersionRepository versionRepository,
            BuildRepository buildRepository,
            ProjectMetadataRepository projectMetadataRepository,
            MetadataRepository metadataRepository) {
        this.configuration = configuration;
        this.projectRepository = projectRepository;
        this.versionRepository = versionRepository;
        this.buildRepository = buildRepository;
        this.projectMetadataRepository = projectMetadataRepository;
        this.metadataRepository = metadataRepository;
    }

    @PutMapping("project/{project}")
    @ResponseBody
    public ResponseEntity<String> updateProjectMetadata(@RequestHeader(HttpHeaders.AUTHORIZATION) String authHeader, @PathVariable("project") String projectName, @RequestBody UpdateMetadataBody body) {
        this.requireAuth(authHeader);

        Project project = projectRepository.findByName(projectName).orElseThrow(ProjectNotFound::new);

        List<ProjectMetadata> oldMetadata = projectMetadataRepository.findByProject(project);
        Map<String, String> newMetadata = body.metadata();

        if (newMetadata.isEmpty()) {
            projectMetadataRepository.deleteAll(oldMetadata);
            return ResponseEntity.ok("");
        }

        if (oldMetadata.isEmpty()) {
            List<ProjectMetadata> metadata = body.metadata().entrySet().stream()
                    .map(entry -> new ProjectMetadata(project, entry.getKey(), entry.getValue()))
                    .toList();
            projectMetadataRepository.saveAll(metadata);
            return ResponseEntity.ok("");
        }

        List<ProjectMetadata> updatedMetadata = oldMetadata.stream()
                .filter(metadata -> newMetadata.containsKey(metadata.getName()))
                .filter(metadata -> !newMetadata.get(metadata.getName()).equals(metadata.getValue()))
                .map(metadata -> {
                    metadata.setValue(newMetadata.get(metadata.getName()));
                    return metadata;
                })
                .toList();
        List<ProjectMetadata> deletedMetadata = oldMetadata.stream()
                .filter(metadata -> !newMetadata.containsKey(metadata.getName()))
                .toList();

        List<String> existingKeys = oldMetadata.stream().map(ProjectMetadata::getName).toList();
        List<ProjectMetadata> addedMetadata = newMetadata.entrySet().stream()
                .filter(entry -> !existingKeys.contains(entry.getKey()))
                .map(entry -> new ProjectMetadata(project, entry.getKey(), entry.getValue()))
                .toList();

        if (!deletedMetadata.isEmpty()) {
            projectMetadataRepository.deleteAll(deletedMetadata);
        }

        projectMetadataRepository.saveAll(updatedMetadata);
        projectMetadataRepository.saveAll(addedMetadata);

        return ResponseEntity.ok("");
    }

    @PutMapping("/build/{project}/{version}/{build}")
    @ResponseBody
    public ResponseEntity<String> updateBuildMetadata(@RequestHeader(HttpHeaders.AUTHORIZATION) String authHeader, @PathVariable("project") String projectName, @PathVariable("version") String versionName, @PathVariable("build") String buildName, @RequestBody UpdateMetadataBody body) {
        this.requireAuth(authHeader);

        Project project = projectRepository.findByName(projectName).orElseThrow(ProjectNotFound::new);
        Version version = versionRepository.findByProjectAndName(project, versionName).orElseThrow(VersionNotFound::new);
        Build build = (buildName.equals("latest")
                ? buildRepository.findLatestByVersionAndFileNotNull(version)
                : buildRepository.findByVersionAndNameAndReady(version, buildName)
        ).orElseThrow(BuildNotFound::new);
        
        List<Metadata> oldMetadata = metadataRepository.findByBuild(build);
        Map<String, String> newMetadata = body.metadata();

        if (newMetadata.isEmpty()) {
            metadataRepository.deleteAll(oldMetadata);
            return ResponseEntity.ok("");
        }

        if (oldMetadata.isEmpty()) {
            List<Metadata> metadata = body.metadata().entrySet().stream()
                    .map(entry -> new Metadata(build, entry.getKey(), entry.getValue()))
                    .toList();
            metadataRepository.saveAll(metadata);
            return ResponseEntity.ok("");
        }

        List<Metadata> updatedMetadata = oldMetadata.stream()
                .filter(metadata -> newMetadata.containsKey(metadata.getName()))
                .filter(metadata -> !newMetadata.get(metadata.getName()).equals(metadata.getValue()))
                .map(metadata -> {
                    metadata.setValue(newMetadata.get(metadata.getName()));
                    return metadata;
                })
                .toList();
        List<Metadata> deletedMetadata = oldMetadata.stream()
                .filter(metadata -> !newMetadata.containsKey(metadata.getName()))
                .toList();

        List<String> existingKeys = oldMetadata.stream().map(Metadata::getName).toList();
        List<Metadata> addedMetadata = newMetadata.entrySet().stream()
                .filter(entry -> !existingKeys.contains(entry.getKey()))
                .map(entry -> new Metadata(build, entry.getKey(), entry.getValue()))
                .toList();

        if (!deletedMetadata.isEmpty()) {
            metadataRepository.deleteAll(deletedMetadata);
        }

        metadataRepository.saveAll(updatedMetadata);
        metadataRepository.saveAll(addedMetadata);

        return ResponseEntity.ok("");
    }

    private void requireAuth(String authHeader) {
        String[] parts = authHeader.trim().split(" ");
        if (parts.length != 2) {
            throw new InvalidAuthToken();
        }

        if (!parts[0].equals("Basic")) {
            throw new InvalidAuthToken();
        }

        if (!parts[1].equals(configuration.getAuthToken())) {
            throw new InvalidAuthToken();
        }
    }

    private record UpdateMetadataBody(Map<String, String> metadata) {
    }
}
