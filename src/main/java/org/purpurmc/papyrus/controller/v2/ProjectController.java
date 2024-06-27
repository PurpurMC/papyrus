package org.purpurmc.papyrus.controller.v2;

import io.swagger.v3.oas.annotations.Operation;
import org.purpurmc.papyrus.db.entity.Project;
import org.purpurmc.papyrus.db.entity.Version;
import org.purpurmc.papyrus.db.repository.ProjectRepository;
import org.purpurmc.papyrus.db.repository.VersionRepository;
import org.purpurmc.papyrus.exception.ProjectNotFound;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping("/v2")
public class ProjectController {
    private final ProjectRepository projectRepository;
    private final VersionRepository versionRepository;

    @Autowired
    public ProjectController(ProjectRepository projectRepository, VersionRepository versionRepository) {
        this.projectRepository = projectRepository;
        this.versionRepository = versionRepository;
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

        return new ProjectResponse(project.getName(), versions.stream().map(Version::getName).toList());
    }

    private record ProjectsResponse(List<String> projects) {
    }

    private record ProjectResponse(String project, List<String> versions) {
    }
}
