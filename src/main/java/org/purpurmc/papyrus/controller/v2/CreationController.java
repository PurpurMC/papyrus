package org.purpurmc.papyrus.controller.v2;

import io.swagger.v3.oas.annotations.Hidden;
import org.purpurmc.papyrus.config.AppConfiguration;
import org.purpurmc.papyrus.db.entity.Build;
import org.purpurmc.papyrus.db.entity.Commit;
import org.purpurmc.papyrus.db.entity.CreationState;
import org.purpurmc.papyrus.db.entity.File;
import org.purpurmc.papyrus.db.entity.Metadata;
import org.purpurmc.papyrus.db.entity.Project;
import org.purpurmc.papyrus.db.entity.Version;
import org.purpurmc.papyrus.db.repository.BuildRepository;
import org.purpurmc.papyrus.db.repository.CommitRepository;
import org.purpurmc.papyrus.db.repository.CreationStateRepository;
import org.purpurmc.papyrus.db.repository.FileRepository;
import org.purpurmc.papyrus.db.repository.MetadataRepository;
import org.purpurmc.papyrus.db.repository.ProjectRepository;
import org.purpurmc.papyrus.db.repository.VersionRepository;
import org.purpurmc.papyrus.exception.BuildAlreadyExists;
import org.purpurmc.papyrus.exception.FileUploadError;
import org.purpurmc.papyrus.exception.InvalidStateKey;
import org.purpurmc.papyrus.util.AuthUtil;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpHeaders;
import org.springframework.http.ResponseEntity;
import org.springframework.util.DigestUtils;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.UUID;

@Hidden
@RestController
@RequestMapping("/v2/create")
public class CreationController {
    private final AppConfiguration configuration;
    private final ProjectRepository projectRepository;
    private final VersionRepository versionRepository;
    private final BuildRepository buildRepository;
    private final CommitRepository commitRepository;
    private final MetadataRepository metadataRepository;
    private final FileRepository fileRepository;
    private final CreationStateRepository creationStateRepository;

    @Autowired
    public CreationController(
            AppConfiguration configuration,
            ProjectRepository projectRepository,
            VersionRepository versionRepository,
            BuildRepository buildRepository,
            CommitRepository commitRepository,
            MetadataRepository metadataRepository,
            FileRepository fileRepository,
            CreationStateRepository creationStateRepository
    ) {
        this.configuration = configuration;
        this.projectRepository = projectRepository;
        this.versionRepository = versionRepository;
        this.buildRepository = buildRepository;
        this.commitRepository = commitRepository;
        this.metadataRepository = metadataRepository;
        this.fileRepository = fileRepository;
        this.creationStateRepository = creationStateRepository;
    }

    @PostMapping
    @ResponseBody
    public CreateBuild createBuild(@RequestHeader(HttpHeaders.AUTHORIZATION) String authHeader, @RequestBody CreateBuildBody body) {
        AuthUtil.requireAuth(configuration, authHeader);

        Project project = null;
        Version version = null;

        Optional<Project> projectOption = projectRepository.findByName(body.project);
        if (projectOption.isPresent()) {
            project = projectOption.get();
            Optional<Version> versionOption = versionRepository.findByProjectAndName(project, body.version);
            if (versionOption.isPresent()) {
                version = versionOption.get();
                if (buildRepository.existsByVersionAndName(version, body.build)) {
                    throw new BuildAlreadyExists();
                }
            }
        }

        if (project == null) {
            project = projectRepository.save(new Project(body.project));
        }

        if (version == null) {
            version = versionRepository.save(new Version(project, body.version));
        }

        int ready = 0;
        if (body.result != Build.BuildResult.SUCCESS) {
            ready = 1;
        }

        Build build = buildRepository.save(new Build(version, body.build, body.result, body.timestamp, body.duration, ready));
        if (body.commits != null) {
            commitRepository.saveAll(body.commits.stream().map(commit -> new Commit(build, commit.author, commit.email, commit.description, commit.hash, commit.timestamp)).toList());
        }

        if (body.metadata != null) {
            List<Metadata> metadata = new ArrayList<>();
            body.metadata.forEach((name, value) -> metadata.add(new Metadata(build, name, value)));
            metadataRepository.saveAll(metadata);
        }

        if (ready == 1) {
            return new CreateBuild(null);
        }

        CreationState id = creationStateRepository.save(new CreationState(build, body.fileExtension));
        return new CreateBuild(id.getId().toString());
    }

    @PostMapping("/upload")
    @ResponseBody
    public ResponseEntity<String> uploadFile(@RequestHeader(HttpHeaders.AUTHORIZATION) String authHeader, @RequestParam("stateKey") String stateKey, @RequestParam("file") MultipartFile uploadFile) {
        AuthUtil.requireAuth(configuration, authHeader);

        CreationState state;
        try {
            state = creationStateRepository.getStateById(UUID.fromString(stateKey)).orElseThrow(InvalidStateKey::new);
        } catch (Exception e) {
            throw new InvalidStateKey();
        }

        byte[] bytes;
        try {
            bytes = uploadFile.getBytes();
        } catch (Exception e) {
            throw new FileUploadError();
        }

        String contentType;
        try {
            Path tempFile = Files.createTempFile("papyrus", state.getId().toString());
            Files.write(tempFile, bytes);
            contentType = Files.probeContentType(tempFile);
            Files.deleteIfExists(tempFile);
        } catch (Exception e) {
            throw new FileUploadError();
        }

        Build build = state.getBuild();

        File file = fileRepository.save(new File(build, contentType, state.getFileExtension()));
        try {
            Path path = Path.of(configuration.getFileStorage(), file.getId().toString());
            Files.write(path, bytes);
        } catch (Exception e) {
            throw new FileUploadError();
        }

        build.setHash(DigestUtils.md5DigestAsHex(bytes));
        build.setReady(1);
        buildRepository.save(build);

        creationStateRepository.delete(state);
        return ResponseEntity.ok("");
    }

    private record CreateBuildBody(String project,
                                   String version,
                                   String build,
                                   Build.BuildResult result,
                                   long timestamp,
                                   long duration,
                                   List<CommitBody> commits,
                                   Map<String, String> metadata,
                                   String fileExtension) {
        public record CommitBody(String author, String email, String description, String hash, long timestamp) {
        }
    }

    private record CreateBuild(String stateKey) {
    }
}
