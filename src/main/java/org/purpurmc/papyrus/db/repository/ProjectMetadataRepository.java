package org.purpurmc.papyrus.db.repository;

import java.util.List;
import java.util.UUID;
import org.purpurmc.papyrus.db.entity.ProjectMetadata;
import org.purpurmc.papyrus.db.entity.Project;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface ProjectMetadataRepository extends JpaRepository<ProjectMetadata, UUID> {

    List<ProjectMetadata> findByProject(Project project);
}
