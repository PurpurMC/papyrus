package org.purpurmc.papyrus.controller.v2;

import org.apache.commons.lang3.RandomStringUtils;
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
    public ListProjects listProjects() {
        List<Project> projects = projectRepository.findAll();
        return new ListProjects(projects.stream().map(Project::getName).toList());
    }

    private record ListProjects(List<String> projects) {
    }

    @GetMapping("/{project}")
    @ResponseBody
    public GetProject getProject(@PathVariable("project") String projectName) {
        Project project = projectRepository.findByName(projectName).orElseThrow(ProjectNotFound::new);
        List<Version> versions = versionRepository.findAllByProject(project);

        return new GetProject(project.getName(), versions.stream().map(Version::getName).toList());
    }

    private record GetProject(String project, List<String> versions) {
    }
}
