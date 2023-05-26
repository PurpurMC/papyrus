package org.purpurmc.papyrus.db.repository;

import org.purpurmc.papyrus.db.entity.Build;
import org.purpurmc.papyrus.db.entity.File;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface FileRepository extends JpaRepository<File, UUID> {
    Optional<File> findByBuild(Build build);
}
