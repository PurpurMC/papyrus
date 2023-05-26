package org.purpurmc.papyrus.db.repository;

import org.purpurmc.papyrus.db.entity.Project;
import org.purpurmc.papyrus.db.entity.Version;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface VersionRepository extends JpaRepository<Version, UUID> {
    Optional<Version> findByProjectAndName(Project project, String name);

    List<Version> findAllByProject(Project project);
}
