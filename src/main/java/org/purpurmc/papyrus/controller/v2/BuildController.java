package org.purpurmc.papyrus.controller.v2;

import org.apache.commons.lang3.RandomStringUtils;
import org.purpurmc.papyrus.config.AppConfiguration;
import org.purpurmc.papyrus.db.entity.Build;
import org.purpurmc.papyrus.db.entity.Commit;
import org.purpurmc.papyrus.db.entity.File;
import org.purpurmc.papyrus.db.entity.Project;
import org.purpurmc.papyrus.db.entity.Version;
import org.purpurmc.papyrus.db.repository.BuildRepository;
import org.purpurmc.papyrus.db.repository.CommitRepository;
import org.purpurmc.papyrus.db.repository.FileRepository;
import org.purpurmc.papyrus.db.repository.ProjectRepository;
import org.purpurmc.papyrus.db.repository.VersionRepository;
import org.purpurmc.papyrus.exception.BuildNotFound;
import org.purpurmc.papyrus.exception.FileDownloadError;
import org.purpurmc.papyrus.exception.ProjectNotFound;
import org.purpurmc.papyrus.exception.VersionNotFound;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.core.io.Resource;
import org.springframework.http.ContentDisposition;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;
import java.util.Random;

@RestController
@RequestMapping("/v2/{project}/{version}")
public class BuildController {
    private final AppConfiguration configuration;
    private final ProjectRepository projectRepository;
    private final VersionRepository versionRepository;
    private final BuildRepository buildRepository;
    private final CommitRepository commitRepository;
    private final FileRepository fileRepository;

    @Autowired
    public BuildController(AppConfiguration configuration, ProjectRepository projectRepository, VersionRepository versionRepository, BuildRepository buildRepository, CommitRepository commitRepository, FileRepository fileRepository) {
        this.configuration = configuration;
        this.projectRepository = projectRepository;
        this.versionRepository = versionRepository;
        this.buildRepository = buildRepository;
        this.commitRepository = commitRepository;
        this.fileRepository = fileRepository;
    }

    @GetMapping("/{build}")
    @ResponseBody
    public GetBuild getBuild(@PathVariable("project") String projectName, @PathVariable("version") String versionName, @PathVariable("build") String buildName) {
        Project project = projectRepository.findByName(projectName).orElseThrow(ProjectNotFound::new);
        Version version = versionRepository.findByProjectAndName(project, versionName).orElseThrow(VersionNotFound::new);
        Build build = (buildName.equals("latest")
                ? buildRepository.findLatestByVersionAndFileNotNull(version)
                : buildRepository.findByVersionAndNameAndFileNotNull(version, buildName)
        ).orElseThrow(BuildNotFound::new);
        List<Commit> commits = commitRepository.findAllByBuild(build);

        List<GetBuild.Commits> responseCommits = commits.stream().map(commit -> new GetBuild.Commits(commit.getAuthor(), commit.getEmail(), commit.getDescription(), commit.getHash(), commit.getTimestamp())).toList();
        return new GetBuild(project.getName(), version.getName(), build.getName(), build.getResult().toString(), build.getTimestamp(), build.getDuration(), responseCommits, build.getHash());
    }

    private record GetBuild(String project, String version, String build, String result, Long timestamp, Long duration,
                            List<Commits> commits, String md5) {
        public record Commits(String author, String email, String description, String hash, Long timestamp) {
        }
    }

    @GetMapping("/{build}/download")
    @ResponseBody
    public ResponseEntity<Resource> downloadBuild(@PathVariable("project") String projectName, @PathVariable("version") String versionName, @PathVariable("build") String buildName) throws IOException {
        Project project = projectRepository.findByName(projectName).orElseThrow(ProjectNotFound::new);
        Version version = versionRepository.findByProjectAndName(project, versionName).orElseThrow(VersionNotFound::new);
        Build build = (buildName.equals("latest")
                ? buildRepository.findLatestByVersionAndFileNotNull(version)
                : buildRepository.findByVersionAndNameAndFileNotNull(version, buildName)
        ).orElseThrow(BuildNotFound::new);
        File file = fileRepository.findByBuild(build).orElseThrow(BuildNotFound::new);

        MediaType mediaType;
        try {
            mediaType = MediaType.parseMediaType(file.getContentType());
        } catch (Exception e) {
            mediaType = MediaType.APPLICATION_OCTET_STREAM;
        }

        ByteArrayResource resource;
        try {
            Path path = Path.of(configuration.getFileStorage(), file.getId().toString());
            byte[] bytes = Files.readAllBytes(path);
            resource = new ByteArrayResource(bytes);
        } catch (Exception e) {
            throw new FileDownloadError();
        }

        String filename = String.format("%s-%s-%s.%s", project.getName(), version.getName(), build.getName(), file.getFileExtension());

        return ResponseEntity.ok()
                .contentType(mediaType)
                .contentLength(resource.contentLength())
                .header(HttpHeaders.CONTENT_DISPOSITION, ContentDisposition.attachment().filename(filename).build().toString())
                .body(resource);
    }
}
