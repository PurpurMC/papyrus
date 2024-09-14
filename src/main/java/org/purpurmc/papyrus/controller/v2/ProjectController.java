package org.purpurmc.papyrus.controller.v2;

import io.swagger.v3.oas.annotations.Operation;
import java.util.HashMap;
import java.util.Map;
import org.purpurmc.papyrus.config.AppConfiguration;
import org.purpurmc.papyrus.db.entity.Project;
import org.purpurmc.papyrus.db.entity.ProjectMetadata;
import org.purpurmc.papyrus.db.entity.Version;
import org.purpurmc.papyrus.db.repository.ProjectMetadataRepository;
import org.purpurmc.papyrus.db.repository.ProjectRepository;
import org.purpurmc.papyrus.db.repository.VersionRepository;
import org.purpurmc.papyrus.exception.ProjectNotFound;
import org.purpurmc.papyrus.util.AuthUtil;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpHeaders;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping("/v2")
public class ProjectController {
    private final ProjectRepository projectRepository;
    private final ProjectMetadataRepository projectMetadataRepository;
    private final VersionRepository versionRepository;
    private final AppConfiguration configuration;

    @Autowired
    public ProjectController(ProjectRepository projectRepository, VersionRepository versionRepository,
                             ProjectMetadataRepository projectMetadataRepository, AppConfiguration configuration) {
        this.projectRepository = projectRepository;
        this.versionRepository = versionRepository;
        this.projectMetadataRepository = projectMetadataRepository;
        this.configuration = configuration;
    }

    @GetMapping
    @ResponseBody
    @Operation(summary = "List all projects")
    public ProjectsResponse listProjects() {
        List<Project> projects = projectRepository.findAll();
        return new ProjectsResponse(projects.stream().map(Project::getName).toList());
    }

    @GetMapping("/{project}")
    @ResponseBody
    @Operation(summary = "Get a project")
    public ProjectResponse getProject(@PathVariable("project") String projectName) {
        Project project = projectRepository.findByName(projectName).orElseThrow(ProjectNotFound::new);
        List<Version> versions = versionRepository.findAllByProject(project);

        List<ProjectMetadata> metadata = projectMetadataRepository.findByProject(project);
        Map<String, String> responseMetadata = new HashMap<>();

        for (ProjectMetadata data : metadata) {
            responseMetadata.put(data.getName(), data.getValue());
        }

        return new ProjectResponse(project.getName(), responseMetadata, versions.stream().map(Version::getName).toList());
    }

    @PutMapping("/{project}/metadata")
    @ResponseBody
    public ResponseEntity<String> updateProjectMetadata(@RequestHeader(HttpHeaders.AUTHORIZATION) String authHeader, @PathVariable("project") String projectName, @RequestBody UpdateMetadataBody body) {
        AuthUtil.requireAuth(configuration, authHeader);

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

    private record ProjectsResponse(List<String> projects) {
    }

    private record ProjectResponse(String project, Map<String, String> metadata, List<String> versions) {
    }

    private record UpdateMetadataBody(Map<String, String> metadata) {
    }
}
